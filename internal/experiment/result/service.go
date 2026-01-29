package result

import (
	"context"
	"nh-be/constant"
	"nh-be/internal/experiment/root"
	"nh-be/internal/permission"
	"nh-be/utils"
	"slices"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetResultByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error)
	CreateResult(ctx context.Context, dto *CreateResultDto) error
	UpdateResult(ctx context.Context, experimentID uuid.UUID, dto *UpdateResultDto) error
}

type service struct {
	resultRepo        Repository
	experimentRepo    root.Repository
	permissionService permission.Service
}

func NewService(resultRepo Repository, experimentRepo root.Repository, permissionService permission.Service) Service {
	return &service{
		resultRepo:        resultRepo,
		experimentRepo:    experimentRepo,
		permissionService: permissionService,
	}
}

func (s *service) GetResultByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error) {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(userPerm, constant.ViewExperiment) && !slices.Contains(userPerm, constant.ManageExperiment) {
		return nil, ErrForbidViewExperimentResult
	}

	// Verify experiment exists
	_, err = s.experimentRepo.FindByID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	result, err := s.resultRepo.FindByExperimentID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) CreateResult(ctx context.Context, dto *CreateResultDto) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidCreateExperimentResult
	}

	// Parse experiment ID
	experimentID, err := uuid.Parse(dto.ExperimentID)
	if err != nil {
		return err
	}

	// Verify experiment exists
	_, err = s.experimentRepo.FindByID(ctx, experimentID)
	if err != nil {
		return err
	}

	// Check if result already exists for this experiment
	existingResult, err := s.resultRepo.FindByExperimentID(ctx, experimentID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existingResult != nil {
		return ErrExperimentResultAlreadyExists
	}

	result := &ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    experimentID,
		Outcome:         Outcome(dto.Outcome),
		Summary:         dto.Summary,
		OutcomeReason:   dto.OutcomeReason,
		ConfidenceLevel: ConfidenceLevel(dto.ConfidenceLevel),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return s.resultRepo.Create(ctx, result)
}

func (s *service) UpdateResult(ctx context.Context, experimentID uuid.UUID, dto *UpdateResultDto) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageExperiment) {
		return ErrForbidUpdateExperimentResult
	}

	// Verify experiment exists
	_, err = s.experimentRepo.FindByID(ctx, experimentID)
	if err != nil {
		return err
	}

	// Check if result exists
	_, err = s.resultRepo.FindByExperimentID(ctx, experimentID)
	if err != nil {
		return err
	}

	result := &ExperimentResult{
		Outcome:         Outcome(dto.Outcome),
		Summary:         dto.Summary,
		OutcomeReason:   dto.OutcomeReason,
		ConfidenceLevel: ConfidenceLevel(dto.ConfidenceLevel),
		UpdatedAt:       time.Now(),
	}

	return s.resultRepo.Update(ctx, experimentID, result)
}
