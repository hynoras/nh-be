package experiment

import (
	"time"

	"nh-be/internal/features/experiment"
	"nh-be/internal/utils/timeutil"

	"github.com/google/uuid"
)

func TestExperiment() experiment.Experiment {
	procedureID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	return experiment.Experiment{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Identifier:  "EXP-0001",
		Title:       "Test Experiment",
		Objective:   "Test Objective",
		Status:      experiment.ExperimentDraft,
		Type:        experiment.ExperimentExploratoryType,
		Version:     1,
		CreatedByID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		ProcedureID: &procedureID,
		CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		UpdatedAt:   timeutil.TimePtr(time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC)),
	}
}

func TestExperimentsQueryDto() []experiment.ExperimentsQueryDto {
	return []experiment.ExperimentsQueryDto{
		{
			Identifier:    "EXP-0001",
			Title:         "Test Experiment",
			Objective:     "Test Objective",
			Status:        experiment.ExperimentDraft,
			Type:          experiment.ExperimentExploratoryType,
			Creator:       "Test Creator",
			Updater:       "Test Updater",
			CreatedAt:     time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			UpdatedAt:     time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			ProcedureName: "Test Procedure",
		},
	}
}

func TestExperimentsResponseDto() []experiment.ExperimentsResponseDto {
	return []experiment.ExperimentsResponseDto{
		{
			Identifier:    "EXP-0001",
			Title:         "Test Experiment",
			Objective:     "Test Objective",
			Status:        string(experiment.ExperimentDraft),
			Type:          string(experiment.ExperimentExploratoryType),
			Creator:       "Test Creator",
			Updater:       "Test Updater",
			CreatedAt:     time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			UpdatedAt:     time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
			ProcedureName: "Test Procedure",
		},
	}
}

func TestExperimentDetailResponseDto() experiment.ExperimentResponseDto {
	exp := TestExperiment()
	return experiment.MapExperimentToDto(exp)
}
