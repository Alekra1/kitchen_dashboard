package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Order struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var ErrOrderNotFound = errors.New("order not found")

var Pool *pgxpool.Pool

// Connect initializes the connection pool using environment variables.
func Connect(ctx context.Context) error {
	_ = godotenv.Load()

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return err
	}

	Pool = pool
	return nil
}

// Close releases database resources.
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

// ListOrders returns all orders ordered by id.
func ListOrders(ctx context.Context) ([]Order, error) {
	rows, err := Pool.Query(ctx, "SELECT id, name, price FROM orders ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Name, &o.Price); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrder returns the order with the given id.
func GetOrder(ctx context.Context, id string) (Order, error) {
	var o Order
	err := Pool.QueryRow(ctx, "SELECT id, name, price FROM orders WHERE id = $1", id).
		Scan(&o.ID, &o.Name, &o.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, ErrOrderNotFound
	}
	return o, err
}

// CreateOrder inserts a new order and returns its stored representation.
func CreateOrder(ctx context.Context, o Order) (Order, error) {
	err := Pool.QueryRow(
		ctx,
		"INSERT INTO orders (id, name, price) VALUES ($1, $2, $3) RETURNING id, name, price",
		o.ID, o.Name, o.Price,
	).Scan(&o.ID, &o.Name, &o.Price)
	return o, err
}
