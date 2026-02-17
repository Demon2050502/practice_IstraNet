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

func (s *ApplicationService) DeleteUserApplication(ctx context.Context, userID, appID int64) error {
	return s.repos.DeleteUserApplication(ctx, userID, appID)
}

func (s *ApplicationService) UpdateUserApplication(ctx context.Context, userID int64, in dto.ChangeApplicationRequest) error {
	return s.repos.UpdateUserApplication(ctx, userID, in)
}
