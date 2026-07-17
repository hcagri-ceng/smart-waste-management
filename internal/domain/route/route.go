package route

import "time"

type Container struct {
	ID               string    `json:"id"`
	Latitude         float64   `json:"latitude"`           // Enlem (GPS)
	Longitude        float64   `json:"longitude"`          // Boylam (GPS)
	Capacity         float64   `json:"capacity"`           // Kutunun maksimum hacmi veya ağırlık kapasitesi
	CurrentFillLevel float64   `json:"current_fill_level"` // Anlık doluluk oranı (%)
	Temperature      float64   `json:"temperature"`        // Kutunun iç sıcaklığı (Koku/bakteri reaksiyon hızını etkiler)
	GasLevelPPM      float64   `json:"gas_level_ppm"`      // Metan/Kötü koku salınım miktarı (Sensörden gelen veri)
	BatteryStatus    float64   `json:"battery_status"`     // IoT cihazının pil seviyesi (%)
	LastEmptiedAt    time.Time `json:"last_emptied_at"`    // En son ne zaman boşaltıldığı (Çok uzun süre bekleyen kutulara öncelik vermek için)
}
