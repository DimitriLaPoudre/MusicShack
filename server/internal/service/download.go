package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	pkg_sync "github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type DownloadService struct {
	cfg      config.DownloadConfig
	tasks    sync.Map // map[uuid.UUID]sync.Map -> map[uuid.UUID]*downloadTask
	limit    pkg_sync.Semaphore
	plugin   *PluginService
	instance model.InstanceRepository
}

type downloadTask struct {
	userID         uuid.UUID
	song           model.EnrichedSong
	status         atomic.Value // model.DownloadStatus
	statusComment  atomic.Value // string
	downloadCancel atomic.Value // context.CancelFunc
	globalLimit    *pkg_sync.Semaphore
	plugin         *PluginService
	instance       model.InstanceRepository
}

func NewDownloadService(cfg config.DownloadConfig, plugin *PluginService, instance model.InstanceRepository) *DownloadService {
	return &DownloadService{
		cfg:      cfg,
		limit:    pkg_sync.NewSemaphore(cfg.Concurrent),
		plugin:   plugin,
		instance: instance,
	}
}

func (s *DownloadService) AddArtist(ctx context.Context, userID uuid.UUID, pluginInstances map[model.Plugin][]model.Instance, artistID string) ([]model.AddDownloadError, error) {
	artist, err := s.plugin.GetArtistFromPluginInstances(ctx, pluginInstances, artistID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get artist %s info for download: %w", artistID, err)
	}

	errList := []model.AddDownloadError{}
	for _, album := range artist.Albums {
		subErrList, err := s.AddAlbum(ctx, userID, pluginInstances, album.Id)
		if err != nil {
			errList = append(errList, model.AddDownloadError{
				Type:   model.TypeAlbum,
				ID:     album.Id,
				Name:   album.Title,
				Reason: err.Error(),
			})
		}
		errList = append(errList, subErrList...)
	}

	return errList, nil
}

func (s *DownloadService) AddAlbum(ctx context.Context, userID uuid.UUID, pluginInstances map[model.Plugin][]model.Instance, albumID string) ([]model.AddDownloadError, error) {
	album, err := s.plugin.GetAlbumFromPluginInstances(ctx, pluginInstances, albumID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get album %s info for download: %w", albumID, err)
	}

	errList := []model.AddDownloadError{}
	for _, song := range album.Songs {
		err := s.AddSong(ctx, userID, pluginInstances, song.Id)
		if err != nil {
			errList = append(errList, model.AddDownloadError{
				Type:   model.TypeSong,
				ID:     song.Id,
				Name:   song.Title,
				Reason: err.Error(),
			})
		}
	}

	return errList, nil
}

func (s *DownloadService) AddPlaylist(ctx context.Context, userID uuid.UUID, pluginInstances map[model.Plugin][]model.Instance, playlistID string) ([]model.AddDownloadError, error) {
	playlist, err := s.plugin.GetPlaylistFromPluginInstances(ctx, pluginInstances, playlistID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get playlist %s info for download: %w", playlistID, err)
	}

	errList := []model.AddDownloadError{}
	for _, song := range playlist.Songs {
		err := s.AddSong(ctx, userID, pluginInstances, song.Id)
		if err != nil {
			errList = append(errList, model.AddDownloadError{
				Type:   model.TypeSong,
				ID:     song.Id,
				Name:   song.Title,
				Reason: err.Error(),
			})
		}
	}

	return errList, nil
}

func (s *DownloadService) AddSong(ctx context.Context, userID uuid.UUID, pluginInstances map[model.Plugin][]model.Instance, songID string) error {
	song, err := s.plugin.GetSongFromPluginInstances(ctx, pluginInstances, songID)
	if err != nil {
		return fmt.Errorf("get song %s info for download: %w", songID, err)
	}

	if err := s.DownloadSong(userID, song); err != nil {
		return fmt.Errorf("download song %s: %w", song.Title, err)
	}

	return nil
}

