package procedure

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/ctxutil"
	"slices"

	"github.com/google/uuid"
)

type Service interface {
	GetAllProcedures(ctx context.Context, search string, offset, limit int) ([]ProcedureListResponseDto, int64, error)
	GetProcedureByID(ctx context.Context, id uuid.UUID) (*ProcedureResponseDto, error)
	CreateProcedure(ctx context.Context, procedure *CreateProcedureDto) error
	UpdateProcedure(ctx context.Context, id uuid.UUID, procedure *UpdateProcedureDto) error
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
	userId, err := ctxutil.GetUserIdFromContext(ctx)
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

func (s *service) CanManageProcedure(ctx context.Context, id uuid.UUID, action constant.ManageAction) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		switch action {
		case constant.Create:
			return constant.ErrForbidCreateProcedure
		case constant.Update:
			return constant.ErrForbidUpdateProcedure
		case constant.Delete:
			return constant.ErrForbidDeleteProcedure
		}
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

func (s *service) GetProcedureByID(ctx context.Context, id uuid.UUID) (*ProcedureResponseDto, error) {
	permErr := s.CanViewProcedure(ctx, id)
	if permErr != nil {
		return nil, permErr
	}
	procedure, repoErr := s.repository.FindByID(ctx, id, true, true)
	if repoErr != nil {
		return nil, repoErr
	}
	mapToProcedure := MapProcedureToDto(procedure)
	return &mapToProcedure, nil
}

func (s *service) CreateProcedure(ctx context.Context, procedure *CreateProcedureDto) error {
	permErr := s.CanManageProcedure(ctx, uuid.Nil, constant.Create)
	if permErr != nil {
		return permErr
	}

	repoErr := s.repository.CreateProcedure(ctx, MapCreateDtoToProcedure(procedure))
	return repoErr
}

func (s *service) UpdateProcedure(ctx context.Context, id uuid.UUID, procedure *UpdateProcedureDto) error {
	permErr := s.CanManageProcedure(ctx, id, constant.Update)
	if permErr != nil {
		return permErr
	}
	repoErr := s.repository.UpdateProcedure(ctx, id, MapUpdateDtoToProcedure(procedure))
	return repoErr
}

// func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
// 	return s.repository.Delete(ctx, id)
// }
