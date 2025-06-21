package queue

import (
	"container/heap"
	"log"
	"notification_service/internal/models"
	"sync"
)

type Item struct {
	Notification *models.Notification
	Index        int
}

type NotificationQueue []*Item

func (nq NotificationQueue) Len() int {
	return len(nq)
}

type Queue struct {
	pq     *NotificationQueue
	mu     *sync.Mutex
	closed bool
}

func NewQueue() *Queue {
	return &Queue{pq: &NotificationQueue{}, mu: &sync.Mutex{}, closed: false}
}

func (q *Queue) Push(notification *models.Notification) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		log.Println("Queue is closed")
		return
	}

	item := &Item{
		Notification: notification,
	}

	heap.Push(q.pq, item)
}

func (q *Queue) Pop() *models.Notification {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.pq.Len() > 0 {
		item := heap.Pop(q.pq).(*Item)
		return item.Notification
	}

	return nil
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
}
