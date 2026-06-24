package storage

//сохранение вакансий и навыков в PostgreSQL
import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) SaveVacancy(ctx context.Context, v models.RawVacancy) error {
	var sourceID int
	err := s.pool.QueryRow(ctx,
		`SELECT id FROM sources WHERE name = $1`, v.Source,
	).Scan(&sourceID)
	if err != nil {
		return fmt.Errorf("get source id: %w", err)
	}

	var companyID int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO companies (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, v.Company).Scan(&companyID)
	if err != nil {
		return fmt.Errorf("upsert company: %w", err)
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO vacancies (
			external_id, source_id, company_id,
			title, description, location,
			salary_min, salary_max, currency,
			remote, url, published_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12
		)
		ON CONFLICT (source_id, external_id) DO NOTHING
	`,
		v.ExternalID, sourceID, companyID,
		v.Title, v.Description, v.Location,
		v.SalaryFrom, v.SalaryTo, v.Currency,
		v.Remote, v.URL, v.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("insert vacancy: %w", err)
	}

	return nil
}

func (s *Storage) CountVacancies(ctx context.Context) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM vacancies`).Scan(&count)
	return count, err
}

// SaveVacancyWithID сохраняет вакансию и возвращает её ID
// Если вакансия уже существует — возвращает 0
func (s *Storage) SaveVacancyWithID(ctx context.Context, v models.RawVacancy) (int, error) {
	var sourceID int
	err := s.pool.QueryRow(ctx,
		`SELECT id FROM sources WHERE name = $1`, v.Source,
	).Scan(&sourceID)
	if err != nil {
		return 0, fmt.Errorf("get source id: %w", err)
	}

	var companyID int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO companies (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, v.Company).Scan(&companyID)
	if err != nil {
		return 0, fmt.Errorf("upsert company: %w", err)
	}

	var vacancyID int
	err = s.pool.QueryRow(ctx, `
		INSERT INTO vacancies (
			external_id, source_id, company_id,
			title, description, location,
			salary_min, salary_max, currency,
			remote, url, published_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12
		)
		ON CONFLICT (source_id, external_id) DO NOTHING
		RETURNING id
	`,
		v.ExternalID, sourceID, companyID,
		v.Title, v.Description, v.Location,
		v.SalaryFrom, v.SalaryTo, v.Currency,
		v.Remote, v.URL, v.PublishedAt,
	).Scan(&vacancyID)

	// Если вакансия уже существовала — Scan вернёт ошибку (нет строки)
	if err != nil {
		return 0, nil
	}

	return vacancyID, nil
}

// SaveVacancySkills сохраняет связи вакансии с навыками
func (s *Storage) SaveVacancySkills(ctx context.Context, vacancyID int, skillIDs []int) error {
	for _, skillID := range skillIDs {
		_, err := s.pool.Exec(ctx, `
			INSERT INTO vacancy_skills (vacancy_id, skill_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, vacancyID, skillID)
		if err != nil {
			return fmt.Errorf("insert vacancy_skill: %w", err)
		}
	}
	return nil
}
