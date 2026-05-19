package result

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/authutil"
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
	if err := authutil.RequirePermission(ctx, s.permissionService, constant.ErrForbidViewExperimentResult, constant.ViewExperiment, constant.ManageExperiment); err != nil {
		return nil, err
	}

	result, err := s.resultRepo.FindByExperimentID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) CreateResult(ctx context.Context, experimentID uuid.UUID, dto *CreateResultDto) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, constant.ErrForbidCreateExperimentResult, constant.ManageExperiment); err != nil {
		return err
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
	if err := authutil.RequirePermission(ctx, s.permissionService, constant.ErrForbidUpdateExperimentResult, constant.ManageExperiment); err != nil {
		return err
	}

	// Use FindByIDAndExperimentID to validate both IDs and get current version
	_, err := s.resultRepo.FindByIDAndExperimentID(ctx, resultID, experimentID)
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
