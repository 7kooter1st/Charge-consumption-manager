package dto

type NumberOfUseCasesResponse struct {
	ConsumptionId    uint  `json:"consumption_id"`
	NumberOfUseCases int64 `json:"use_cases_in_consumption"`
}

type ConsumptionsResponse struct {
	ID          uint   `json:"id"`
	Status      string `json:"status"`
	Creator     string `json:"creator"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	ModeratedAt int64  `json:"moderated_at"`
	Moderator   string `json:"moderator"`
	TotalPower  uint   `json:"total_power"`
	UserID      uint   `json:"user_id"`
}

type OneConsumptionResponse struct {
	ID         uint                             `json:"id"`
	TotalPower uint                             `json:"total_power"`
	Status     string                           `json:"status"`
	UserID     uint                             `json:"user_id"`
	UseCases   []UseCaseFromConsumptionResponse `json:"use_cases"`
}

type ChangeConsumptionDuration struct {
	UseCaseID uint `json:"use_case_id" binding:"required"`
	Duration  uint `json:"duration" binding:"required"`
}

type ModeratorAction struct {
	Action string `json:"action" binding:"required"`
	// TotalConsumption float64 `json:"total_consumption" binding:"required"`
	ModeratorID uint `json:"moderator_id" binding:"required"`
}
