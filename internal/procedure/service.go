package procedure

import (
	"context"
	"nh-be/constant"
	"nh-be/internal/permission"
	"nh-be/utils"
	"slices"

	"github.com/google/uuid"
)

type Service interface {
	GetAllProcedures(ctx context.Context, search string, offset, limit int) ([]ProcedureListResponseDto, int64, error)
	// FindByID(ctx context.Context, id uuid.UUID) (*Procedure, error)
	// Create(ctx context.Context, procedure *Procedure) error
	// Update(ctx context.Context, procedure *Procedure) error
	// Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repository        Repository
	permissionService permission.Service
}

func NewService(repository Repository, permissionService permission.Service) Service {
	return &service{repository: repository, permissionService: permissionService}
}

func (s *service) CanViewProcedure(ctx context.Context, id uuid.UUID) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ViewExperiment) && !slices.Contains(userPerm, constant.ManageExperiment) {
		return constant.ErrForbidViewProcedure
	}

	return nil
}

func (s *service) GetAllProcedures(ctx context.Context, search string, offset, limit int) ([]ProcedureListResponseDto, int64, error) {
	permErr := s.CanViewProcedure(ctx, uuid.Nil)
	if permErr != nil {
		return nil, 0, permErr
	}
	procedures, length, repoErr := s.repository.FindAll(ctx, search, offset, limit, true)
	if repoErr != nil {
		return nil, 0, repoErr
	}

	return MapProceduresToDto(procedures), length, nil
}

// func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*Procedure, error) {
// 	return s.repository.FindByID(ctx, id)
// }

// func (s *service) Create(ctx context.Context, procedure *Procedure) error {
// 	return s.repository.Create(ctx, procedure)
// }

// func (s *service) Update(ctx context.Context, procedure *Procedure) error {
// 	return s.repository.Update(ctx, procedure)
// }

// func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
// 	return s.repository.Delete(ctx, id)
// }
