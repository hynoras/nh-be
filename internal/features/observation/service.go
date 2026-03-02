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
	CreateObservation(
		ctx context.Context,
		expId, procStepId uuid.UUID,
		input CreateObservationInput,
	) (CreatedObservationResponseDto, error)
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

	mappedObs := MapObservationsMetadataToDto(obs)

	return mappedObs, length, nil
}

func (s *service) CreateObservation(
	ctx context.Context,
	expId, procStepId uuid.UUID,
	input CreateObservationInput,
) (CreatedObservationResponseDto, error) {
	permErr := s.CanCreateObservation(ctx)
	if permErr != nil {
		return CreatedObservationResponseDto{}, permErr
	}

	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return CreatedObservationResponseDto{}, err
	}

	observation := MapCreateInputToObservation(input, userId, expId, procStepId)
	createdObs, createErr := s.observationRepo.CreateObservation(ctx, expId, procStepId, observation)
	if createErr != nil {
		return CreatedObservationResponseDto{}, createErr
	}

	mappedObs := MapObsToCreatedObsResponseDto(createdObs)

	return mappedObs, nil
}
