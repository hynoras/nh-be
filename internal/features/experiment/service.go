package experiment

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/ctxutil"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAllExperiments(ctx context.Context, search string, page, pageSize int) ([]Experiment, int64, error)
	GetExperimentByID(ctx context.Context, id uuid.UUID) (*Experiment, error)
	CreateExperiment(ctx context.Context, dto *CreateExperimentDto) error
	UpdateExperiment(ctx context.Context, id uuid.UUID, dto *UpdateExperimentDto) error
	UpdateExperimentStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus) error
	DeleteExperiment(ctx context.Context, id uuid.UUID) error
}

type service struct {
	experimentRepo    Repository
	permissionService permission.Service
}

func NewService(experimentRepo Repository, permissionService permission.Service) Service {
	return &service{
		experimentRepo:    experimentRepo,
		permissionService: permissionService,
	}
}

func (s *service) GetAllExperiments(ctx context.Context, search string, page, pageSize int) ([]Experiment, int64, error) {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	if !slices.Contains(userPerm, constant.ViewExperiment) && !slices.Contains(userPerm, constant.ManageExperiment) {
		return nil, 0, ErrForbidViewExperiments
	}

	experiments, length, err := s.experimentRepo.FindAll(ctx, search, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return experiments, length, nil
}

func (s *service) GetExperimentByID(ctx context.Context, id uuid.UUID) (*Experiment, error) {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(userPerm, constant.ViewExperiment) && !slices.Contains(userPerm, constant.ManageExperiment) {
		return nil, ErrForbidViewExperiment
	}

	experiment, err := s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return experiment, nil
}

func (s *service) CreateExperiment(ctx context.Context, dto *CreateExperimentDto) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidCreateExperiment
	}

	experiment := &Experiment{
		ID:          uuid.New(),
		Title:       dto.Title,
		Objective:   dto.Objective,
		Status:      ExperimentDraft,
		Type:        ExperimentType(dto.Type),
		CreatedByID: userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.experimentRepo.Create(ctx, experiment)
}

func (s *service) UpdateExperiment(ctx context.Context, id uuid.UUID, dto *UpdateExperimentDto) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidUpdateExperiment
	}

	// Check if experiment exists
	_, err = s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	experiment := &Experiment{
		Title:     dto.Title,
		Objective: dto.Objective,
		Type:      ExperimentType(dto.Type),
		UpdatedAt: time.Now(),
	}

	return s.experimentRepo.Update(ctx, id, experiment)
}

func (s *service) UpdateExperimentStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidUpdateExperiment
	}

	// Check if experiment exists
	var exp *Experiment
	exp, err = s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if exp.Status == status {
		return ErrAlreadyInTargetState
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

func (s *service) DeleteExperiment(ctx context.Context, id uuid.UUID) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidDeleteExperiment
	}

	// Check if experiment exists
	_, err = s.experimentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.experimentRepo.Delete(ctx, id)
}
