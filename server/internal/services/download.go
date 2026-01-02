package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
)

type downloadManager struct {
	mTasks sync.Mutex
	tasks  map[uint]map[uint]*downloadTask
	mIds   sync.Mutex
	ids    map[uint]uint
}

type downloadTask struct {
	mu             sync.Mutex
	userId         uint
	api            models.Plugin
	songId         string
	quality        string
	songData       models.SongData
	status         models.Status
	downloadCancel context.CancelFunc
}

var DownloadManager = downloadManager{
	mTasks: sync.Mutex{},
	tasks:  make(map[uint]map[uint]*downloadTask),
	mIds:   sync.Mutex{},
	ids:    make(map[uint]uint),
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

func (m *downloadManager) AddArtist(userId uint, api models.Plugin, artistId string, quality string) {
	artist, err := api.Artist(context.Background(), userId, artistId)
	if err != nil {
		fmt.Printf("downloadManager.AddArtist: %v\n", err)
		return
	}
	for _, album := range artist.Albums {
		m.AddAlbum(userId, api, album.Id, quality)
	}
}

func (m *downloadManager) AddAlbum(userId uint, api models.Plugin, albumId string, quality string) {
	album, err := api.Album(context.Background(), userId, albumId)
	if err != nil {
		fmt.Printf("downloadManager.AddAlbum: %v\n", err)
		return
	}
	for _, song := range album.Songs {
		m.AddSong(userId, api, song.Id, quality)
	}
}

func (m *downloadManager) AddSong(userId uint, api models.Plugin, songId string, quality string) {
	if quality == "" {
		user, err := repository.GetUserByID(userId)
		if err != nil {
			fmt.Printf("downloadManager.AddSong: %v\n", err)
		} else if !user.BestQuality {
			quality = "HIGH"
		}
	}

	taskId := m.generateId(userId)

	newTask := &downloadTask{
		mu:             sync.Mutex{},
		userId:         userId,
		api:            api,
		songId:         songId,
		quality:        quality,
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

func saveSong(userId uint, reader io.ReadCloser, extension string, data models.SongData) error {
	defer reader.Close()

	user, err := repository.GetUserByID(userId)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}

	root, err := os.OpenRoot(config.DOWNLOAD_FOLDER)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer root.Close()

	if err := root.Mkdir(user.Username, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	root, err = os.OpenRoot(filepath.Join(config.DOWNLOAD_FOLDER, user.Username))
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer root.Close()

	artistName := strings.ReplaceAll(data.Artists[0].Name, "/", "_")
	albumTitle := strings.ReplaceAll(data.Album.Title, "/", "_")
	songTitle := strings.ReplaceAll(data.Title, "/", "_")

	filename := filepath.Join(artistName, albumTitle, fmt.Sprintf("%d - %s.%s", data.TrackNumber, songTitle, extension))

	if err := root.Mkdir(data.Artists[0].Name, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	if err := root.Mkdir(filepath.Join(data.Artists[0].Name, data.Album.Title), 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	file, err := root.Create(filename)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		path := file.Name()
		file.Close()
		os.Remove(path)
		return fmt.Errorf("saveSong: %w", err)
	}

	if err := metadata.FormatMetadata(userId, filepath.Join(root.Name(), filename), data); err != nil {
		return fmt.Errorf("saveSong: %w", err)
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
	t.status = models.StatusRunning
	t.mu.Unlock()

	go t.run(ctx)
}

func (t *downloadTask) run(ctx context.Context) {
	go func() {
		song, err := t.api.Song(ctx, t.userId, t.songId)
		t.mu.Lock()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				t.status = models.StatusCancel
			} else {
				t.cancel()
				t.status = models.StatusFailed
				fmt.Println("downloadTask.run: ", err)
			}
		} else {
			t.songData = song
		}
		t.mu.Unlock()
	}()

	reader, extension, err := t.api.Download(ctx, t.userId, t.songId, t.quality)

	if err == nil {
		err = saveSong(t.userId, reader, extension, t.songData)
	}

	t.mu.Lock()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			t.status = models.StatusCancel
		} else {
			t.cancel()
			t.status = models.StatusFailed
			fmt.Println("downloadTask.run: ", err)
		}
	} else {
		t.status = models.StatusDone
	}
	t.mu.Unlock()
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
	m.mTasks.Lock()
	var doneList []uint
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
			Id:     id,
			Data:   task.songData,
			Api:    task.api.Name(),
			Status: task.status,
		}
		task.mu.Unlock()
		tasks = append(tasks, tmp)
	}
	m.mTasks.Unlock()
	return tasks
}
