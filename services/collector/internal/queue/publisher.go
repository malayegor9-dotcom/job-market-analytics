package queue

//отправка вакансий в NATS очередь
import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"github.com/nats-io/nats.go"
)

const TopicRawVacancies = "vacancies.raw"

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) Publish(ctx context.Context, v models.RawVacancy) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal vacancy: %w", err)
	}

	if err := p.nc.Publish(TopicRawVacancies, data); err != nil {
		return fmt.Errorf("nats publish: %w", err)
	}

	return nil
}
