package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

type downloadManager struct {
	mu    sync.Mutex
	tasks map[uint]map[uint]*downloadTask
	ids   map[uint]uint
}

type downloadTask struct {
	id             uint
	userId         uint
	api            models.Plugin
	songId         string
	quality        string
	songData       models.SongData
	status         models.Status // depend-on a chan to be change
	retryDownload  chan struct{}
	cancelDownload chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
}

var DownloadManager = downloadManager{
	mu:    sync.Mutex{},
	tasks: make(map[uint]map[uint]*downloadTask),
	ids:   make(map[uint]uint),
}

func (m *downloadManager) generateId(userId uint) uint {
	id, ok := m.ids[userId]
	if !ok {
		id = 0
	}
	id++
	m.ids[userId] = id

	return id
}

func (m *downloadManager) Add(userId uint, api models.Plugin, songId string, quality string) {
	taskId := m.generateId(userId)
	ctx, cancel := context.WithCancel(context.Background())

	newTask := &downloadTask{
		id:             taskId,
		userId:         userId,
		api:            api,
		songId:         songId,
		quality:        quality,
		status:         models.StatusPending,
		retryDownload:  make(chan struct{}),
		cancelDownload: make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}

	m.mu.Lock()
	userTasks, ok := m.tasks[userId]
	if !ok {
		userTasks = make(map[uint]*downloadTask)
		m.tasks[userId] = userTasks
	}
	userTasks[taskId] = newTask
	m.mu.Unlock()

	go m.startMaster(newTask)
}

func (m *downloadManager) startMaster(t *downloadTask) {
	ctx, cancel := context.WithCancel(context.Background())
	status := make(chan models.Status)
	t.status = models.StatusPending

	go m.run(t, ctx, status)

	for {
		select {
		case res := <-status:
			t.status = res
			if res == models.StatusDone {
				cancel()
				return
			}
		case <-t.retryDownload:
			if t.status == models.StatusCancel || t.status == models.StatusFailed {
				ctx, cancel = context.WithCancel(context.Background())
				status := make(chan models.Status)
				t.status = models.StatusPending

				go m.run(t, ctx, status)
			}
		case <-t.cancelDownload:
			cancel()
		case <-t.ctx.Done():
			cancel()
			return
		}
	}
}

func (m *downloadManager) run(t *downloadTask, ctx context.Context, status chan models.Status) {
	if t.songData.Id == "" {
		song, err := t.api.Song(ctx, t.songId)
		select {
		case <-ctx.Done():
			status <- models.StatusCancel
			return
		default:
		}
		if err != nil {
			status <- models.StatusFailed
			return
		}
		t.songData = song
	}

	status <- models.StatusRunning

	req, _ := http.NewRequestWithContext(ctx, "GET", t.songData.DownloadUrl, nil)
	resp, err := http.DefaultClient.Do(req)
	select {
	case <-ctx.Done():
		status <- models.StatusCancel
		return
	default:
	}
	if err != nil {
		status <- models.StatusFailed
		return
	}
	defer resp.Body.Close()
	filename := fmt.Sprintf("%d - %s.flac", t.songData.TrackNumber, t.songData.Title)
	out, err := os.Create(filename)
	if err != nil {
		status <- models.StatusFailed
		return
	}
	defer out.Close()

	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			out.Close()
			os.Remove(filename)
			status <- models.StatusCancel
			return
		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				out.Write(buf[:n])
			}
			if err == io.EOF {
				status <- models.StatusDone
				return
			}
			if err != nil {
				status <- models.StatusFailed
				return
			}
		}
	}
}

func (m *downloadManager) Retry(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	task.retryDownload <- struct{}{}
	m.mu.Unlock()
	return nil
}

func (m *downloadManager) Cancel(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	task.cancelDownload <- struct{}{}
	m.mu.Unlock()
	return nil
}

func (m *downloadManager) Remove(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	task.cancel()
	delete(m.tasks[userId], taskId)
	m.mu.Unlock()
	return nil
}

func (m *downloadManager) List(userId uint) []models.DownloadData {
	m.mu.Lock()
	tasks := make([]models.DownloadData, 0, len(m.tasks))
	for _, value := range m.tasks[userId] {
		tmp := models.DownloadData{
			Id:     value.id,
			Data:   value.songData,
			Api:    value.api.Name(),
			Status: value.status,
		}
		tasks = append(tasks, tmp)
	}
	m.mu.Unlock()
	return tasks
}
