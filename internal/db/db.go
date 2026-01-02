package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectAndMigrate(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect error: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	schemaPath := filepath.Join("internal", "db", "schema.sql")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		schemaPath = filepath.Join("..", "..", "internal", "db", "schema.sql")
	}

	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("read schema: %w", err)
	}

	sqlString := string(sqlBytes)
	if sqlString == "" {
		pool.Close()
		return nil, fmt.Errorf("schema.sql is empty")
	}

	fmt.Printf("Running migration...\n")

	_, err = pool.Exec(ctx, sqlString)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	fmt.Println("Migration completed successfully")
	return pool, nil
}
