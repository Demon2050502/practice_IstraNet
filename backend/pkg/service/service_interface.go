package service

import (
	"practice_IstraNet/pkg/repository"
)


type Service struct {

}

func NewService(repos *repository.Repository) *Service {
	return &Service{
	}
}