package services

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (hs *HealthService) CheckHealth() (string, bool) {
	return "Healthy", true
}
