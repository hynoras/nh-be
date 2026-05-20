package observation

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/authutil"
	"nh-be/internal/utils/ctxutil"

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
	return authutil.RequirePermission(ctx, s.permissionService, ErrForbidViewObservation, constant.ViewExperiment, constant.ManageExperiment)
}

func (s *service) CanCreateObservation(ctx context.Context) error {
	return authutil.RequirePermission(ctx, s.permissionService, ErrForbidCreateObservation, constant.ManageExperiment)
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
	createdObs, createErr := s.observationRepo.CreateObservation(ctx, observation)
	if createErr != nil {
		return CreatedObservationResponseDto{}, createErr
	}

	mappedObs := MapObsToCreatedObsResponseDto(createdObs)

	return mappedObs, nil
}
