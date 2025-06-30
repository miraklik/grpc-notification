package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notification_service/internal/models"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	streamName = "NOTIFICATIONS"
	subjTmpl   = "notifications.p%d"
)

type NatsQueue struct {
	jes       nats.JetStreamContext
	consumers [5]*nats.Subscription
}

func NewNatsQueue(url string) (*NatsQueue, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Printf("Error connecting to nats: %v", err)
		return nil, err
	}

	js, _ := nc.JetStream()

	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"notifications.*"},
		Storage:  nats.FileStorage,
	})

	q := &NatsQueue{jes: js}

	for p := 1; p <= 5; p++ {
		subj := fmt.Sprintf(subjTmpl, p)
		cons, err := js.PullSubscribe(subj, fmt.Sprintf("workers-p%d", p),
			nats.PullMaxWaiting(128),
			nats.ManualAck(),
		)
		if err != nil {
			return nil, err
		}
		q.consumers[p-1] = cons
	}
	return q, nil
}

func (q *NatsQueue) Push(n *models.Notification) error {
	data, _ := json.Marshal(n)
	subj := fmt.Sprintf(subjTmpl, n.Priority)
	_, err := q.jes.Publish(subj, data)
	return err
}

func (q *NatsQueue) Pop(ctx context.Context) (*nats.Msg, error) {
	for pri := 0; pri < 5; pri++ {
		msgs, err := q.consumers[pri].Fetch(1, nats.Context(ctx), nats.MaxWait(100*time.Millisecond))
		if err == nats.ErrTimeout {
			continue
		}
		if err != nil {
			return nil, err
		}
		if len(msgs) == 1 {
			return msgs[0], nil
		}
	}
	return nil, nil
}
