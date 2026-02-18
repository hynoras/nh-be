package procedure

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/procedure"
	"nh-be/internal/utils/timeutil"
	"time"

	"github.com/google/uuid"
)

func ContextWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), constant.CtxUserId, userID)
}

func TestProcedureList() []procedure.Procedure {
	timeNow := time.Now()
	return []procedure.Procedure{
		{
			ID:          uuid.New(),
			Title:       "Test Procedure",
			Description: "Test Procedure Description",
			Version:     1,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   &timeNow,
		},
		{
			ID:          uuid.New(),
			Title:       "Test Procedure 2",
			Description: "Test Procedure Description 2",
			Version:     2,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   timeutil.TimePtr(time.Now()),
		},
	}
}

func TestProcedureDetail() procedure.Procedure {
	return procedure.Procedure{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Title:       "Test Procedure",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
		CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
	}
}

func TestStepList() []procedure.ProcedureStep {
	return []procedure.ProcedureStep{
		{
			ID:          uuid.MustParse("33333333-1234-1234-1234-444433332222"),
			ProcedureID: uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			Index:       1,
			Title:       "Test Step",
			Description: "Test Step Description",
			CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
		},
		{
			ID:          uuid.MustParse("33333333-1234-1234-1234-555533332222"),
			ProcedureID: uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			Index:       2,
			Title:       "Test Step 2",
			Description: "Test Step Description 2",
			CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
		},
	}
}

func TestProcedureDetailWithRelations() procedure.Procedure {
	procID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	return procedure.Procedure{
		ID:          procID,
		Title:       "Test Procedure",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
		CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
		Steps:       TestStepList(),
	}
}

func CreateValidProcedure() *procedure.Procedure {
	return &procedure.Procedure{
		ID:          uuid.New(),
		Title:       "Test Procedure",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
	}
}

func CreateProcedureWithSteps() *procedure.Procedure {
	procID := uuid.New()
	return &procedure.Procedure{
		ID:          procID,
		Title:       "Test Procedure with Steps",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
		Steps: []procedure.ProcedureStep{
			{
				ID:          uuid.New(),
				ProcedureID: procID,
				Index:       1,
				Title:       "Step 1",
				Description: "Description 1",
				StepType:    "manual",
			},
			{
				ID:          uuid.New(),
				ProcedureID: procID,
				Index:       2,
				Title:       "Step 2",
				Description: "Description 2",
				StepType:    "manual",
			},
			{
				ID:          uuid.New(),
				ProcedureID: procID,
				Index:       3,
				Title:       "Step 3",
				Description: "Description 3",
				StepType:    "manual",
			},
		},
	}
}

func CreateProcedureWithExperiments() *procedure.Procedure {
	procID := uuid.New()

	return &procedure.Procedure{
		ID:          procID,
		Title:       "Test Procedure with Experiments",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
	}
}

func CreateProcedureWithParent(parentID uuid.UUID) *procedure.Procedure {
	return &procedure.Procedure{
		ID:          uuid.New(),
		Title:       "Child Procedure",
		Description: "Child Procedure Description",
		Version:     2,
		ParentID:    &parentID,
	}
}

func TestCreateProcedureDto() *procedure.CreateProcedureDto {
	return &procedure.CreateProcedureDto{
		Title:       "Test Procedure",
		Description: "Test Procedure Description",
		Steps: []procedure.CreateStepDto{
			{
				Index:       1,
				Title:       "Step 1",
				Description: "Step 1 Description",
				StepType:    "action",
			},
			{
				Index:       2,
				Title:       "Step 2",
				Description: "Step 2 Description",
				StepType:    "wait",
			},
		},
		ExperimentAssignments: []procedure.CreateExperimentAssignmentDto{
			{
				ID: uuid.New().String(),
			},
		},
	}
}

func TestUpdateProcedureDto() *procedure.UpdateProcedureDto {
	return &procedure.UpdateProcedureDto{
		Title:       "Updated Procedure Title",
		Description: "Updated Procedure Description",
		Version:     1,
	}
}
