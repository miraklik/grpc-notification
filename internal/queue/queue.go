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

func (nq NotificationQueue) Less(i, j int) bool {
	return nq[i].Notification.Priority > nq[j].Notification.Priority
}

func (nq NotificationQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
	nq[i].Index = i
	nq[j].Index = j
}

func (nq *NotificationQueue) Push(x any) {
	n := len(*nq)
	item := x.(*Item)
	item.Index = n
	*nq = append(*nq, item)
}

func (nq *NotificationQueue) Pop() any {
	old := *nq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*nq = old[0 : n-1]
	return item
}

type Queue struct {
	pq     *NotificationQueue
	mu     *sync.Mutex
	closed bool
}

func NewQueue() *Queue {
	pq := make(NotificationQueue, 0)
	heap.Init(&pq)
	q := &Queue{
		pq:     &pq,
		mu:     &sync.Mutex{},
		closed: false,
	}

	return q
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

	if q.closed {
		log.Println("Queue is closed")
		return nil
	}

	return nil
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
}
