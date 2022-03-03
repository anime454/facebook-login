package service

type facebookService struct{}

func NewFacebookService() facebookService {
	return facebookService{}
}

func (fbService facebookService) Callback() error {
	return nil
}
