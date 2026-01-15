package views

type HealthView struct {
	Message string `json:"message"`
}

func NewHealthView(message string) *HealthView {
	return &HealthView{
		Message: message,
	}
}
