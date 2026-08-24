package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	pkg_path "github.com/DimitriLaPoudre/MusicShack/internal/pkg/path"
	pkg_sync "github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type DownloadService struct {
	cfgDownload    config.DownloadConfig
	cfgPlugin      config.PluginConfig
	tasks          sync.Map // map[uuid.UUID]sync.Map -> map[uuid.UUID]*downloadTask
	limit          pkg_sync.Semaphore
	plugin         *PluginService
	metadata       *MetadataService
	instance       model.InstanceRepository
	userRepository model.UserRepository
}

type downloadTask struct {
	user           model.User
	songInfo       model.EnrichedSongInfo
	running        atomic.Bool
	status         atomic.Value // model.DownloadStatus
	statusComment  atomic.Value // string
	downloadCancel atomic.Value // context.CancelFunc
	download       *DownloadService
}

func NewDownloadService(cfgDownload config.DownloadConfig, cfgPlugin config.PluginConfig, plugin *PluginService, metadata *MetadataService, instance model.InstanceRepository, userRepository model.UserRepository) *DownloadService {
	return &DownloadService{
		cfgDownload:    cfgDownload,
		cfgPlugin:      cfgPlugin,
		limit:          pkg_sync.NewSemaphore(cfgDownload.Concurrent),
		plugin:         plugin,
		metadata:       metadata,
		instance:       instance,
		userRepository: userRepository,
	}
}

