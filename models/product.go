package models

type Product struct {
	ID    int  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	UpdateCount int `json:"update_count"`
	UpdatedAt string `json:"updated_at"`
}
