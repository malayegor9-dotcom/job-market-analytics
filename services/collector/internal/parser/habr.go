package parser

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

// HabrParser парсит вакансии с Habr Career через RSS feed
type HabrParser struct {
	client *http.Client
}

func NewHabrParser() *HabrParser {
	return &HabrParser{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *HabrParser) Source() string {
	return "habr"
}

// RSS структура
type habrRSS struct {
	Channel struct {
		Items []habrItem `xml:"item"`
	} `xml:"channel"`
}

type habrItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
}

func (p *HabrParser) Fetch(ctx context.Context, page int) ([]models.RawVacancy, error) {
	// Habr Career RSS для разных специализаций
	feeds := []string{
		"https://career.habr.com/vacancies/rss",
		"https://career.habr.com/vacancies/rss?type=all&q=golang",
		"https://career.habr.com/vacancies/rss?type=all&q=python",
	}

	// Берём только нужный feed по номеру страницы
	if page >= len(feeds) {
		return []models.RawVacancy{}, nil
	}
	url := feeds[page]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "JobMarketAnalytics/1.0 (malayegor9@gmail.com)")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var rss habrRSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("decode rss: %w", err)
	}

	vacancies := make([]models.RawVacancy, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		// Извлекаем ID из URL
		// URL вида: https://career.habr.com/vacancies/1234567
		parts := strings.Split(strings.TrimRight(item.Link, "/"), "/")
		externalID := parts[len(parts)-1]
		if externalID == "" {
			continue
		}

		// Извлекаем компанию из поля author
		company := strings.TrimSpace(item.Author)
		if company == "" {
			company = "Unknown"
		}

		v := models.RawVacancy{
			ExternalID:  externalID,
			Source:      "habr",
			Title:       strings.TrimSpace(item.Title),
			Description: item.Description,
			Company:     company,
			URL:         item.Link,
			Remote: strings.Contains(strings.ToLower(item.Description), "удалённо") ||
				strings.Contains(strings.ToLower(item.Description), "remote"),
		}

		if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate); err == nil {
			v.PublishedAt = t
		} else if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			v.PublishedAt = t
		}

		vacancies = append(vacancies, v)
	}

	return vacancies, nil
}
