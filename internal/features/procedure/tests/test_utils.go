package procedure

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/procedure"
	"nh-be/internal/utils/intutil"
	"nh-be/internal/utils/stringutil"
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

func TestProcedureStep() *procedure.ProcedureStep {
	return &procedure.ProcedureStep{
		ID:          uuid.MustParse("33333333-1234-1234-1234-444433332222"),
		ProcedureID: uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Index:       1,
		Title:       "Test Step",
		Description: "Test Step Description",
		StepType:    "manual",
		Version:     1,
		CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
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

// Fixed UUIDs for UpdateProcedureStep tests
var (
	ExistingStepID1 = uuid.MustParse("aaaaaaaa-1111-1111-1111-111111111111")
	ExistingStepID2 = uuid.MustParse("bbbbbbbb-2222-2222-2222-222222222222")
	ExistingStepID3 = uuid.MustParse("cccccccc-3333-3333-3333-333333333333")
)

func TestExistingStepMetadata() []procedure.StepMetadata {
	return []procedure.StepMetadata{
		{ID: ExistingStepID1, Version: 1},
		{ID: ExistingStepID2, Version: 2},
	}
}

func TestNewStepInput() procedure.UpdateProcedureStepInput {
	return procedure.UpdateProcedureStepInput{
		ID:          uuid.Nil,
		Index:       3,
		Title:       "New Step",
		Description: "New Step Description",
		StepType:    "action",
	}
}

func TestUpdateStepInput(id uuid.UUID) procedure.UpdateProcedureStepInput {
	return procedure.UpdateProcedureStepInput{
		ID:          id,
		Index:       1,
		Title:       "Updated Step",
		Description: "Updated Step Description",
		StepType:    "action",
		Version:     1,
	}
}

func TestUpdateProcedureStepDto() []procedure.UpdateProcedureStepDto {
	stepID := uuid.New().String()
	return []procedure.UpdateProcedureStepDto{
		{
			ID:          &stepID,
			Index:       intutil.IntPtr(1),
			Title:       stringutil.StringPtr("Updated Step 1"),
			Description: stringutil.StringPtr("Updated Step 1 Description"),
			StepType:    stringutil.StringPtr("action"),
			Version:     1,
		},
		{
			Index:    intutil.IntPtr(2),
			Title:    stringutil.StringPtr("New Step 2"),
			StepType: stringutil.StringPtr("wait"),
			Version:  1,
		},
	}
}
