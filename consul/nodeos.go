package consul

import (
	"time"
	"context"
)

func (s *Service) ResumeProduction() error {

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	s.PPMan.ResumeProduction(&ctx)

	if s.PPMan.Err != nil {

		return s.PPMan.Err
	}

	return nil
}

func (s *Service) PauseProduction() error {

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	s.PPMan.PauseProduction(&ctx)

	if s.PPMan.Err != nil {

		return s.PPMan.Err
	}

	return nil
}

func (s *Service) ProductionState() error {

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	s.PPMan.ProductionState(&ctx)

	if s.PPMan.Err != nil {

		return s.PPMan.Err
	}

	return nil
}

func (s *Service) Check() (bool, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	result := s.PPMan.CheckIfApiIsReachable(&ctx)

	if s.PPMan.Err != nil {

		return false, s.PPMan.Err
	}

	return result, nil
}