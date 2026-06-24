package parser

//интерфейс Parser, который реализует каждый источник
import (
	"context"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

type Parser interface {
	Source() string
	Fetch(ctx context.Context, page int) ([]models.RawVacancy, error)
}
