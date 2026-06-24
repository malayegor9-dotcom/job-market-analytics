package models

//общая структура RawVacancy, используется всеми сервисами
import "time"

type RawVacancy struct {
	ExternalID  string
	Source      string
	Title       string
	Description string
	Company     string
	Location    string
	SalaryFrom  *int
	SalaryTo    *int
	Currency    string
	Remote      bool
	URL         string
	PublishedAt time.Time
}
