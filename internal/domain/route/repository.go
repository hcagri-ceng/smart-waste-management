package route

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetOptimalRoutes(ctx context.Context) ([]Container, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) GetOptimalRoutes(ctx context.Context) ([]Container, error) {
	query := `
		SELECT 
			id, latitude, longitude, capacity, current_fill_level, 
			temperature, gas_level_ppm, battery_status, last_emptied_at
		FROM containers
		WHERE 
			temperature >= 65.0 OR gas_level_ppm >= 20.0 
			OR current_fill_level > 50.0
		ORDER BY 
			-- Acil durumlara (1), normal durumlara (0) vererek listeyi force et.
			CASE WHEN temperature >= 65.0 OR gas_level_ppm >= 20.0 THEN 1 ELSE 0 END DESC,
			-- Normalize edilmiş öncelik skoru (Puanlama Algoritması)
			((current_fill_level / 100.0) * 0.60) + 
			((gas_level_ppm / 100.0) * 0.15) + 
			((temperature / 100.0) * 0.10) DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("optimal rotalar hesaplanırken veritabanı hatası oluştu: %w", err)
	}
	defer rows.Close()

	var routes []Container

	// Dönen verileri sırayla gezip Container struct'larına dönüştürüyoruz (Mapping).
	for rows.Next() {
		var c Container
		if err := rows.Scan(
			&c.ID, &c.Latitude, &c.Longitude, &c.Capacity,
			&c.CurrentFillLevel, &c.Temperature, &c.GasLevelPPM,
			&c.BatteryStatus, &c.LastEmptiedAt,
		); err != nil {
			return nil, fmt.Errorf("satır verisi struct'a dönüştürülürken hata: %w", err)
		}
		routes = append(routes, c)
	}

	// Döngü sırasında gizli bir hata olup olmadığını kontrol et.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("satırlar işlenirken görünmez bir hata oluştu: %w", err)
	}

	return routes, nil
}
