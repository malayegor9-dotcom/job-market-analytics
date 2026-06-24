package handler

//эндпоинты аналитики: топ навыков, статистика зарплат
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SkillStat struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// TopSkills GET /api/v1/analytics/skills
func (h *Handler) TopSkills(c *gin.Context) {
	rows, err := h.db.Query(c, `
		SELECT s.name, s.category, COUNT(*) as count
		FROM vacancy_skills vs
		JOIN skills s ON s.id = vs.skill_id
		GROUP BY s.name, s.category
		ORDER BY count DESC
		LIMIT 20
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	skills := []SkillStat{}
	for rows.Next() {
		var s SkillStat
		if err := rows.Scan(&s.Name, &s.Category, &s.Count); err != nil {
			continue
		}
		skills = append(skills, s)
	}

	c.JSON(http.StatusOK, gin.H{"data": skills})
}

var stats struct {
	Total      int     `json:"total_vacancies"`
	WithSalary int     `json:"with_salary"`
	AvgMin     float64 `json:"avg_salary_min"`
	AvgMax     float64 `json:"avg_salary_max"`
}

// SalaryStats GET /api/v1/analytics/salary
func (h *Handler) SalaryStats(c *gin.Context) {

	err := h.db.QueryRow(c, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE salary_min IS NOT NULL),
			COALESCE(AVG(salary_min) FILTER (WHERE salary_min IS NOT NULL), 0),
			COALESCE(AVG(salary_max) FILTER (WHERE salary_max IS NOT NULL), 0)
		FROM vacancies
	`).Scan(&stats.Total, &stats.WithSalary, &stats.AvgMin, &stats.AvgMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
