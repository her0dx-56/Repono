package service

type HealthService struct{
}

func NewHealthService() *HealthService{
	return &HealthService{}
}

func (s *HealthService) Check() string{
	return "ok"
}