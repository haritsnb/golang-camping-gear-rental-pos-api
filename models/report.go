package models

type RevenueReport struct {
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
	TotalRentalRevenue   float64 `json:"total_rental_revenue"`
	TotalPenaltyRevenue  float64 `json:"total_penalty_revenue"`
	TotalGrossRevenue    float64 `json:"total_gross_revenue"`
	TotalMaintenanceCost float64 `json:"total_maintenance_cost"`
	TotalNetProfit       float64 `json:"total_net_profit"`
	TotalTransactions    int     `json:"total_transactions"`
}

type RentalSummaryReport struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	TotalBooked      int    `json:"total_booked"`
	TotalActive      int    `json:"total_active"`
	TotalReturned    int    `json:"total_returned"`
	TotalCompleted   int    `json:"total_completed"`
	TotalCancelled   int    `json:"total_cancelled"`
	TotalTransaction int    `json:"total_transaction"`
}

type TopProductReport struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryName string  `json:"category_name"`
	BrandName    string  `json:"brand_name"`
	TotalRented  int     `json:"total_rented"`
	TotalRevenue float64 `json:"total_revenue"`
}

type InventoryReport struct {
	TotalAssetsValue float64 `json:"total_assets_value"`
	TotalUnits       int     `json:"total_units"`
	AvailableUnits   int     `json:"available_units"`
	BookedUnits      int     `json:"booked_units"`
	RentedUnits      int     `json:"rented_units"`
	MaintenanceUnits int     `json:"maintenance_units"`
	LostUnits        int     `json:"lost_units"`
}
