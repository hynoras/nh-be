package observation

import (
	"fmt"
	"nh-be/internal/features/observation"
	"nh-be/internal/utils/stringutil"
	"time"

	"github.com/google/uuid"
)

func TestObservationMetadata() observation.ObservationMetadata {
	return observation.ObservationMetadata{
		ID:                    uuid.MustParse("aaaa1111-1111-1111-1111-111111111111"),
		ObservedAt:            time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		Title:                 "Test Observation",
		Notes:                 stringutil.StringPtr("Test Notes"),
		PreviousObservationID: nil,
		CreatedBy:             uuid.MustParse("bbbb2222-2222-2222-2222-222222222222"),
		CreatedAt:             time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
	}
}

func TestObservationMetadataList(count int) []observation.ObservationMetadata {
	list := make([]observation.ObservationMetadata, count)
	for i := 0; i < count; i++ {
		metadata := TestObservationMetadata()
		metadata.ID = uuid.MustParse(fmt.Sprintf("aaaa1111-1111-1111-1111-%012d", i+1))
		list[i] = metadata
	}
	return list
}

func TestObservation() observation.Observation {
	procStepID := uuid.MustParse("cccc3333-3333-3333-3333-333333333333")
	return observation.Observation{
		ID:              uuid.MustParse("aaaa1111-1111-1111-1111-111111111111"),
		ObservedAt:      time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		Title:           "Test Observation",
		Notes:           stringutil.StringPtr("Test Notes"),
		CreatedBy:       uuid.MustParse("bbbb2222-2222-2222-2222-222222222222"),
		CreatedAt:       time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
		ExperimentID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		ProcedureStepID: &procStepID,
	}
}

func TestObservationsResponseDtoList(count int) []observation.ObservationsResponseDto {
	list := make([]observation.ObservationsResponseDto, count)
	for i := 0; i < count; i++ {
		list[i] = observation.ObservationsResponseDto{
			ID:                    fmt.Sprintf("aaaa1111-1111-1111-1111-%012d", i+1),
			ObservedAt:            "2026-02-03T19:15:10Z",
			Title:                 "Test Observation",
			Notes:                 stringutil.StringPtr("Test Notes"),
			PreviousObservationID: nil,
			CreatedBy:             "bbbb2222-2222-2222-2222-222222222222",
			CreatedAt:             "2026-02-03T19:15:10Z",
		}
	}
	return list
}

func TestCreateObservationDto() observation.CreateObservationDto {
	return observation.CreateObservationDto{
		Title:                 "Test Observation",
		Notes:                 stringutil.StringPtr("Test Notes"),
		PreviousObservationID: nil,
		ObservedAt:            time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
	}
}

func TestCreatedObservationResponseDto() observation.CreatedObservationResponseDto {
	procStepID := "cccc3333-3333-3333-3333-333333333333"
	return observation.CreatedObservationResponseDto{
		ID:                    "aaaa1111-1111-1111-1111-111111111111",
		ExperimentID:          "11111111-1111-1111-1111-111111111111",
		ProcedureStepID:       &procStepID,
		ObservedAt:            "2026-02-03T19:15:10Z",
		Title:                 "Test Observation",
		Notes:                 stringutil.StringPtr("Test Notes"),
		PreviousObservationID: nil,
		CreatedBy:             "bbbb2222-2222-2222-2222-222222222222",
		CreatedAt:             "2026-02-03T19:15:10Z",
	}
}
