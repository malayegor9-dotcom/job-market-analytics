package extractor

//извлечение навыков из текста вакансии по словарю
import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SkillExtractor struct {
	// map[keyword]skill_id — загружается из БД при старте
	index map[string]int
}

// New загружает все навыки из БД в память
func New(ctx context.Context, pool *pgxpool.Pool) (*SkillExtractor, error) {
	rows, err := pool.Query(ctx, `SELECT id, name FROM skills`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	index := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		// Ключ в нижнем регистре для case-insensitive матчинга
		index[strings.ToLower(name)] = id
	}

	return &SkillExtractor{index: index}, nil
}

// Extract находит навыки в тексте вакансии
// Возвращает список skill_id
func (e *SkillExtractor) Extract(title, description string) []int {
	// Объединяем заголовок и описание
	text := strings.ToLower(title + " " + description)

	seen := make(map[int]bool)
	var result []int

	for keyword, skillID := range e.index {
		// Ищем навык как отдельное слово
		// Например "go" не должен матчиться в "golang"
		if containsWord(text, keyword) {
			if !seen[skillID] {
				seen[skillID] = true
				result = append(result, skillID)
			}
		}
	}

	return result
}

// containsWord проверяет что слово встречается отдельно
// а не как часть другого слова
func containsWord(text, word string) bool {
	idx := strings.Index(text, word)
	if idx == -1 {
		return false
	}

	// Проверяем символ перед словом
	if idx > 0 {
		before := text[idx-1]
		if isLetter(before) {
			return false
		}
	}

	// Проверяем символ после слова
	end := idx + len(word)
	if end < len(text) {
		after := text[end]
		if isLetter(after) {
			return false
		}
	}

	return true
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
