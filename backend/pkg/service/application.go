package service

import (
	"context"

	dbmodel "practice_IstraNet/pkg/DB_model"
	"practice_IstraNet/pkg/dto"
	repository "practice_IstraNet/pkg/repository"
)

type ApplicationService struct {
	repos repository.Applications
}

func NewApplicationService(repos repository.Applications) *ApplicationService {
	return &ApplicationService{repos: repos}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, userID int64, in dto.CreateApplicationRequest) (dto.ApplicationResponse, error) {
	app, err := s.repos.CreateApplication(
		ctx,
		userID,
		in.Title,
		in.Description,
		in.PriorityCode,
		in.CategoryID,
		in.ContactPhone,
		in.ContactAddress,
	)
	if err != nil {
		return dto.ApplicationResponse{}, err
	}
	return mapApplicationDB(app), nil
}

func (s *ApplicationService) GetAllApplications(ctx context.Context) (dto.ApplicationsListResponse, error) {
	items, err := s.repos.GetAllApplications(ctx)
	if err != nil {
		return dto.ApplicationsListResponse{}, err
	}

	out := dto.ApplicationsListResponse{Items: make([]dto.ApplicationResponse, 0, len(items))}
	for _, a := range items {
		out.Items = append(out.Items, mapApplicationDB(a))
	}
	return out, nil
}

func mapApplicationDB(a dbmodel.ApplicationDB) dto.ApplicationResponse {
	var out dto.ApplicationResponse

	out.ID = a.ID
	out.Title = a.Title
	out.Description = a.Description

	out.Status.Code = a.StatusCode
	out.Status.Name = a.StatusName

	out.Priority.Code = a.PriorityCode
	out.Priority.Name = a.PriorityName
	out.Priority.Weight = a.PriorityWeight

	if a.CategoryID != nil && a.CategoryName != nil {
		out.Category = &struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}{ID: *a.CategoryID, Name: *a.CategoryName}
	}

	out.CreatedBy.ID = a.CreatedByID
	out.CreatedBy.Name = a.CreatedByName

	if a.AssignedToID != nil && a.AssignedToName != nil {
		out.AssignedTo = &struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}{ID: *a.AssignedToID, Name: *a.AssignedToName}
	}

	out.ContactPhone = a.ContactPhone
	out.ContactAddress = a.ContactAddress

	out.CreatedAt = a.CreatedAt
	out.UpdatedAt = a.UpdatedAt
	out.ClosedAt = a.ClosedAt

	return out
}
