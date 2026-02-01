package model

// Produk represents a product with optional category relationship
type Produk struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Harga      int       `json:"harga"`
	Stok       int       `json:"stok"`
	CategoryID *int      `json:"category_id,omitempty"`
	Category   *Category `json:"category,omitempty"`
}
