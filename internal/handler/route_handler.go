package handler

import (
	"log"
	"smartwaste/internal/domain/route"

	"github.com/gofiber/fiber/v2"
)

// RouteHandler, rota işlemleri için gelen HTTP isteklerini karşılayan yapıdır.
type RouteHandler struct {
	repo route.Repository // Bağımlılık (Priz)
}

// NewRouteHandler, dependency injection ile yeni bir handler oluşturur.
func NewRouteHandler(repo route.Repository) *RouteHandler {
	return &RouteHandler{
		repo: repo,
	}
}

// GetOptimalRoutes, en acil ve optimal rotaları hesaplayıp JSON olarak dışarı fırlatır.
func (h *RouteHandler) GetOptimalRoutes(c *fiber.Ctx) error {
	// İstek iptal olursa diye fiber context'ini standart Go context'ine çeviriyoruz
	ctx := c.Context()

	// Repository'deki o güçlü algoritmayı (motoru) çalıştır
	routes, err := h.repo.GetOptimalRoutes(ctx)
	if err != nil {
		log.Printf("Rota hesaplama hatası: %v\n", err)
		// Dışarıya iç sistem hatasını belli etmeden profesyonel bir 500 dönüyoruz
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Optimal rotalar hesaplanırken sistemsel bir hata oluştu.",
		})
	}

	// Her şey başarılıysa, hesaplanan rotaları 200 OK koduyla JSON olarak döndür
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    routes,
		"count":   len(routes),
	})
}

type TelemetryRequest struct {
	Temperature   float64 `json:"temperature"`
	GasLevelPPM   float64 `json:"gas_level_ppm"`
	BatteryStatus float64 `json:"battery_status"`
}

// UpdateTelemetry, çöp kutusunun donanımsal durumunu günceller.
func (h *RouteHandler) UpdateTelemetry(c *fiber.Ctx) error {
	id := c.Params("id") // URL'den container ID'sini alıyoruz (Örn: /containers/123/telemetry)

	var req TelemetryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Geçersiz telemetri verisi.",
		})
	}

	// Repository üzerinden veritabanını güncelle
	if err := h.repo.UpdateTelemetry(c.Context(), id, req.Temperature, req.GasLevelPPM, req.BatteryStatus); err != nil {
		log.Printf("Telemetri güncelleme hatası: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Sunucu hatası.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Telemetri güncellendi",
	})
}
