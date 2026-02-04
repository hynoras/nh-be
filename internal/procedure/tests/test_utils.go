package procedure

import (
	"context"
	"nh-be/constant"
	"nh-be/internal/experiment/root"
	"nh-be/internal/procedure"
	"time"

	"github.com/google/uuid"
)

func ContextWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), constant.CtxUserId, userID)
}

func TestProcedureList() []procedure.Procedure {
	return []procedure.Procedure{
		{
			ID:          uuid.New(),
			Title:       "Test Procedure",
			Description: "Test Procedure Description",
			Version:     1,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Test Procedure 2",
			Description: "Test Procedure Description 2",
			Version:     2,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
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
		UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
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
			UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("33333333-1234-1234-1234-555533332222"),
			ProcedureID: uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			Index:       2,
			Title:       "Test Step 2",
			Description: "Test Step Description 2",
			CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		},
	}
}

func TestProcedureDetailWithRelations() procedure.Procedure {
	procID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	exp1ID := uuid.MustParse("aaaaaaaa-1111-1111-1111-111111111111")
	exp2ID := uuid.MustParse("bbbbbbbb-2222-2222-2222-222222222222")

	return procedure.Procedure{
		ID:          procID,
		Title:       "Test Procedure",
		Description: "Test Procedure Description",
		Version:     1,
		ParentID:    nil,
		CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		Steps:       TestStepList(),
		Experiments: []procedure.ProcedureExperimentAssignment{
			{
				ID:           uuid.MustParse("cccccccc-3333-3333-3333-333333333333"),
				ProcedureID:  procID,
				ExperimentID: exp1ID,
				Experiment: root.Experiment{
					ID:        exp1ID,
					Title:     "Test Experiment 1",
					Objective: "Test Objective 1",
				},
				CreatedAt: time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
				UpdatedAt: time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			},
			{
				ID:           uuid.MustParse("dddddddd-4444-4444-4444-444444444444"),
				ProcedureID:  procID,
				ExperimentID: exp2ID,
				Experiment: root.Experiment{
					ID:        exp2ID,
					Title:     "Test Experiment 2",
					Objective: "Test Objective 2",
				},
				CreatedAt: time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
				UpdatedAt: time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			},
		},
	}
}
