package dto

type ChangeUseCaseDurationRequest struct {
	Duration uint `json:"duration" binding:"required"`
}
