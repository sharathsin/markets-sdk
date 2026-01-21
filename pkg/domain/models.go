package domain

import (
	"time"
)

// AssetType represents the type of financial asset
type AssetType string

const (
	AssetTypeStock  AssetType = "STOCK"
	AssetTypeCrypto AssetType = "CRYPTO"
)

// Quote represents a simplified pricing quote for an asset
type Quote struct {
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Change24h   float64   `json:"change_24h,omitempty"`
	Volume      float64   `json:"volume,omitempty"`
	LastUpdated time.Time `json:"last_updated"`
	Source      string    `json:"source"`
}
