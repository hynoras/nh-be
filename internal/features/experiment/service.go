package experiment

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/procedure"
	"nh-be/internal/utils/authutil"
	"nh-be/internal/utils/ctxutil"
	"nh-be/internal/utils/timeutil"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAllExperiments(ctx context.Context,
		sortBy,
		search string,
		sortOrder constant.Order,
		experimentStatus *ExperimentStatus,
		experimentType *ExperimentType,
		page, pageSize int) ([]ExperimentsResponseDto, int64, error)
	GetExperimentDetail(ctx context.Context, identifier string) (*ExperimentResponseDto, error)
	CreateExperiment(ctx context.Context, dto *CreateExperimentDto) error
	UpdateExperiment(ctx context.Context, id uuid.UUID, dto *UpdateExperimentDto) error
	UpdateExperimentStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus) error
	AssignProcedureToExperiment(ctx context.Context, experimentId uuid.UUID, procedureId uuid.UUID, version int) error
	DeleteExperiment(ctx context.Context, id uuid.UUID) error
}

type service struct {
	experimentRepo    Repository
	permissionService permission.Service
	procedureService  procedure.Service
}

func NewService(experimentRepo Repository, permissionService permission.Service, procedureService procedure.Service) Service {
	return &service{
		experimentRepo:    experimentRepo,
		permissionService: permissionService,
		procedureService:  procedureService,
	}
}

func (s *service) CanManageExperiment(ctx context.Context, id uuid.UUID, action constant.ManageAction) error {
	var forbidErr error
	switch action {
	case constant.Create:
		forbidErr = ErrForbidCreateExperiment
	case constant.Update:
		forbidErr = ErrForbidUpdateExperiment
	case constant.Delete:
		forbidErr = ErrForbidDeleteExperiment
	default:
		forbidErr = ErrForbidManageExperiment
	}

	return authutil.RequirePermission(ctx, s.permissionService, forbidErr, constant.ManageExperiment)
}

func (s *service) GetAllExperiments(ctx context.Context,
	sortBy, search string,
	sortOrder constant.Order,
	experimentStatus *ExperimentStatus,
	experimentType *ExperimentType,
	page, pageSize int) ([]ExperimentsResponseDto, int64, error) {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidViewExperiments, constant.ViewExperiment, constant.ManageExperiment); err != nil {
		return nil, 0, err
	}

	userId, getUserIdErr := ctxutil.GetUserIdFromContext(ctx)
	if getUserIdErr != nil {
		return nil, 0, getUserIdErr
	}

	var queryExperiments []ExperimentsQueryDto
	var count int64

	//NOTE: for now retrieve count that is belong to current logged in user
	count, countErr := s.experimentRepo.CountExperiments(ctx, &userId)
	if countErr != nil {
		return nil, 0, countErr
	}

	//NOTE: for now retrieve experiments that is belong to current logged in user
	queryExperiments, queryErr := s.experimentRepo.FindAllExperiments(ctx, sortBy, search, sortOrder, experimentStatus, experimentType, &userId, page, pageSize)
	if queryErr != nil {
		return nil, 0, queryErr
	}

	experiments := MapExperimentsQueryToDto(queryExperiments)
	return experiments, count, nil
}

func (s *service) GetExperimentDetail(ctx context.Context, identifier string) (*ExperimentResponseDto, error) {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidViewExperiment, constant.ViewExperiment, constant.ManageExperiment); err != nil {
		return nil, err
	}

	userId, getUserIdErr := ctxutil.GetUserIdFromContext(ctx)
	if getUserIdErr != nil {
		return nil, getUserIdErr
	}

	experiment, err := s.experimentRepo.FindByIdentifierAndCreatedBy(ctx, identifier, userId)
	if err != nil {
		return nil, err
	}

	mappedExperiment := MapExperimentToDto(*experiment)
	return &mappedExperiment, nil
}

func (s *service) CreateExperiment(ctx context.Context, dto *CreateExperimentDto) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidCreateExperiment, constant.ManageExperiment); err != nil {
		return err
	}

	experiment := &Experiment{
		ID:          uuid.New(),
		Title:       dto.Title,
		Objective:   dto.Objective,
		Status:      ExperimentDraft,
		Type:        ExperimentType(dto.Type),
		CreatedByID: userId,
		CreatedAt:   time.Now(),
	}

	return s.experimentRepo.Create(ctx, experiment)
}

func (s *service) UpdateExperiment(ctx context.Context, id uuid.UUID, dto *UpdateExperimentDto) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidUpdateExperiment, constant.ManageExperiment); err != nil {
		return err
	}

	// Check if experiment exists
	_, err := s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	experiment := &Experiment{
		Title:     dto.Title,
		Objective: dto.Objective,
		Type:      ExperimentType(dto.Type),
		UpdatedAt: timeutil.TimePtr(time.Now()),
	}

	return s.experimentRepo.Update(ctx, id, experiment)
}

func (s *service) UpdateExperimentStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidUpdateExperiment, constant.ManageExperiment); err != nil {
		return err
	}

	// Check if experiment exists
	var exp *Experiment
	exp, err := s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if exp.Status == status {
		return ErrExperimentAlreadyInTargetState
	}

	// validate status
	if exp.Status == ExperimentDraft && status != ExperimentPlanning {
		return ErrStatusTransitionFromDraftToPlanning
	}
	if exp.Status == ExperimentPlanning && status != ExperimentRunning {
		return ErrStatusTransitionFromPlanningToRunning
	}
	if exp.Status == ExperimentRunning && (status != ExperimentCompleted && status != ExperimentAborted) {
		return ErrStatusTransitionFromRunningToCompletedOrAborted
	}

	return s.experimentRepo.UpdateStatus(ctx, id, status, exp.Version)
}

func (s *service) AssignProcedureToExperiment(ctx context.Context, experimentId uuid.UUID, procedureId uuid.UUID, version int) error {
	permErr := s.CanManageExperiment(ctx, experimentId, constant.Update)
	if permErr != nil {
		return permErr
	}

	// Check if procedure exists
	_, getProcErr := s.procedureService.GetProcedureByID(ctx, procedureId)
	if getProcErr != nil {
		return getProcErr
	}

	// Check if experiment exists
	assignedProcedureId, getExpErr := s.experimentRepo.GetProcedureIDByID(ctx, experimentId)
	if getExpErr != nil {
		return getExpErr
	}

	if assignedProcedureId == procedureId {
		return ErrDuplicateProcedureAssignment
	}

	return s.experimentRepo.UpdateProcedureID(ctx, experimentId, procedureId, version)
}

func (s *service) DeleteExperiment(ctx context.Context, id uuid.UUID) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidDeleteExperiment, constant.ManageExperiment); err != nil {
		return err
	}

	// Check if experiment exists
	_, err := s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.experimentRepo.Delete(ctx, id)
}
