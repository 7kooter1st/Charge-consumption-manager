package dto

type UseCaseFromConsumptionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Consumption uint   `json:"consumption"`
	Image       string `json:"image"`
	IsDelete    bool   `json:"is_deleted"`
	Duration    uint   `json:"duration"` // длительность использования в минутах
}

type AddUseCase struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Consumption uint   `json:"consumption" binding:"required"`
}

type ChangeUseCase struct {
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Consumption uint   `json:"consumption,omitempty"`
}

type AddImageRequest struct {
	ImageURL string `json:"image_url" binding:"required,url"`
}
