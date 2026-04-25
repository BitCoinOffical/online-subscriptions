package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/models"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscriptionHandler(pool *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) GetSubscriptionsById(ctx context.Context, id int) (*models.Subscription, error) {
	var sub models.Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.Created_at, &sub.Updated_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found subscription by id: %w", err)
		}
		return nil, fmt.Errorf("r.pool.QueryRow: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepo) UpdateSubscriptions(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.pool.Exec(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)

	if err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}

func (r *SubscriptionRepo) DeleteSubscriptions(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) GetSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	query := `SELECT * FROM subscriptions ORDER BY id DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription

	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.Created_at, &sub.Updated_at); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return subs, nil
}
func (r *SubscriptionRepo) GetSubscriptionsFilter(ctx context.Context, from, to, user_id, service_name string) (int, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := psql.Select("COALESCE(SUM(price), 0)").From("subscriptions")
	if from != "" {
		q = q.Where(squirrel.GtOrEq{"start_date": from})
	}
	if to != "" {
		q = q.Where(squirrel.LtOrEq{"start_date": to})
	}
	if user_id != "" {
		q = q.Where(squirrel.Eq{"user_id": user_id})
	}
	if service_name != "" {
		q = q.Where(squirrel.Eq{"service_name": service_name})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf("q.ToSql():%w", err)
	}
	var total int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("r.pool.QueryRow: %w", err)
	}
	return total, nil
}