func (s *DownloadService) DownloadSong(userID uuid.UUID, song model.EnrichedSong) error {
	// check if song already owned

	taskID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("create id for new download task: %w", err)
	}

	newTask := &downloadTask{
		userID:         userID,
		song:           song,
		status:         atomic.Value{},
		statusComment:  atomic.Value{},
		downloadCancel: atomic.Value{},
		globalLimit:    &s.limit,
		plugin:         s.plugin,
		instance:       s.instance,
	}

	newTask.status.Store(model.DownloadStatusPending)
	newTask.statusComment.Store("")

	newTask.start()

	actual, _ := s.tasks.LoadOrStore(userID, &sync.Map{})
	userTasks := actual.(*sync.Map)
	userTasks.Store(taskID, newTask)

	return nil
}

func saveSong(ctx context.Context, userID uuid.UUID, reader io.ReadCloser, extension string, data model.Song) error {
	defer reader.Close()

	// user, err := repository.GetUserByID(userID)
	// if err != nil {
	// 	return fmt.Errorf("saveSong: %w", err)
	// }
	//
	// root, err := os.OpenRoot(config.LIBRARY_PATH)
	// if err != nil {
	// 	return fmt.Errorf("saveSong: os.OpenRoot: %w", err)
	// }
	// defer root.Close()
	//
	// if err := root.Mkdir(user.Username, 0755); err != nil && !os.IsExist(err) {
	// 	return fmt.Errorf("saveSong: root.Mkdir: 1: %w", err)
	// }
	//
	// rootUser, err := root.OpenRoot(user.Username)
	// if err != nil {
	// 	return fmt.Errorf("saveSong: root.OpenRoot: %w", err)
	// }
	// defer rootUser.Close()
	//
	// artistName := strings.ReplaceAll(data.Artists[0].Name, "/", "_")
	// albumTitle := strings.ReplaceAll(data.Album.Title, "/", "_")
	// songTitle := strings.ReplaceAll(data.Title, "/", "_")
	//
	// filename := filepath.Join(artistName, albumTitle, fmt.Sprintf("%d - %s.%s", data.TrackNumber, songTitle, extension))
	// dirFile := filepath.Dir(filename)
	//
	// if err := rootUser.MkdirAll(dirFile, 0755); err != nil {
	// 	return fmt.Errorf("saveSong: rootUser.MkdirAll: %w", err)
	// }
	//
	// file, err := rootUser.Create(filename)
	// if err != nil {
	// 	return fmt.Errorf("saveSong: root.Create: %w", err)
	// }
	//
	// if _, err := io.Copy(file, reader); err != nil {
	// 	_ = file.Close()
	// 	if removeErr := rootUser.Remove(filename); removeErr != nil {
	// 		return fmt.Errorf("saveSong: io.Copy: %w: %w", err, removeErr)
	// 	} else {
	// 		return fmt.Errorf("saveSong: io.Copy: %w", err)
	// 	}
	// }
	// if err := file.Close(); err != nil {
	// 	return fmt.Errorf("saveSong: file.Close: %w", err)
	// }
	//
	// path := filepath.Join(rootUser.Name(), filename)
	// if err := metadata.FormatMetadata(ctx, userID, path, data); err != nil {
	// 	_ = file.Close()
	// 	if removeErr := rootUser.Remove(filename); removeErr != nil {
	// 		return fmt.Errorf("saveSong: %w: %w", err, removeErr)
	// 	} else {
	// 		return fmt.Errorf("saveSong: %w", err)
	// 	}
	// }
	//
	// _ = repository.AddSong(models.Song{UserID: userID, Path: filename, Isrc: data.Isrc, MTime: time.Now()})

	return nil
}

