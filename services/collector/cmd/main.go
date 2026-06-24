package main

//точка входа collector'а: инициализирует зависимости, запускает парсинг
import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/dedup"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/parser"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/queue"
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

	redisClient, err := db.NewRedisClient(ctx, os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("✓ Redis connected")

	natsConn, err := db.NewNATSConn(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer natsConn.Close()
	log.Println("✓ NATS connected")

	deduplicator := dedup.New(redisClient)
	publisher := queue.NewPublisher(natsConn)
	remoteok := parser.NewRemoteOKParser()

	log.Println("Fetching vacancies from RemoteOK...")
	vacancies, err := remoteok.Fetch(ctx, 0)
	if err != nil {
		log.Fatalf("fetch: %v", err)
	}
	log.Printf("Fetched %d vacancies", len(vacancies))

	published := 0
	skipped := 0
	for _, v := range vacancies {
		// Проверяем дубль через Redis
		isDup, err := deduplicator.IsDuplicate(ctx, v.Source, v.ExternalID)
		if err != nil {
			log.Printf("dedup check %s: %v", v.ExternalID, err)
			continue
		}
		if isDup {
			skipped++
			continue
		}

		// Отправляем в NATS
		if err := publisher.Publish(ctx, v); err != nil {
			log.Printf("publish %s: %v", v.ExternalID, err)
			continue
		}
		published++
	}

	log.Printf("Published: %d, Skipped (dedup): %d", published, skipped)
	_ = pgPool

	// Habr Career
	log.Println("Fetching vacancies from Habr Career...")
	habr := parser.NewHabrParser()

	// Собираем по трём фидам (go, python, java)
	for feedPage := 0; feedPage < 3; feedPage++ {
		habrVacs, err := habr.Fetch(ctx, feedPage)
		if err != nil {
			log.Printf("habr fetch page %d error: %v", feedPage, err)
			continue
		}

		for _, v := range habrVacs {
			isDup, err := deduplicator.IsDuplicate(ctx, v.Source, v.ExternalID)
			if err != nil || isDup {
				skipped++
				continue
			}
			if err := publisher.Publish(ctx, v); err != nil {
				log.Printf("publish habr %s: %v", v.ExternalID, err)
				continue
			}
			published++
		}
		log.Printf("Habr feed %d: fetched %d vacancies", feedPage, len(habrVacs))
	}
}
