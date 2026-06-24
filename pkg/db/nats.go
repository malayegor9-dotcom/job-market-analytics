package db

//создаёт соединение с NATS
import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func NewNATSConn(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(
		url,
		nats.Name("job-market-analytics"),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	return nc, nil
}
