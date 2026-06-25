# 📊 Job Market Analytics

Платформа для сбора и анализа вакансий с различных источников.  
Go backend + Python визуализация на Streamlit.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)
![Python](https://img.shields.io/badge/Python-3.12-3776AB?style=flat&logo=python)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)

---

## ✨ Возможности

- **Сбор вакансий** с RemoteOK и Habr Career
- **Дедупликация** через Redis — не сохраняем одно и то же дважды
- **Асинхронная обработка** через очередь NATS
- **Извлечение навыков** из текста вакансий по словарю
- **REST API** для доступа к данным
- **Дашборд** с графиками, картой навыков и фильтрами

---

## 🚀 Быстрый старт

### Требования

- Go 1.21+
- Python 3.12+
- Docker Desktop

### 1. Клонирование репозитория

```bash
git clone https://github.com/malayegor9-dotcom/job-market-analytics
cd job-market-analytics
```

### 2. Настройка переменных окружения

```bash
cp .env.example .env
```

Открой `.env` и заполни:

```env
DATABASE_URL=postgres://app@127.0.0.1:5433/jobmarket?sslmode=disable
REDIS_URL=redis://localhost:6379/0
NATS_URL=nats://localhost:4222
API_PORT=8080
```

### 3. Настройка инфраструктуры

```bash
cd infra
docker compose up -d
```

### 4. Применение миграций

```bash
# Из папки infra
Get-Content ./migrations/001_sources.up.sql | docker exec -i jobmarket_pg psql -U app -d jobmarket
Get-Content ./migrations/002_companies.up.sql | docker exec -i jobmarket_pg psql -U app -d jobmarket
Get-Content ./migrations/003_create_skills.up.sql | docker exec -i jobmarket_pg psql -U app -d jobmarket
Get-Content ./migrations/004_create_vacancies.up.sql | docker exec -i jobmarket_pg psql -U app -d jobmarket
Get-Content ./migrations/005_create_vacancy_skills.up.sql | docker exec -i jobmarket_pg psql -U app -d jobmarket
```

### 5. Запуск сервисов

Открой три терминала:

```bash
# Терминал 1 — Processor
go run ./services/processor/cmd/main.go

# Терминал 2 — Collector
go run ./services/collector/cmd/main.go

# Терминал 3 — API
go run ./services/api/cmd/main.go
```

### 6. Запуск дашборда

```bash
cd analytics
py -3.12 -m streamlit run dashboard.py
```

---

## 🔌 API эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/vacancies` | Список вакансий с пагинацией |
| GET | `/api/v1/vacancies/:id` | Детали вакансии |
| GET | `/api/v1/vacancies/search?q=golang` | Полнотекстовый поиск |
| GET | `/api/v1/analytics/skills` | Топ навыков |
| GET | `/api/v1/analytics/salary` | Статистика зарплат |

---

## 🛠️ Технологии

| Слой | Технологии |
|------|-----------|
| Backend | Go, Gin, pgx/v5 |
| Очередь | NATS JetStream |
| Кэш | Redis |
| База данных | PostgreSQL 16 |
| Визуализация | Python, Streamlit, Plotly |
| Инфраструктура | Docker Compose |

---

## 📦 Источники вакансий

| Источник | Тип | Статус |
|----------|-----|--------|
| RemoteOK | JSON API | ✅ Активен |
| Habr Career | RSS | ✅ Активен |
| HH.ru | OAuth 2.0 API | 🔧 В разработке |

---

## 📈 Дашборд

- **Обзор** — метрики, динамика публикаций, источники
- **Навыки** — топ технологий, карта, категории
- **Компании** — рейтинг работодателей
- **Вакансии** — полный список с поиском и фильтрами

**Таблицы:** `sources`, `companies`, `vacancies`, `skills`, `vacancy_skills`

---

## 👤 Автор

**Egor Malay** — [@malayegor9-dotcom](https://github.com/malayegor9-dotcom)