func (s *DownloadService) DownloadArtist(ctx context.Context, user model.User, provider string, id string) error {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := s.plugin.InstancesToMapPluginInstances(instances)

	err = s.AddArtist(ctx, user, pluginInstances, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DownloadService) DownloadAlbum(ctx context.Context, user model.User, provider string, id string) error {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := s.plugin.InstancesToMapPluginInstances(instances)

	err = s.AddAlbum(ctx, user, pluginInstances, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DownloadService) DownloadPlaylist(ctx context.Context, user model.User, provider string, id string) error {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := s.plugin.InstancesToMapPluginInstances(instances)

	err = s.AddPlaylist(ctx, user, pluginInstances, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DownloadService) DownloadSong(ctx context.Context, user model.User, provider string, id string) error {
	instances, err := s.instance.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &user.ID, Provider: &provider})
	if err != nil {
		return fmt.Errorf("list instances of user %s for provider %s: %w", user.ID.String(), provider, err)
	}

	pluginInstances := s.plugin.InstancesToMapPluginInstances(instances)

	err = s.AddSong(ctx, user, pluginInstances, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DownloadService) AddArtist(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, artistID string) error {
	artistAlbums, err := s.plugin.GetArtistAlbumsFromPluginInstances(ctx, user, pluginInstances, artistID, s.cfgPlugin.Pagination.Limit, s.cfgPlugin.Pagination.Offset)
	if err != nil {
		return fmt.Errorf("get artist %s info for download: %w", artistID, err)
	}

	var allAlbums []model.EnrichedAlbumInfo
	allAlbums = append(allAlbums, artistAlbums.Albums.Albums...)
	allAlbums = append(allAlbums, artistAlbums.EPs.Albums...)
	allAlbums = append(allAlbums, artistAlbums.Singles.Albums...)

	for _, album := range allAlbums {
		_ = s.AddAlbum(ctx, user, pluginInstances, album.ID)
	}

	return nil
}

func (s *DownloadService) AddAlbum(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, albumID string) error {
	albumSongs, err := s.plugin.GetAlbumSongsFromPluginInstances(ctx, user, pluginInstances, albumID, s.cfgPlugin.Pagination.Limit, s.cfgPlugin.Pagination.Offset)
	if err != nil {
		return fmt.Errorf("get album %s info for download: %w", albumID, err)
	}

	for _, song := range albumSongs.Songs {
		_ = s.AddSong(ctx, user, pluginInstances, song.ID)
	}

	return nil
}

func (s *DownloadService) AddPlaylist(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, playlistID string) error {
	playlistSongs, err := s.plugin.GetPlaylistSongsFromPluginInstances(ctx, user, pluginInstances, playlistID, s.cfgPlugin.Pagination.Limit, s.cfgPlugin.Pagination.Offset)
	if err != nil {
		return fmt.Errorf("get playlist %s info for download: %w", playlistID, err)
	}

	for _, song := range playlistSongs.Songs {
		_ = s.AddSong(ctx, user, pluginInstances, song.ID)
	}

	return nil
}

func (s *DownloadService) AddSong(ctx context.Context, user model.User, pluginInstances map[model.Plugin][]model.Instance, songID string) error {
	songInfo, err := s.plugin.GetSongInfoFromPluginInstances(ctx, user, pluginInstances, songID)
	if err != nil {
		return fmt.Errorf("get song %s info for download: %w", songID, err)
	}

	albumInfo, err := s.plugin.GetAlbumInfoFromPluginInstances(ctx, user, pluginInstances, songInfo.Album.ID)
	if err != nil {
		return fmt.Errorf("get song %s album info for download: %w", songID, err)
	}

	songInfo.Album = albumInfo.AlbumInfo

	if err := s.AddDownloadTask(user, songInfo); err != nil {
		return fmt.Errorf("download song %s: %w", songInfo.Title, err)
	}

	return nil
}

func (s *DownloadService) AddDownloadTask(user model.User, songInfo model.EnrichedSongInfo) error {
	// check if song already owned

	taskID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("create id for new download task: %w", err)
	}

	newTask := &downloadTask{
		user:           user,
		songInfo:       songInfo,
		running:        atomic.Bool{},
		status:         atomic.Value{},
		statusComment:  atomic.Value{},
		downloadCancel: atomic.Value{},
		download:       s,
	}

	newTask.status.Store(model.DownloadStatusPending)
	newTask.statusComment.Store("")

	newTask.start()

	actual, _ := s.tasks.LoadOrStore(newTask.user.ID, &sync.Map{})
	userTasks := actual.(*sync.Map)
	userTasks.Store(taskID, newTask)

	return nil
}

func (t *downloadTask) start() {
	if ok := t.running.CompareAndSwap(false, true); !ok {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.downloadCancel.Store(cancel)

	if user, err := t.download.userRepository.GetUserByFilter(ctx, model.UserFilter{ID: &t.user.ID}); err != nil {
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
	if err := t.download.limit.Acquire(ctx); err != nil {
		t.running.CompareAndSwap(true, false)
		t.status.Store(model.DownloadStatusCancel)
		t.statusComment.Store(fmt.Sprintf("wait for download: %s", err.Error()))
		return
	}
	defer t.download.limit.Release()

	t.status.Store(model.DownloadStatusRunning)

	reader, extension, err := t.download.plugin.Download(ctx, t.user, t.songInfo.Provider, t.songInfo.ID, t.user.HiRes)
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

	err = t.SaveDownload(ctx, reader, extension)
	if err != nil {
		t.running.CompareAndSwap(true, false)
		if errors.Is(err, context.Canceled) {
			t.status.Store(model.DownloadStatusCancel)
		} else {
			t.status.Store(model.DownloadStatusFailed)
		}
		t.statusComment.Store(fmt.Sprintf("save downloaded song: %s", err.Error()))
		return
	}

	t.status.Store(model.DownloadStatusDone)
}

func (t *downloadTask) SaveDownload(ctx context.Context, reader io.ReadCloser, extension string) error {
	path, err := t.download.SaveSong(ctx, t.user, t.songInfo.SongInfo, reader, extension)
	if err != nil {
		return fmt.Errorf("save song: %w", err)
	}

	if err := t.download.metadata.FormatMetadata(ctx, t.user, t.songInfo.Provider, path, t.songInfo.SongInfo); err != nil {
		os.Remove(path)
		return fmt.Errorf("format song metadata: %w", err)
	}

	// _ = repository.AddSong(models.Song{UserID: userID, Path: filename, ISRC: data.ISRC, MTime: time.Now()})

	return nil
}

func (s *DownloadService) SaveSong(ctx context.Context, user model.User, songInfo model.SongInfo, reader io.ReadCloser, extension string) (string, error) {
	defer reader.Close()

	root, err := os.OpenRoot(s.cfgDownload.Path)
	if err != nil {
		return "", fmt.Errorf("open download folder: %w", err)
	}
	defer root.Close()

	rootUser, err := root.OpenRoot(user.ID.String())
	if err != nil {
		return "", fmt.Errorf("open user folder: %w", err)
	}
	defer rootUser.Close()

	artistName := "UnknownArtist"
	if len(songInfo.Album.Artists) > 0 {
		artistName = pkg_path.SanitizeName(songInfo.Album.Artists[0].Name)
	}
	albumTitle := pkg_path.SanitizeName(songInfo.Album.Title)
	songTitle := pkg_path.SanitizeName(songInfo.Title)

	dirFile := filepath.Join(artistName, albumTitle)
	filename := filepath.Join(dirFile, fmt.Sprintf("%d - %s.%s", songInfo.TrackNumber, songTitle, extension))

	if err := rootUser.MkdirAll(dirFile, 0755); err != nil {
		return "", fmt.Errorf("create song folders: %w", err)
	}

	path := filepath.Join(rootUser.Name(), filename)
	file, err := rootUser.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("create song file: %w", err)
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("copy song content into song file: %w", err)
	}

	return path, nil
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
			Data:          task.songInfo,
			Status:        task.status.Load().(model.DownloadStatus),
			StatusComment: task.statusComment.Load().(string),
		})

		return true
	})

	return tasks
}
