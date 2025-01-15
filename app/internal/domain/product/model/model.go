package model

import "time"

type Product struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	ImageId       *string    `json:"image_id"`
	Price         int64      `json:"price"`
	CurrencyId    int32      `json:"currency_id"`
	Rating        int32      `json:"rating"`
	CategoryId    int32      `json:"category_id"`
	Specification *string    `json:"specification"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
