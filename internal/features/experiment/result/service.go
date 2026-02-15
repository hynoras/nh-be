package result

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
	GetResultByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error)
	CreateResult(ctx context.Context, experimentID uuid.UUID, dto *CreateResultDto) error
	UpdateResult(ctx context.Context, resultID uuid.UUID, experimentID uuid.UUID, dto *UpdateResultDto) error
}

type service struct {
	resultRepo        Repository
	permissionService permission.Service
}

func NewService(resultRepo Repository, permissionService permission.Service) Service {
	return &service{
		resultRepo:        resultRepo,
		permissionService: permissionService,
	}
}

func (s *service) GetResultByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error) {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
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

	result, err := s.resultRepo.FindByExperimentID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) CreateResult(ctx context.Context, experimentID uuid.UUID, dto *CreateResultDto) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
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

	result := &ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    experimentID,
		Outcome:         Outcome(dto.Outcome),
		Summary:         dto.Summary,
		OutcomeReason:   dto.OutcomeReason,
		ConfidenceLevel: ConfidenceLevel(dto.ConfidenceLevel),
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return s.resultRepo.Create(ctx, result)
}

func (s *service) UpdateResult(ctx context.Context, resultID uuid.UUID, experimentID uuid.UUID, dto *UpdateResultDto) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
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

	// Use FindByIDAndExperimentID to validate both IDs and get current version
	_, err = s.resultRepo.FindByIDAndExperimentID(ctx, resultID, experimentID)
	if err != nil {
		return err
	}

	// Build update fields with pointers
	fields := &UpdateFields{}
	if dto.Outcome != nil {
		outcome := Outcome(*dto.Outcome)
		fields.Outcome = &outcome
	}
	if dto.Summary != nil {
		fields.Summary = dto.Summary
	}
	if dto.OutcomeReason != nil {
		fields.OutcomeReason = dto.OutcomeReason
	}
	if dto.ConfidenceLevel != nil {
		confidenceLevel := ConfidenceLevel(*dto.ConfidenceLevel)
		fields.ConfidenceLevel = &confidenceLevel
	}

	return s.resultRepo.Update(ctx, resultID, experimentID, fields, dto.Version)
}
