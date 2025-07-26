package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Roflan4eg/test_work/internal/models"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"time"
)

type SubscriptionStorage struct {
	db *sql.DB
}

func New(db *sql.DB) (*SubscriptionStorage, error) {
	storage := &SubscriptionStorage{db: db}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect db: %w", err)
	}

	if err := storage.createMigrations(); err != nil {
		return nil, fmt.Errorf("migrations error: %w", err)
	}

	return storage, nil
}

func (s *SubscriptionStorage) createMigrations() error {
	var tableExists bool
	err := s.db.QueryRow(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'subscriptions'
        )`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if !tableExists {
		log.Println("Initializing database...")
		_, err := s.db.Exec(`
            CREATE TABLE subscriptions (
                id SERIAL PRIMARY KEY,
                service_name VARCHAR(255) NOT NULL,
                price INTEGER NOT NULL,
                user_id UUID NOT NULL,
                start_date DATE NOT NULL,
                end_date DATE
            )`)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
		log.Println("Database initialized")
	}
	return nil
}

func (s *SubscriptionStorage) Create(sub *models.Subscription) error {
	query := `INSERT INTO subscriptions
			(service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (s *SubscriptionStorage) GetByID(id int) (*models.Subscription, error) {
	sub := &models.Subscription{}
	query := `SELECT * FROM subscriptions WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription not found or deleted")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return sub, err
}

func (s *SubscriptionStorage) Delete(id int) error {
	res, err := s.db.Exec(`DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error, failed to check deletion: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found or deleted")
	}
	return err
}

func (s *SubscriptionStorage) Update(sub *models.Subscription) error {
	query := `UPDATE subscriptions 
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6`
	res, err := s.db.Exec(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error, failed to check deletion: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found or deleted")
	}
	return err
}

func (s *SubscriptionStorage) List() ([]models.Subscription, error) {
	query := `SELECT * FROM subscriptions`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()
	subs := []models.Subscription{}
	for rows.Next() {
		sub := models.Subscription{}
		if err = rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, fmt.Errorf("database error, row scan error: %w", err)
		}
		sub.FormatDates()
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SubscriptionStorage) GetSubsForPeriod(start, end time.Time, uuid string, serviceName string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0) 
        FROM subscriptions 
        WHERE start_date BETWEEN $1 AND $2`
	params := []interface{}{start, end}
	paramCounter := 3
	if uuid != "" {
		query += fmt.Sprintf(" AND user_id = $%v", paramCounter)
		params = append(params, uuid)
		paramCounter++
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%v", paramCounter)
		params = append(params, serviceName)
	}
	totalPrice := 0
	err := s.db.QueryRow(query, params...).Scan(&totalPrice)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}

	return totalPrice, nil
}
