package models

import "strings"

type Receipt struct {
	Retailer     string `json:"retailer" validate:"required"`
	PurchaseDate string `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string `json:"purchaseTime" validate:"required"`
	Items        []Item `json:"items" validate:"required,dive"`
	Total        string `json:"total" validate:"required"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

// Sanitize trims strings to clean input
func (r *Receipt) Sanitize() {
	r.Retailer = strings.TrimSpace(r.Retailer)
	r.PurchaseDate = strings.TrimSpace(r.PurchaseDate)
	r.PurchaseTime = strings.TrimSpace(r.PurchaseTime)
	r.Total = strings.TrimSpace(r.Total)
	for i := range r.Items {
		r.Items[i].ShortDescription = strings.TrimSpace(r.Items[i].ShortDescription)
		r.Items[i].Price = strings.TrimSpace(r.Items[i].Price)
	}
}
