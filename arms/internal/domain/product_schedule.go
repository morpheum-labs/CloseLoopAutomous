package domain

import "time"

// ProductSchedule gates whether autopilot cadence ticks apply to a product (Mission Control–style product_schedules row).
// No row in the store means “enabled” for backward compatibility.
type ProductSchedule struct {
	ProductID ProductID
	Enabled   bool
	SpecJSON  string
	UpdatedAt time.Time
}
