package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// Sahadaki donanımları temsil eden UUID'ler (Kendi veritabanından aldığın geçerli ID'leri buraya yaz)
var containerIDs = []string{
	"0e05ac73-d420-4755-955b-f3f5cb09613a", // Eskiler
	"1124e3d4-30ff-4b54-8228-1e64a7803018", // Eskiler
	"640f0a99-d018-4beb-86a5-2dc75b0c0c5e", // Eskiler
	"7d2b3c4a-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
	"8e3c4d5b-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
	"9f4d5e6c-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
	"a05e6f7d-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
	"b16f7a8e-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
	"c27a8b9f-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
	"d38b9c0a-1e2f-3a4b-5c6d-7e8f9a0b1c2d",
	"e49c0d1b-2f3a-4b5c-6d7e-8f9a0b1c2d3e",
	"f50d1e2c-3a4b-5c6d-7e8f-9a0b1c2d3e4f",
	"0a1e2f3d-4b5c-6d7e-8f9a-0b1c2d3e4f50",
	"1b2f3a4e-5c6d-7e8f-9a0b-1c2d3e4f5061",
	"2c3a4b5f-6d7e-8f9a-0b1c-2d3e4f506172",
}

// TelemetryPayload, API'a gönderilecek DTO yapısı
type TelemetryPayload struct {
	Temperature   float64 `json:"temperature"`
	GasLevelPPM   float64 `json:"gas_level_ppm"`
	BatteryStatus float64 `json:"battery_status"`
}

func main() {
	log.Println("IoT Donanım Simülasyonu başlatılıyor...")

	// Her bir çöp kutusunu bağımsız bir Goroutine (eşzamanlı process) olarak ayağa kaldır
	for _, id := range containerIDs {
		go simulateDevice(id)
	}

	// Ana uygulamanın kapanmasını engellemek için boş bir select (Sonsuz dinleme)
	select {}
}

// simulateDevice, sahada çalışan tek bir IoT kartının davranış modelidir.
func simulateDevice(containerID string) {
	// Sensör her 5 saniyede bir uyanıp merkeze veri gönderecek
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 1. Normal Şartlar Altındaki Sensör Verileri
		temp := 20.0 + rand.Float64()*10.0 // 20 - 30 derece arası normal sıcaklık
		gas := rand.Float64() * 5.0        // 0 - 5 ppm arası normal metan gazı
		battery := 100.0 - rand.Float64()*10.0

		// 2. Acil Durum (Kaos) Simülasyonu - %10 ihtimalle cihaz hata versin
		anomaly := rand.Intn(100)
		if anomaly < 5 {
			temp = 75.0 + rand.Float64()*15.0 // Sıcaklık birden 75-90 arasına fırlar (YANGIN!)
			log.Printf("[ALARM - YANGIN] Cihaz: %s", containerID[:8])
		} else if anomaly >= 5 && anomaly < 10 {
			gas = 30.0 + rand.Float64()*15.0 // Gaz birden 30-45 arasına fırlar (Sızıntı!)
			log.Printf("[ALARM - GAZ] Cihaz: %s", containerID[:8])
		}

		// Verileri virgülden sonra tek haneli olacak şekilde yuvarla
		payload := TelemetryPayload{
			Temperature:   math.Round(temp*10) / 10,
			GasLevelPPM:   math.Round(gas*10) / 10,
			BatteryStatus: math.Round(battery*10) / 10,
		}

		// 3. Veriyi Go Fiber API'ımıza gönder
		sendData(containerID, payload)
	}
}

// sendData, HTTP PUT isteğini gerçek bir donanım gibi backend'e iletir.
func sendData(id string, payload TelemetryPayload) {
	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("http://localhost:3000/api/v1/containers/%s/telemetry", id)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("İletişim Hatası (%s): %v\n", id[:8], err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Gönderildi -> Kutu: %s | Sıcaklık: %.1f | Gaz: %.1f", id[:8], payload.Temperature, payload.GasLevelPPM)
	}
}
