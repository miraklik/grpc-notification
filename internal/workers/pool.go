package workers

import (
	"context"
	"encoding/json"
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
	queue    *queue.NatsQueue
	services *service.NotificationsService
	channel  map[models.NotificationType]channel.NotificationChannel
}

type Pool struct {
	size     int
	queue    *queue.NatsQueue
	services *service.NotificationsService
	channels map[models.NotificationType]channel.NotificationChannel
	workers  []*Worker
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorker(id int, queue *queue.NatsQueue, services *service.NotificationsService, channels map[models.NotificationType]channel.NotificationChannel) *Worker {
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
			msg, err := w.queue.Pop(ctx)
			if err != nil {
				log.Printf("Error popping message: %v", err)
				continue
			}

			if msg == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			var n models.Notification
			if err := json.Unmarshal(msg.Data, &n); err != nil {
				_ = msg.Nak()
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			if err := w.proccessWorker(ctx, &n); err != nil {
				_ = msg.Nak()
				log.Printf("Error processing worker: %v", err)
				continue
			} else {
				_ = msg.Ack()
			}
		}
	}
}

func (w *Worker) proccessWorker(ctx context.Context, notification *models.Notification) error {
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

	if err := channel.Send(ctx, notification); err != nil {
		notification.LastError = err.Error()
		notification.UpdatedAt = time.Now()
		w.services.UpdateNote(notification.ID, notification)
		return err
	}

	notification.Status = models.StatusSent
	notification.UpdatedAt = time.Now()
	w.services.UpdateNote(notification.ID, notification)

	return nil
}

func NewPool(size int, queue *queue.NatsQueue, services *service.NotificationsService) *Pool {
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
	p.wg.Wait()
	log.Println("Pool stopped")
}
