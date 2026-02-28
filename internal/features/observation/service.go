package observation

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/ctxutil"
	"slices"

	"github.com/google/uuid"
)

type Service interface {
	GetAllObservations(
		ctx context.Context,
		expId, procId uuid.UUID,
		offset, limit int,
		sortBy *string,
		sortOrder *constant.Order,
	) ([]ObservationsResponseDto, int64, error)
}

type service struct {
	observationRepo   Repository
	permissionService permission.Service
}

func NewService(observationRepo Repository, permissionService permission.Service) Service {
	return &service{
		observationRepo:   observationRepo,
		permissionService: permissionService,
	}
}

func (s *service) CanViewObservation(ctx context.Context) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ViewExperiment) && !slices.Contains(userPerm, constant.ManageExperiment) {
		return constant.ErrForbidViewObservation
	}

	return nil
}

func (s *service) CanCreateObservation(ctx context.Context) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return constant.ErrForbidCreateObservation
	}

	return nil
}

func (s *service) GetAllObservations(
	ctx context.Context,
	expId, procId uuid.UUID,
	offset, limit int,
	sortBy *string,
	sortOrder *constant.Order,
) ([]ObservationsResponseDto, int64, error) {
	permErr := s.CanViewObservation(ctx)
	if permErr != nil {
		return nil, 0, permErr
	}

	obs, length, getObsErr := s.observationRepo.GetAllObsByExpIDAndProcID(ctx, expId, procId, offset, limit, sortBy, sortOrder)
	if getObsErr != nil {
		return nil, 0, getObsErr
	}

	mappedObs := MapObservationsToDto(obs)

	return mappedObs, length, nil
}
