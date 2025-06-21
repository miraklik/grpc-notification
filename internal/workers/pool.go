package workers

import (
	"context"
	"log"
	"notification_service/internal/models"
	"notification_service/internal/queue"
	"notification_service/internal/service"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	id      int
	queue   *queue.Queue
	storage *service.NotificationsService
}

type Pool struct {
	size    int
	queue   *queue.Queue
	workers []*Worker
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWorker(id int, queue *queue.Queue, storage *service.NotificationsService) *Worker {
	return &Worker{
		id:      id,
		queue:   queue,
		storage: storage,
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
			if notification != nil {
				log.Printf("Worker %d received notification: %s", w.id, notification.Message)
			}
		}
	}
}

func (w *Worker) proccessWorker(ctx context.Context, notification *models.Notification) {
	notification.Status = models.StatusInProgress
	notification.UpdatedAt = time.Now()
	notificationID, err := strconv.Atoi(notification.ID)
	if err != nil {
		log.Println("Error converting notification ID to int: ", err)
		return
	}
	w.storage.UpdateNote(notificationID, notification)

}

func NewPool(size int, queue *queue.Queue) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		size:   size,
		queue:  queue,
		wg:     sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Pool) Start() {
	p.workers = make([]*Worker, p.size)

	for i := 0; i < p.size; i++ {
		worker := &Worker{}

		p.workers[i] = worker
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			worker.Run(p.ctx)
		}()
	}
}

func (p *Pool) Stop() {
	p.cancel()
	p.queue.Close()
	p.wg.Wait()
	log.Println("Pool stopped")
}
