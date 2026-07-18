package waste

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateWaste(ctx context.Context, w *Waste) error
}
type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) CreateWaste(ctx context.Context, w *Waste) error {
	query := `
		INSERT INTO wastes (id, container_id, type, weight, carbon_footprint, fill_level_at_drop, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// pgxpool üzerinden SQL sorgusunu çalıştırıyoruz.
	_, err := r.db.Exec(ctx, query,
		w.ID,
		w.ContainerID,
		w.Type,
		w.Weight,
		w.CarbonFootprint,
		w.FillLevelAtDrop,
		w.CreatedAt,
	)

	if err != nil {
		// Hata mesajını yutmadan (swallow), nerede patladığını belirterek üst katmana iletiyoruz (Error Wrapping).
		return fmt.Errorf("atık kaydedilirken veritabanı hatası oluştu: %w", err)
	}

	return nil
}
