package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/giannuccilli/user-api/internal/domain"
)

type FailedEventRepository struct {
	pool *pgxpool.Pool
}

func NewFailedEventRepository(pool *pgxpool.Pool) *FailedEventRepository {
	return &FailedEventRepository{pool: pool}
}

func (r *FailedEventRepository) Save(ctx context.Context, event *domain.FailedEvent) error {
	query := `
		INSERT INTO failed_events (event_id, event_type, user_id, payload, error, attempts, created_at, last_error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		event.EventID,
		event.EventType,
		event.UserID,
		event.Payload,
		event.Error,
		event.Attempts,
		event.CreatedAt,
		event.LastError,
	).Scan(&event.ID)

	return err
}

func (r *FailedEventRepository) List(ctx context.Context, limit, offset int) ([]domain.FailedEvent, int, error) {
	countQuery := `SELECT COUNT(*) FROM failed_events`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, event_id, event_type, user_id, payload, error, attempts, created_at, last_error
		FROM failed_events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []domain.FailedEvent
	for rows.Next() {
		var e domain.FailedEvent
		if err := rows.Scan(
			&e.ID,
			&e.EventID,
			&e.EventType,
			&e.UserID,
			&e.Payload,
			&e.Error,
			&e.Attempts,
			&e.CreatedAt,
			&e.LastError,
		); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *FailedEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM failed_events WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
