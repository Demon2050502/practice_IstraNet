package service

import (
	"context"
	"practice_IstraNet/pkg/dto"
)

func (s *ApplicationService) GetUserApplications(ctx context.Context, userID int64) (dto.UserAppsResponse, error) {
	apps, err := s.repos.GetUserApplications(ctx, userID)
	if err != nil {
		return dto.UserAppsResponse{}, err
	}

	var result dto.UserAppsResponse
	for _, a := range apps {
		result.Items = append(result.Items, dto.UserAppShort{
			ID:        a.ID,
			Title:     a.Title,
			Status:    a.StatusCode,
			Priority:  a.PriorityCode,
			CreatedAt: a.CreatedAt,
		})
	}

	return result, nil
}

func (s *ApplicationService) GetUserApplicationByID(ctx context.Context, userID, appID int64) (dto.UserAppDetailsResponse, error) {
	app, err := s.repos.GetUserApplicationByID(ctx, userID, appID)
	if err != nil {
		return dto.UserAppDetailsResponse{}, err
	}

	commentsDB, _ := s.repos.GetApplicationComments(ctx, appID)

	var comments []dto.AppComment
	for _, c := range commentsDB {
		comments = append(comments, dto.AppComment{
			ID:        c.ID,
			Author:    c.Author,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
		})
	}

	return dto.UserAppDetailsResponse{
		ApplicationResponse: mapApplicationDB(app),
		Comments:            comments,
	}, nil
}

func (s *ApplicationService) GetOperatorApplicationByID(ctx context.Context, appID int64) (dto.UserAppDetailsResponse, error) {
	app, err := s.repos.GetApplicationByID(ctx, appID)
	if err != nil {
		return dto.UserAppDetailsResponse{}, err
	}

	commentsDB, _ := s.repos.GetApplicationComments(ctx, appID)

	var comments []dto.AppComment
	for _, c := range commentsDB {
		comments = append(comments, dto.AppComment{
			ID:        c.ID,
			Author:    c.Author,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
		})
	}

	return dto.UserAppDetailsResponse{
		ApplicationResponse: mapApplicationDB(app),
		Comments:            comments,
	}, nil
}

func (s *ApplicationService) DeleteUserApplication(ctx context.Context, userID, appID int64) error {
	return s.repos.DeleteUserApplication(ctx, userID, appID)
}

func (s *ApplicationService) UpdateUserApplication(ctx context.Context, userID int64, in dto.ChangeApplicationRequest) error {
	return s.repos.UpdateUserApplication(ctx, userID, in)
}

func (s *ApplicationService) GetApplicationHistory(ctx context.Context, appID int64) (dto.ApplicationHistoryResponse, error) {
	if _, err := s.repos.GetApplicationByID(ctx, appID); err != nil {
		return dto.ApplicationHistoryResponse{}, err
	}

	items, err := s.repos.GetApplicationHistory(ctx, appID)
	if err != nil {
		return dto.ApplicationHistoryResponse{}, err
	}

	out := dto.ApplicationHistoryResponse{
		Items: make([]dto.ApplicationHistoryItem, 0, len(items)),
	}
	for _, it := range items {
		out.Items = append(out.Items, dto.ApplicationHistoryItem{
			ID:       it.ID,
			Action:   it.Action,
			Field:    it.Field,
			OldValue: it.OldValue,
			NewValue: it.NewValue,
			Actor: dto.HistoryActor{
				ID:   it.ActorID,
				Name: it.ActorName,
			},
			CreatedAt: it.CreatedAt,
		})
	}
	return out, nil
}

func (s *ApplicationService) TakeApplication(ctx context.Context, operatorID, appID int64) error {
	return s.repos.TakeApplication(ctx, operatorID, appID)
}

func (s *ApplicationService) ChangeApplicationStatus(ctx context.Context, operatorID int64, in dto.ChangeStatusRequest) error {
	return s.repos.ChangeApplicationStatus(ctx, operatorID, in)
}

func (s *ApplicationService) CloseApplication(ctx context.Context, operatorID int64, in dto.CloseApplicationRequest) error {
	return s.repos.CloseApplication(ctx, operatorID, in)
}
