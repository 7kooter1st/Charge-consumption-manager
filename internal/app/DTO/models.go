package dto

import "time"

type ConsumptionFilter struct {
	Status     string
	Start_date time.Time
	End_date   time.Time
}
