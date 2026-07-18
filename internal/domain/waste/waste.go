package waste

import (
	"time"
)

type Waste struct {
	ID              string    `json:"id"`                 // Güvenlik ve dağıtık sistemler için UUID (Universally Unique Identifier) kullanmak standarttır.
	ContainerID     string    `json:"container_id"`       // Atığın bırakıldığı akıllı çöp kutusunun veya konteynerin ID'si.
	Type            string    `json:"type"`               // Plastik, Kağıt, Cam, Elektronik vb.
	Weight          float64   `json:"weight"`             // Kilogram cinsinden ağırlık (sensörlerden gelebilir).
	CarbonFootprint float64   `json:"carbon_footprint"`   // Bu atığın geri dönüştürülmesiyle kurtarılan veya hesaplanan karbon ayak izi miktarı.
	CreatedAt       time.Time `json:"created_at"`         // Verinin sisteme işlendiği an.
	FillLevelAtDrop float64   `json:"fill_level_at_drop"` // Atığın bırakıldığı andaki çöp kutusunun doluluk seviyesi (sensörlerden gelebilir).
}
