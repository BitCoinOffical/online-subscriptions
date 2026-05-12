package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	pool *pgxpool.Pool
}

func NewSubscription(pool *pgxpool.Pool) *SubscriptionRepo {
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
			return nil, fmt.Errorf("subscription not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("r.pool.QueryRow: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepo) UpdateSubscriptionsById(ctx context.Context, sub *models.PatchSubscription) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := psql.Update("subscriptions")
	if sub.ServiceName != "" {
		q = q.Set("service_name", sub.ServiceName)
	}
	if sub.Price != nil {
		q = q.Set("price", sub.Price)
	}
	if sub.UserID != uuid.Nil {
		q = q.Set("user_id", sub.UserID)
	}
	if !sub.StartDate.IsZero() {
		q = q.Set("start_date", sub.StartDate)
	}
	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		q = q.Set("end_date", sub.EndDate)
	}

	sql, args, err := q.Where(squirrel.Eq{"id": sub.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("q.ToSql():%w", err)
	}
	res, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *SubscriptionRepo) FullUpdateSubscriptionsById(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW() WHERE id = $6`
	res, err := r.pool.Exec(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)

	if err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *SubscriptionRepo) DeleteSubscriptions(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *SubscriptionRepo) GetSubscriptions(ctx context.Context, limit int, offset int) ([]models.Subscription, error) {

	query := `
		SELECT *
		FROM subscriptions
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Query: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription

	for rows.Next() {
		var sub models.Subscription

		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.Created_at,
			&sub.Updated_at,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
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
	fmt.Println("SQL:", query, "ARGS:", args)
	var total int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("r.pool.QueryRow: %w", err)
	}
	return total, nil
}
