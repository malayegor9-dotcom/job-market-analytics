package parser

//парсер RemoteOK API
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

type RemoteOKParser struct {
	client *http.Client
}

func NewRemoteOKParser() *RemoteOKParser {
	return &RemoteOKParser{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *RemoteOKParser) Source() string {
	return "remoteok"
}

type remoteOKJob struct {
	ID          string   `json:"id"`
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
	Date        string   `json:"date"`
	SalaryMin   int      `json:"salary_min"`
	SalaryMax   int      `json:"salary_max"`
}

func (p *RemoteOKParser) Fetch(ctx context.Context, page int) ([]models.RawVacancy, error) {
	if page > 0 {
		return []models.RawVacancy{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://remoteok.com/api", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	vacancies := make([]models.RawVacancy, 0)
	for i, item := range raw {
		if i == 0 {
			continue
		}

		var job remoteOKJob
		if err := json.Unmarshal(item, &job); err != nil {
			continue
		}
		if job.Position == "" {
			continue
		}

		v := models.RawVacancy{
			ExternalID:  job.ID,
			Source:      "remoteok",
			Title:       job.Position,
			Company:     job.Company,
			Description: job.Description,
			Location:    job.Location,
			Remote:      true,
			URL:         job.URL,
		}

		if job.SalaryMin > 0 {
			v.SalaryFrom = &job.SalaryMin
		}
		if job.SalaryMax > 0 {
			v.SalaryTo = &job.SalaryMax
		}
		if job.SalaryMin > 0 || job.SalaryMax > 0 {
			v.Currency = "USD"
		}
		if t, err := time.Parse(time.RFC3339, job.Date); err == nil {
			v.PublishedAt = t
		}

		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}
