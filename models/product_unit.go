package models

type ProductUnit struct {
	ID           int     `json:"id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name,omitempty"`
	UnitCode     string  `json:"unit_code"`
	SerialNumber *string `json:"serial_number,omitempty"`
	Condition    string  `json:"condition"`
	Status       string  `json:"status"` // available, booked, rented, maintenance, lost
	Notes        *string `json:"notes,omitempty"`
	AuditTrail
}

type ProductUnitCreateDTO struct {
	ProductID    int     `json:"product_id" binding:"required"`
	UnitCode     string  `json:"unit_code" binding:"required"`
	SerialNumber *string `json:"serial_number"`
	Condition    string  `json:"condition"`
	Notes        *string `json:"notes"`
}

type ProductUnitUpdateDTO struct {
	ProductID    *int    `json:"product_id"`
	UnitCode     *string `json:"unit_code"`
	SerialNumber *string `json:"serial_number"`
	Condition    *string `json:"condition"`
	Status       *string `json:"status"`
	Notes        *string `json:"notes"`
}
