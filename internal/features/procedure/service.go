package procedure

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/ctxutil"
	"nh-be/internal/utils/timeutil"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAllProcedures(ctx context.Context, search string, offset, limit int) ([]ProcedureListResponseDto, int64, error)
	GetProcedureByID(ctx context.Context, id uuid.UUID) (*ProcedureResponseDto, error)
	CreateProcedure(ctx context.Context, procedure *CreateProcedureDto) error
	UpdateProcedure(ctx context.Context, id uuid.UUID, procedure *UpdateProcedureDto) error
	UpdateProcedureStep(ctx context.Context, procedureId uuid.UUID, steps []UpdateProcedureStepInput) error
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

func (s *service) UpdateProcedureStep(
	ctx context.Context,
	procedureId uuid.UUID,
	stepInput []UpdateProcedureStepInput,
) error {
	permErr := s.CanManageProcedure(ctx, procedureId, constant.Update)
	if permErr != nil {
		return permErr
	}

	transactionErr := s.repository.WithTransaction(ctx, func(repository Repository) error {
		existingSteps, getErr := repository.GetStepIDsByProcID(ctx, procedureId)
		if getErr != nil {
			return getErr
		}

		existingStepIds := make(map[uuid.UUID]int)

		for _, step := range existingSteps {
			existingStepIds[step.ID] = step.Version
		}

		incomingIDs := make(map[uuid.UUID]bool)
		now := timeutil.TimePtr(time.Now())

		for _, input := range stepInput {
			if input.ID == uuid.Nil {
				step := MapCreateProcStepInputToProcStep(&input)
				step.ProcedureID = procedureId

				if err := repository.CreateProcedureStep(ctx, step); err != nil {
					return err
				}
				continue
			}

			incomingIDs[input.ID] = true

			version, exists := existingStepIds[input.ID]
			if !exists {
				return constant.ErrProcedureStepNotFound
			}

			step := MapUpdateProcStepInputToProcStep(&input, now)
			step.Version = version

			if err := repository.UpdateProcedureStep(ctx, input.ID, procedureId, step); err != nil {
				return err
			}

		}

		for id := range existingStepIds {
			if !incomingIDs[id] {
				if err := repository.DeleteProcedureStep(ctx, id, procedureId); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if transactionErr != nil {
		return transactionErr
	}

	return nil
}

// func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
// 	return s.repository.Delete(ctx, id)
// }
