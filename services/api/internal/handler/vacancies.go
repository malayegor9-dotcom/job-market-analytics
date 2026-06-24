package handler

//эндпоинты для вакансий: список, детали, поиск
import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ListVacancies GET /api/v1/vacancies
// Параметры: page, limit, location, remote
func (h *Handler) ListVacancies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	location := c.Query("location")
	remote := c.Query("remote")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `
		SELECT v.id, v.title, v.location, v.salary_min, v.salary_max,
		       v.currency, v.remote, v.url, v.published_at,
		       c.name as company
		FROM vacancies v
		LEFT JOIN companies c ON c.id = v.company_id
		WHERE 1=1
	`
	args := []any{}
	argN := 1

	if location != "" {
		query += ` AND v.location ILIKE $` + strconv.Itoa(argN)
		args = append(args, "%"+location+"%")
		argN++
	}
	if remote == "true" {
		query += ` AND v.remote = true`
	}

	query += ` ORDER BY v.published_at DESC LIMIT $` + strconv.Itoa(argN) +
		` OFFSET $` + strconv.Itoa(argN+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(c, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type VacancyRow struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Location    string     `json:"location"`
		SalaryMin   *int64     `json:"salary_min"`
		SalaryMax   *int64     `json:"salary_max"`
		Currency    string     `json:"currency"`
		Remote      bool       `json:"remote"`
		URL         string     `json:"url"`
		PublishedAt *time.Time `json:"published_at"`
		Company     string     `json:"company"`
	}

	vacancies := []VacancyRow{}
	for rows.Next() {
		var v VacancyRow
		err := rows.Scan(
			&v.ID, &v.Title, &v.Location,
			&v.SalaryMin, &v.SalaryMax, &v.Currency,
			&v.Remote, &v.URL, &v.PublishedAt, &v.Company,
		)
		if err != nil {
			log.Println("scan error:", err)
			continue
		}
		vacancies = append(vacancies, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  vacancies,
		"page":  page,
		"limit": limit,
	})
}

// GetVacancy GET /api/v1/vacancies/:id
func (h *Handler) GetVacancy(c *gin.Context) {
	id := c.Param("id")

	type VacancyDetail struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Location    string     `json:"location"`
		SalaryMin   *int64     `json:"salary_min"`
		SalaryMax   *int64     `json:"salary_max"`
		Currency    string     `json:"currency"`
		Remote      bool       `json:"remote"`
		URL         string     `json:"url"`
		PublishedAt *time.Time `json:"published_at"`
		Company     string     `json:"company"`
	}

	var v VacancyDetail
	err := h.db.QueryRow(c, `
		SELECT v.id, v.title, v.description, v.location,
		       v.salary_min, v.salary_max, v.currency,
		       v.remote, v.url, v.published_at, c.name
		FROM vacancies v
		LEFT JOIN companies c ON c.id = v.company_id
		WHERE v.id = $1
	`, id).Scan(
		&v.ID, &v.Title, &v.Description, &v.Location,
		&v.SalaryMin, &v.SalaryMax, &v.Currency,
		&v.Remote, &v.URL, &v.PublishedAt, &v.Company,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}

	c.JSON(http.StatusOK, v)
}

// SearchVacancies GET /api/v1/vacancies/search?q=golang
func (h *Handler) SearchVacancies(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	rows, err := h.db.Query(c, `
		SELECT v.id, v.title, v.location, v.remote, v.url, c.name,
		       ts_rank(v.search_vector, to_tsquery('russian', $1)) as rank
		FROM vacancies v
		LEFT JOIN companies c ON c.id = v.company_id
		WHERE v.search_vector @@ to_tsquery('russian', $1)
		ORDER BY rank DESC
		LIMIT 20
	`, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type SearchResult struct {
		ID       int     `json:"id"`
		Title    string  `json:"title"`
		Location string  `json:"location"`
		Remote   bool    `json:"remote"`
		URL      string  `json:"url"`
		Company  string  `json:"company"`
		Rank     float64 `json:"rank"`
	}

	results := []SearchResult{}
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.ID, &r.Title, &r.Location,
			&r.Remote, &r.URL, &r.Company, &r.Rank); err != nil {
			continue
		}
		results = append(results, r)
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "query": q})
}
