package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
)

type downloadManager struct {
	mTasks      sync.Mutex
	tasks       map[uint]map[uint]*downloadTask
	mIds        sync.Mutex
	ids         map[uint]uint
	limitGlobal chan struct{}
}

type downloadTask struct {
	mu             sync.Mutex
	userId         uint
	provider       string
	songId         string
	songData       models.SongData
	status         models.Status
	downloadCancel context.CancelFunc
}

var DownloadManager = downloadManager{
	mTasks:      sync.Mutex{},
	tasks:       make(map[uint]map[uint]*downloadTask),
	mIds:        sync.Mutex{},
	ids:         make(map[uint]uint),
	limitGlobal: make(chan struct{}, 3),
}

func (m *downloadManager) acquireLimitGlobal(ctx context.Context) error {
	select {
	case m.limitGlobal <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *downloadManager) releaseLimitGlobal() {
	<-m.limitGlobal
}

func (m *downloadManager) generateId(userId uint) uint {
	m.mIds.Lock()
	id, ok := m.ids[userId]
	if !ok {
		id = 0
	}
	id++
	m.ids[userId] = id
	m.mIds.Unlock()

	return id
}

func (m *downloadManager) AddArtist(userId uint, provider string, artistId string) {
	artist, err := plugins.GetArtist(context.Background(), userId, provider, artistId)
	if err != nil {
		log.Println("downloadManager.AddArtist: ", err)
		return
	}
	for _, album := range artist.Albums {
		m.AddAlbum(userId, provider, album.Id)
	}
}

func (m *downloadManager) AddAlbum(userId uint, provider string, albumId string) {
	album, err := plugins.GetAlbum(context.Background(), userId, provider, albumId)
	if err != nil {
		log.Println("downloadManager.AddAlbum: ", err)
		return
	}
	for _, song := range album.Songs {
		m.AddSong(userId, provider, song.Id)
	}
}

func (m *downloadManager) AddSong(userId uint, provider string, songId string) {
	taskId := m.generateId(userId)

	newTask := &downloadTask{
		mu:             sync.Mutex{},
		userId:         userId,
		provider:       provider,
		songId:         songId,
		songData:       models.SongData{},
		status:         models.StatusPending,
		downloadCancel: nil,
	}

	m.mTasks.Lock()
	userTasks, ok := m.tasks[userId]
	if !ok {
		userTasks = make(map[uint]*downloadTask)
		m.tasks[userId] = userTasks
	}
	userTasks[taskId] = newTask
	m.mTasks.Unlock()

	go newTask.start()
}

func saveSong(ctx context.Context, userId uint, reader io.ReadCloser, extension string, data models.SongData) error {
	defer reader.Close()

	user, err := repository.GetUserByID(userId)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}

	root, err := os.OpenRoot(config.LIBRARY_PATH)
	if err != nil {
		return fmt.Errorf("saveSong: os.OpenRoot: 1: %w", err)
	}
	defer root.Close()

	if err := root.Mkdir(user.Username, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: root.Mkdir: 1: %w", err)
	}

	root, err = os.OpenRoot(filepath.Join(config.LIBRARY_PATH, user.Username))
	if err != nil {
		return fmt.Errorf("saveSong: os.OpenRoot: 2: %w", err)
	}
	defer root.Close()

	artistName := strings.ReplaceAll(data.Artists[0].Name, "/", "_")
	albumTitle := strings.ReplaceAll(data.Album.Title, "/", "_")
	songTitle := strings.ReplaceAll(data.Title, "/", "_")

	filename := filepath.Join(artistName, albumTitle, fmt.Sprintf("%d - %s.%s", data.TrackNumber, songTitle, extension))

	if err := root.Mkdir(data.Artists[0].Name, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: root.Mkdir: 2: %w", err)
	}

	if err := root.Mkdir(filepath.Join(data.Artists[0].Name, data.Album.Title), 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: root.Mkdir: 3: %w", err)
	}

	file, err := root.Create(filename)
	if err != nil {
		return fmt.Errorf("saveSong: root.Create: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		filepath := file.Name()
		_ = file.Close()
		if removeErr := os.Remove(filepath); removeErr != nil {
			return fmt.Errorf("saveSong: io.Copy: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("saveSong: io.Copy: %w", err)
		}
	}

	if err := metadata.FormatMetadata(ctx, userId, filepath.Join(root.Name(), filename), data); err != nil {
		filepath := file.Name()
		_ = file.Close()
		if removeErr := os.Remove(filepath); removeErr != nil {
			return fmt.Errorf("saveSong: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("saveSong: %w", err)
		}
	}

	return nil
}

func (t *downloadTask) start() {
	t.mu.Lock()
	if t.status == models.StatusRunning {
		t.mu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.downloadCancel = cancel
	t.mu.Unlock()

	go t.run(ctx)
}

func (t *downloadTask) run(ctx context.Context) {
	song, err := plugins.GetSong(ctx, t.userId, t.provider, t.songId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			t.mu.Lock()
			t.status = models.StatusCancel
			t.mu.Unlock()
		} else {
			t.mu.Lock()
			t.status = models.StatusFailed
			t.mu.Unlock()
		}
		log.Println("downloadTask.run: ", err)
		return
	} else {
		t.mu.Lock()
		t.songData = song
		t.mu.Unlock()
	}

	// check if the user already own this song
	// if false {
	//    return
	// }

	if err := DownloadManager.acquireLimitGlobal(ctx); err != nil {
		t.mu.Lock()
		t.status = models.StatusCancel
		t.mu.Unlock()
		return
	} else {
		t.mu.Lock()
		t.status = models.StatusRunning
		t.mu.Unlock()
	}
	defer DownloadManager.releaseLimitGlobal()

	reader, extension, err := plugins.Download(ctx, t.userId, t.provider, t.songId)

	if err == nil {
		err = saveSong(ctx, t.userId, reader, extension, t.songData)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			t.mu.Lock()
			if t.status != models.StatusFailed {
				t.status = models.StatusCancel
			}
			t.mu.Unlock()
		} else {
			t.downloadCancel()
			t.mu.Lock()
			t.status = models.StatusFailed
			t.mu.Unlock()
		}
		log.Println("downloadTask.run: ", err)
	} else {
		t.mu.Lock()
		t.status = models.StatusDone
		t.mu.Unlock()
	}
}

func (t *downloadTask) cancel() {
	t.mu.Lock()
	if t.status == models.StatusRunning {
		t.downloadCancel()
		t.status = models.StatusCancel
	}
	t.mu.Unlock()
}

func (t *downloadTask) retry() {
	t.mu.Lock()
	if t.status != models.StatusCancel && t.status != models.StatusFailed {
		t.mu.Unlock()
		return
	}
	t.status = models.StatusPending
	t.mu.Unlock()

	t.start()
}

func (m *downloadManager) Retry(userId uint, taskId uint) error {
	m.mTasks.Lock()
	task, ok := m.tasks[userId][taskId]
	m.mTasks.Unlock()
	if !ok {
		return errors.New("download not found")
	}
	task.retry()
	return nil
}

func (m *downloadManager) RetryAll(userId uint) {
	m.mTasks.Lock()
	for _, task := range m.tasks[userId] {
		task.retry()
	}
	m.mTasks.Unlock()
}

func (m *downloadManager) Cancel(userId uint, taskId uint) error {
	m.mTasks.Lock()
	task, ok := m.tasks[userId][taskId]
	m.mTasks.Unlock()
	if !ok {
		return errors.New("download not found")
	}
	task.cancel()
	return nil
}

func (m *downloadManager) Done(userId uint) {
	var doneList []uint
	m.mTasks.Lock()
	for id, task := range m.tasks[userId] {
		task.mu.Lock()
		if task.status == models.StatusDone {
			doneList = append(doneList, id)
		}
		task.mu.Unlock()
	}
	for _, id := range doneList {
		delete(m.tasks[userId], id)
	}
	m.mTasks.Unlock()
}

func (m *downloadManager) Remove(userId uint, taskId uint) error {
	m.mTasks.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		m.mTasks.Unlock()
		return errors.New("download not found")
	}
	task.cancel()
	delete(m.tasks[userId], taskId)
	m.mTasks.Unlock()
	return nil
}

func (m *downloadManager) List(userId uint) []models.DownloadData {
	m.mTasks.Lock()
	tasks := make([]models.DownloadData, 0, len(m.tasks))
	for id, task := range m.tasks[userId] {
		task.mu.Lock()
		tmp := models.DownloadData{
			Id:       id,
			Data:     task.songData,
			Provider: task.provider,
			Status:   task.status,
		}
		task.mu.Unlock()
		tasks = append(tasks, tmp)
	}
	m.mTasks.Unlock()
	return tasks
}
