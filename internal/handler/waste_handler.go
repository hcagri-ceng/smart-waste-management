package handler

import (
	"log"
	"smartwaste/internal/domain/waste"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WasteHandler struct {
	repo waste.Repository
}

// NewWasteHandler, dependency injection ile yeni bir handler oluşturur.
func NewWasteHandler(repo waste.Repository) *WasteHandler {
	return &WasteHandler{
		repo: repo,
	}
}

// CreateWasteRequest, dışarıdan gelen JSON verisinin şemasıdır (DTO - Data Transfer Object).
// Kullanıcıdan sadece bu 3 veriyi bekliyoruz, geri kalan her şeyi backend kendisi hesaplayacak.
type CreateWasteRequest struct {
	ContainerID string  `json:"container_id"`
	Type        string  `json:"type"`
	Weight      float64 `json:"weight"`
}

// HandleCreateWaste, yeni atık girişini alır, işler ve veritabanına gönderir.
func (h *WasteHandler) HandleCreateWaste(c *fiber.Ctx) error {
	// 1. Gelen JSON verisini DTO'ya ayrıştır (Parse)
	var req CreateWasteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Geçersiz veri formatı. Lütfen JSON yapısını kontrol edin.",
		})
	}

	// 2. Container UUID Doğrulaması
	containerUUID, err := uuid.Parse(req.ContainerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Geçersiz Container ID formatı.",
		})
	}

	// 3. İş Mantığı (Business Logic) - Karbon Ayak İzini Hesapla
	carbonFootprint := waste.CalculateCarbonFootprint(req.Type, req.Weight)

	// 4. Veritabanı Modelini (Struct) Oluştur
	// Dışarıdan gelmeyen ID, Tarih ve Karbon verilerini sistem kendisi atıyor.
	newWaste := &waste.Waste{
		ID:              uuid.New().String(),    // <-- .String() eklendi
		ContainerID:     containerUUID.String(), // <-- .String() eklendi
		Type:            req.Type,
		Weight:          req.Weight,
		CarbonFootprint: carbonFootprint,
		FillLevelAtDrop: 0,
		CreatedAt:       time.Now(),
	}

	// 5. Repository üzerinden işlemi gerçekleştir (Transaction bloğu burada çalışır)
	if err := h.repo.CreateWaste(c.Context(), newWaste); err != nil {
		log.Printf("Atık kaydedilemedi: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Atık işlenirken sunucu tarafında bir hata oluştu.",
		})
	}

	// 6. Başarılı Yanıt Dön (201 Created)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Atık başarıyla kaydedildi ve kutu doluluğu güncellendi.",
		"data": fiber.Map{
			"waste_id":         newWaste.ID,
			"carbon_footprint": newWaste.CarbonFootprint,
		},
	})
}
