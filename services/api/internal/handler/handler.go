package handler

//структура Handler, регистрация маршрутов
import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Handler держит все зависимости — пул БД и Redis клиент
type Handler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func New(db *pgxpool.Pool, redis *redis.Client) *Handler {
	return &Handler{db: db, redis: redis}
}

// RegisterRoutes регистрирует все маршруты API
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Вакансии
		api.GET("/vacancies", h.ListVacancies)
		api.GET("/vacancies/:id", h.GetVacancy)
		api.GET("/vacancies/search", h.SearchVacancies)

		// Аналитика
		api.GET("/analytics/skills", h.TopSkills)
		api.GET("/analytics/salary", h.SalaryStats)
	}
}
