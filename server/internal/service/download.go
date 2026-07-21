package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	pkg_sync "github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type DownloadService struct {
	cfg            config.DownloadConfig
	tasks          sync.Map // map[uuid.UUID]sync.Map -> map[uuid.UUID]*downloadTask
	limit          pkg_sync.Semaphore
	plugin         *PluginService
	userRepository model.UserRepository
}

type downloadTask struct {
	user           model.User
	song           model.EnrichedSong
	running        atomic.Bool
	status         atomic.Value // model.DownloadStatus
	statusComment  atomic.Value // string
	downloadCancel atomic.Value // context.CancelFunc
	globalLimit    *pkg_sync.Semaphore
	cfg            config.DownloadConfig
	plugin         *PluginService
	userRepository model.UserRepository
}

func NewDownloadService(cfg config.DownloadConfig, plugin *PluginService, instance model.InstanceRepository, userRepository model.UserRepository) *DownloadService {
	return &DownloadService{
		cfg:            cfg,
		limit:          pkg_sync.NewSemaphore(cfg.Concurrent),
		plugin:         plugin,
		userRepository: userRepository,
	}
}

func (s *DownloadService) AddArtist(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, artistID string) ([]model.AddDownloadError, error) {
	artist, err := s.plugin.GetArtistFromPluginInstances(ctx, pluginInstances, artistID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get artist %s info for download: %w", artistID, err)
	}

	errList := []model.AddDownloadError{}
	for _, album := range artist.Albums {
		subErrList, err := s.AddAlbum(ctx, user, pluginInstances, album.Id)
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

func (s *DownloadService) AddAlbum(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, albumID string) ([]model.AddDownloadError, error) {
	album, err := s.plugin.GetAlbumFromPluginInstances(ctx, pluginInstances, albumID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get album %s info for download: %w", albumID, err)
	}

	errList := []model.AddDownloadError{}
	for _, song := range album.Songs {
		err := s.AddSong(ctx, user, pluginInstances, song.Id)
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

func (s *DownloadService) AddPlaylist(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, playlistID string) ([]model.AddDownloadError, error) {
	playlist, err := s.plugin.GetPlaylistFromPluginInstances(ctx, pluginInstances, playlistID)
	if err != nil {
		return []model.AddDownloadError{}, fmt.Errorf("get playlist %s info for download: %w", playlistID, err)
	}

	errList := []model.AddDownloadError{}
	for _, song := range playlist.Songs {
		err := s.AddSong(ctx, user, pluginInstances, song.Id)
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

func (s *DownloadService) AddSong(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, songID string) error {
	song, err := s.plugin.GetSongFromPluginInstances(ctx, pluginInstances, songID)
	if err != nil {
		return fmt.Errorf("get song %s info for download: %w", songID, err)
	}

	if err := s.DownloadSong(user, song); err != nil {
		return fmt.Errorf("download song %s: %w", song.Title, err)
	}

	return nil
}

func (s *DownloadService) DownloadSong(user model.User, song model.EnrichedSong) error {
	// check if song already owned

	taskID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("create id for new download task: %w", err)
	}

	newTask := &downloadTask{
		user:           user,
		song:           song,
		running:        atomic.Bool{},
		status:         atomic.Value{},
		statusComment:  atomic.Value{},
		downloadCancel: atomic.Value{},
		globalLimit:    &s.limit,
		cfg:            s.cfg,
		plugin:         s.plugin,
		userRepository: s.userRepository,
	}

	newTask.status.Store(model.DownloadStatusPending)
	newTask.statusComment.Store("")

	newTask.start()

	actual, _ := s.tasks.LoadOrStore(user.ID, &sync.Map{})
	userTasks := actual.(*sync.Map)
	userTasks.Store(taskID, newTask)

	return nil
}

func saveSong(ctx context.Context, downloadPath string, user model.User, song model.Song, reader io.ReadCloser, extension string) error {
	defer reader.Close()

	root, err := os.OpenRoot(downloadPath)
	if err != nil {
		return fmt.Errorf("open download folder: %w", err)
	}
	defer root.Close()

	rootUser, err := root.OpenRoot(user.Username)
	if err != nil {
		return fmt.Errorf("open user folder: %w", err)
	}
	defer rootUser.Close()

	artistName := strings.ReplaceAll(song.Artists[0].Name, "/", "_")
	albumTitle := strings.ReplaceAll(song.Album.Title, "/", "_")
	songTitle := strings.ReplaceAll(song.Title, "/", "_")

	dirFile := filepath.Join(artistName, albumTitle)
	filename := filepath.Join(dirFile, fmt.Sprintf("%d - %s.%s", song.TrackNumber, songTitle, extension))

	if err := rootUser.MkdirAll(dirFile, 0755); err != nil {
		return fmt.Errorf("create song folders: %w", err)
	}

	file, err := rootUser.Create(filename)
	if err != nil {
		return fmt.Errorf("create song file: %w", err)
	}
	defer func() {
		file.Close()
		if err != nil {
			path := filepath.Join(rootUser.Name(), filename)
			os.Remove(path)
		}
	}()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("copy song content into song file: %w", err)
	}

	// if err := metadata.FormatMetadata(ctx, userID, path, data); err != nil {
	// 		return fmt.Errorf("format song metadata: %w", err)
	// }

	// _ = repository.AddSong(models.Song{UserID: userID, Path: filename, Isrc: data.Isrc, MTime: time.Now()})

	return nil
}

func (t *downloadTask) start() {
	if ok := t.running.CompareAndSwap(false, true); !ok {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.downloadCancel.Store(cancel)

	if user, err := t.userRepository.GetUserByFilter(ctx, model.UserFilter{ID: &t.user.ID}); err != nil {
		t.status.Store(model.DownloadStatusFailed)
		t.statusComment.Store(fmt.Sprintf("refresh user info for download: %s", err.Error()))
		t.running.CompareAndSwap(true, false)
		return
	} else {
		t.user = user
	}

	t.status.Store(model.DownloadStatusPending)
	t.statusComment.Store("")

	go t.run(ctx)
}

func (t *downloadTask) run(ctx context.Context) {
	if err := t.globalLimit.Acquire(ctx); err != nil {
		t.running.CompareAndSwap(true, false)
		t.status.Store(model.DownloadStatusCancel)
		t.statusComment.Store(fmt.Sprintf("wait for download: %s", err.Error()))
		return
	}
	defer t.globalLimit.Release()

	t.status.Store(model.DownloadStatusRunning)

	reader, extension, err := t.plugin.DownloadWithUserID(ctx, t.user.ID, t.song.Provider, t.song.Id)
	if err != nil {
		t.running.CompareAndSwap(true, false)
		if errors.Is(err, context.Canceled) {
			t.status.Store(model.DownloadStatusCancel)
		} else {
			t.status.Store(model.DownloadStatusFailed)
		}
		t.statusComment.Store(fmt.Sprintf("download: %s", err.Error()))
		return
	}

	err = saveSong(ctx, t.cfg.Path, t.user, t.song.Song, reader, extension)
	if err != nil {
		t.running.CompareAndSwap(true, false)
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