func (t *downloadTask) start() {
	if t.status.Load().(model.DownloadStatus) == model.DownloadStatusRunning {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.downloadCancel.Store(cancel)
	t.status.Store(model.DownloadStatusPending)
	t.statusComment.Store("")

	go t.run(ctx)
}

func (t *downloadTask) run(ctx context.Context) {
	if err := t.globalLimit.Acquire(ctx); err != nil {
		t.status.Store(model.DownloadStatusCancel)
		t.statusComment.Store(fmt.Sprintf("wait for download: %s", err.Error()))
		return
	}
	defer t.globalLimit.Release()

	t.status.Store(model.DownloadStatusRunning)

	reader, extension, err := t.plugin.DownloadWithUserID(ctx, t.userID, t.song.Provider, t.song.Id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			t.status.Store(model.DownloadStatusCancel)
		} else {
			t.status.Store(model.DownloadStatusFailed)
		}
		t.statusComment.Store(fmt.Sprintf("download: %s", err.Error()))
		return
	}

	err = saveSong(ctx, t.userID, reader, extension, t.song.Song)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			t.status.Store(model.DownloadStatusCancel)
		} else {
			t.status.Store(model.DownloadStatusFailed)
		}
		t.statusComment.Store(fmt.Sprintf("save song: %s", err.Error()))
		return
	}

	t.status.Store(model.DownloadStatusDone)
}

func (t *downloadTask) cancel() {
	if status := t.status.Load().(model.DownloadStatus); status != model.DownloadStatusRunning && status != model.DownloadStatusPending {
		return
	}

	cancel, ok := t.downloadCancel.Load().(context.CancelFunc)
	if !ok {
		return
	}

	cancel()
}

func (t *downloadTask) retry() {
	if status := t.status.Load().(model.DownloadStatus); status != model.DownloadStatusCancel && status != model.DownloadStatusFailed {
		return
	}

	t.start()
}

func (s *DownloadService) Retry(userID uuid.UUID, taskID uuid.UUID) error {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return model.ErrDownloadNotFound
	}
	userTasks := actual.(*sync.Map)

	actual, ok = userTasks.Load(taskID)
	if !ok {
		return model.ErrDownloadNotFound
	}
	task := actual.(*downloadTask)

	task.retry()
	return nil
}

func (s *DownloadService) RetryAll(userID uuid.UUID) {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return
	}
	userTasks := actual.(*sync.Map)

	userTasks.Range(func(key any, value any) bool {
		task := value.(*downloadTask)
		task.retry()
		return true
	})
}

func (s *DownloadService) Cancel(userID uuid.UUID, taskID uuid.UUID) error {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return model.ErrDownloadNotFound
	}
	userTasks := actual.(*sync.Map)

	actual, ok = userTasks.Load(taskID)
	if !ok {
		return model.ErrDownloadNotFound
	}
	task := actual.(*downloadTask)

	task.cancel()
	return nil
}

func (s *DownloadService) RemoveDone(userID uuid.UUID) {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return
	}
	userTasks := actual.(*sync.Map)

	doneList := []uuid.UUID{}
	userTasks.Range(func(key any, value any) bool {
		taskID := key.(uuid.UUID)
		task := value.(*downloadTask)
		if task.status.Load().(model.DownloadStatus) == model.DownloadStatusDone {
			doneList = append(doneList, taskID)
		}
		return true
	})

	for _, id := range doneList {
		userTasks.Delete(id)
	}
}

func (s *DownloadService) Remove(userID uuid.UUID, taskID uuid.UUID) {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return
	}
	userTasks := actual.(*sync.Map)

	userTasks.Delete(taskID)
}

func (s *DownloadService) List(userID uuid.UUID) []model.DownloadTaskInfo {
	actual, ok := s.tasks.Load(userID)
	if !ok {
		return []model.DownloadTaskInfo{}
	}
	userTasks := actual.(*sync.Map)

	tasks := []model.DownloadTaskInfo{}
	userTasks.Range(func(key any, value any) bool {
		taskID := key.(uuid.UUID)
		task := value.(*downloadTask)

		tasks = append(tasks, model.DownloadTaskInfo{
			ID:            taskID,
			UserID:        userID,
			Data:          task.song,
			Status:        task.status.Load().(model.DownloadStatus),
			StatusComment: task.statusComment.Load().(string),
		})

		return true
	})

	return tasks
}
