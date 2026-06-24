package parser

//парсер HH.ru API
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"golang.org/x/time/rate"
)

type HHParser struct {
	client      *http.Client
	limiter     *rate.Limiter
	accessToken string
}

func NewHHParser(accessToken string) *HHParser {
	return &HHParser{
		client:      &http.Client{Timeout: 10 * time.Second},
		limiter:     rate.NewLimiter(rate.Every(500*time.Millisecond), 1),
		accessToken: accessToken,
	}
}

func (p *HHParser) Source() string {
	return "hh"
}

type hhResponse struct {
	Items []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Employer struct {
			Name string `json:"name"`
		} `json:"employer"`
		Area struct {
			Name string `json:"name"`
		} `json:"area"`
		Salary *struct {
			From     *int   `json:"from"`
			To       *int   `json:"to"`
			Currency string `json:"currency"`
		} `json:"salary"`
		Schedule struct {
			ID string `json:"id"`
		} `json:"schedule"`
		AlternateURL string `json:"alternate_url"`
		PublishedAt  string `json:"published_at"`
	} `json:"items"`
	Pages int `json:"pages"`
	Found int `json:"found"`
}

func (p *HHParser) Fetch(ctx context.Context, page int) ([]models.RawVacancy, error) {
	if err := p.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.hh.ru/vacancies?text=golang+python+java&area=1&page=%d&per_page=20",
		page,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	req.Header.Set("HH-User-Agent", "JobMarketAnalytics/1.0 (malayegor9@gmail.com)")
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

	var result hhResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	vacancies := make([]models.RawVacancy, 0, len(result.Items))
	for _, item := range result.Items {
		v := models.RawVacancy{
			ExternalID: item.ID,
			Source:     "hh",
			Title:      item.Name,
			Company:    item.Employer.Name,
			Location:   item.Area.Name,
			Remote:     item.Schedule.ID == "remote",
			URL:        item.AlternateURL,
		}
		if item.Salary != nil {
			v.SalaryFrom = item.Salary.From
			v.SalaryTo = item.Salary.To
			v.Currency = item.Salary.Currency
		}
		if t, err := time.Parse(time.RFC3339, item.PublishedAt); err == nil {
			v.PublishedAt = t
		}
		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}
