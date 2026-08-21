package models

type RentalItem struct {
	ID             int     `json:"id"`
	RentalID       int     `json:"rental_id"`
	ProductUnitID  int     `json:"product_unit_id"`
	UnitCode       string  `json:"unit_code,omitempty"`
	ProductName    string  `json:"product_name,omitempty"`
	DailyRate      float64 `json:"daily_rate"`
	DurationDays   int     `json:"duration_days"`
	Subtotal       float64 `json:"subtotal"`
	ConditionOut   string  `json:"condition_out"`
	ConditionIn    *string `json:"condition_in,omitempty"`
	ItemPenaltyFee float64 `json:"item_penalty_fee"`
	ReturnNotes    *string `json:"return_notes,omitempty"`
	AuditTrail
}
