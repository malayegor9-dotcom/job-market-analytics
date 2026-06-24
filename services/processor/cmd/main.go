package main

//точка входа processor'а: читает из NATS, сохраняет в БД
import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/storage"
	"github.com/malayegor9-dotcom/job-market-analytics/services/processor/internal/extractor"
	"github.com/nats-io/nats.go"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file, reading from environment")
	}

	ctx := context.Background()

	pgPool, err := db.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pgPool.Close()
	log.Println("✓ PostgreSQL connected")

	natsConn, err := db.NewNATSConn(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer natsConn.Close()
	log.Println("✓ NATS connected")

	store := storage.New(pgPool)

	// Загружаем навыки из БД в память
	skillExtractor, err := extractor.New(ctx, pgPool)
	if err != nil {
		log.Fatalf("skill extractor: %v", err)
	}
	log.Println("✓ Skill extractor loaded")

	processed := 0
	sub, err := natsConn.Subscribe("vacancies.raw", func(msg *nats.Msg) {
		var v models.RawVacancy
		if err := json.Unmarshal(msg.Data, &v); err != nil {
			log.Printf("unmarshal error: %v", err)
			return
		}

		// Сохраняем вакансию и получаем её ID
		vacancyID, err := store.SaveVacancyWithID(ctx, v)
		if err != nil {
			log.Printf("save vacancy %s: %v", v.ExternalID, err)
			return
		}

		// Если вакансия новая (ID > 0) — извлекаем навыки
		if vacancyID > 0 {
			skillIDs := skillExtractor.Extract(v.Title, v.Description)
			if len(skillIDs) > 0 {
				if err := store.SaveVacancySkills(ctx, vacancyID, skillIDs); err != nil {
					log.Printf("save skills for vacancy %d: %v", vacancyID, err)
				}
			}
		}

		processed++
		log.Printf("[%d] saved: %s — %s (skills: %d)",
			processed, v.Company, v.Title, len(skillExtractor.Extract(v.Title, v.Description)))
	})
	if err != nil {
		log.Fatalf("nats subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Processor started, waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down. Processed: %d vacancies", processed)
}
