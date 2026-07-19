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
	// 1. İşlem Bloğunu (Transaction) Başlat
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction başlatılamadı: %w", err)
	}

	// defer ile fonksiyon bittiğinde eğer işlem onaylanmamışsa (commit edilmemişse)
	// her şeyi otomatik olarak geri almasını (Rollback) garanti ediyoruz.
	defer tx.Rollback(ctx)

	// 2. Birinci İşlem: Yeni Atığı Wastes Tablosuna Ekle
	insertWasteQuery := `
		INSERT INTO wastes (id, container_id, type, weight, carbon_footprint, fill_level_at_drop, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertWasteQuery,
		w.ID, w.ContainerID, w.Type, w.Weight, w.CarbonFootprint, w.FillLevelAtDrop, w.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("atık kaydedilemedi (işlem iptal ediliyor): %w", err)
	}

	// 3. İkinci İşlem: İlgili Çöp Kutusunun Doluluk Oranını Güncelle
	// Burada, eklenen atığın ağırlığına göre kutunun mevcut doluluğunu artırıyoruz.
	// LEAST komutu, doluluğun %100'ü geçmesini engeller (Veritabanı seviyesinde güvenlik).
	updateContainerQuery := `
		UPDATE containers 
		SET current_fill_level = LEAST(current_fill_level + $1, 100.0)
		WHERE id = $2
	`
	// Simülasyon için 1 kg atığın %1 doluluk artırdığını varsayalım
	// (İleride kapasite formülüyle geliştirilebilir)
	fillIncrease := w.Weight

	_, err = tx.Exec(ctx, updateContainerQuery, fillIncrease, w.ContainerID)
	if err != nil {
		return fmt.Errorf("kutu doluluğu güncellenemedi (işlem iptal ediliyor): %w", err)
	}

	// 4. Her şey başarılıysa işlemleri kalıcı olarak onayla (Commit)
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction onaylanamadı: %w", err)
	}

	return nil
}
