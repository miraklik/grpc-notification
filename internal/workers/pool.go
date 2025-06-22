package workers

import (
	"context"
	"log"
	"notification_service/internal/channel"
	"notification_service/internal/models"
	"notification_service/internal/queue"
	"notification_service/internal/service"
	"sync"
	"time"
)

type Worker struct {
	id       int
	queue    *queue.Queue
	services *service.NotificationsService
	channel  map[models.NotificationType]channel.NotificationChannel
}

type Pool struct {
	size     int
	queue    *queue.Queue
	services *service.NotificationsService
	channels map[models.NotificationType]channel.NotificationChannel
	workers  []*Worker
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorker(id int, queue *queue.Queue, services *service.NotificationsService, channels map[models.NotificationType]channel.NotificationChannel) *Worker {
	return &Worker{
		id:       id,
		queue:    queue,
		services: services,
		channel:  channels,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", w.id)
			return
		default:
			notification := w.queue.Pop()
			if notification == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			log.Printf("Worker %d received notification: %s", w.id, notification.Message)
			go w.proccessWorker(ctx, notification)
		}
	}
}

func (w *Worker) proccessWorker(ctx context.Context, notification *models.Notification) {
	notification.Status = models.StatusInProgress
	notification.UpdatedAt = time.Now()
	w.services.UpdateNote(notification.ID, notification)

	channel, exists := w.channel[notification.Type]
	if !exists {
		notification.Status = models.StatusFailed
		notification.LastError = "Invalid notification type: " + string(notification.Type)
		notification.UpdatedAt = time.Now()
		w.services.UpdateNote(notification.ID, notification)
	}

	err := channel.Send(ctx, notification)
	if err != nil {
		notification.LastError = err.Error()
		notification.UpdatedAt = time.Now()

		if notification.Attempts < 3 {
			notification.Status = models.StatusRetrying
			w.services.UpdateNote(notification.ID, notification)
			w.queue.Push(notification)
		} else {
			notification.Status = models.StatusFailed
			w.services.UpdateNote(notification.ID, notification)
		}
	}
	notification.Status = models.StatusSent
	notification.UpdatedAt = time.Now()
	w.services.UpdateNote(notification.ID, notification)
}

func NewPool(size int, queue *queue.Queue, services *service.NotificationsService) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	channelMap := make(map[models.NotificationType]channel.NotificationChannel)
	channelMap[models.TypeEmail] = channel.NewEmailChannel()
	channelMap[models.TypePush] = channel.NewPushChannel()
	channelMap[models.TypeSMS] = channel.NewSmsChannel()

	return &Pool{
		size:     size,
		queue:    queue,
		services: services,
		channels: channelMap,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (p *Pool) Start() {
	p.workers = make([]*Worker, p.size)

	for i := 0; i < p.size; i++ {
		worker := &Worker{
			id:       i,
			queue:    p.queue,
			services: p.services,
			channel:  p.channels,
		}

		p.workers[i] = worker
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Run(p.ctx)
		}(worker)
	}

	log.Printf("Started %d worker", p.size)
}

func (p *Pool) Stop() {
	p.cancel()
	p.queue.Close()
	p.wg.Wait()
	log.Println("Pool stopped")
}
