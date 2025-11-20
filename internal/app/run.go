package app

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	postgres "Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	"context"
	"fmt"
)

/*
func InitDB(db *postgres.DatabaseSource, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	if _, err = db.Pool.Exec(context.Background(), string(file)); err != nil {
		return fmt.Errorf("error executing sql: %v", err)
	}
	givenOrder, err := utils.GetGivenOrder()
	if err != nil {
		return fmt.Errorf("error getting givenOrder: %v", err)
	}

	if err = postgres.AddOrdersToDB(db, givenOrder); err != nil {
		return fmt.Errorf("error adding orders to database: %v", err)
	}
	return nil
}*/

func initPostgres(cfg *config.Config) (*database.DatabaseSource, error) {
	db, err := database.NewStorage(
		database.GetConnection(&cfg.Database),
		database.SetMaxPoolSize(cfg.Database.MaxPoolSize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}
	if err = db.Pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}
	return db, nil
}

func initDBRepository(db *database.DatabaseSource) *postgres.Repository {
	return postgres.CreateNewDBRepository(db)
}

func Run(cfg *config.Config) error {
	db, err := initPostgres(cfg)
	if err != nil {
		return err
	}
	defer func(db *database.DatabaseSource) {
		db.Close()
	}(db)

	pgRepository := initDBRepository(db)
	service := CreateNewOrderService(pgRepository)
	orderController := controller.CreateNewOrderController(orderService)
	if err = startServer(cfg, orderController); err != nil {
		return fmt.Errorf("start server error: %v", err)
	}
	return nil
}
