package observation

import (
	"nh-be/internal/utils/stringutil"
	"time"

	"github.com/google/uuid"
)

func MapObservationMetadataToDto(observation ObservationMetadata) ObservationsResponseDto {
	var previousObservationID *string
	if observation.PreviousObservationID != nil {
		previousObservationID = stringutil.StringPtr(observation.PreviousObservationID.String())
	}

	return ObservationsResponseDto{
		ID:                    observation.ID.String(),
		Title:                 observation.Title,
		Notes:                 observation.Notes,
		PreviousObservationID: previousObservationID,
		CreatedBy:             observation.CreatedBy.String(),
		ObservedAt:            observation.ObservedAt.Format(time.RFC3339),
		CreatedAt:             observation.CreatedAt.Format(time.RFC3339),
	}
}

func MapObservationsMetadataToDto(observation []ObservationMetadata) []ObservationsResponseDto {
	obs := make([]ObservationsResponseDto, 0, len(observation))
	for _, observation := range observation {
		obs = append(obs, MapObservationMetadataToDto(observation))
	}
	return obs
}

func MapCreateInputToObservation(input CreateObservationInput, createdBy, expId, procStepId uuid.UUID) Observation {
	obs := Observation{
		CreatedBy:             createdBy,
		Title:                 input.Title,
		Notes:                 input.Notes,
		ObservedAt:            input.ObservedAt,
		PreviousObservationID: input.PreviousObservationID,
		ProcedureStepID:       &procStepId,
		ExperimentID:          expId,
	}

	return obs
}

func MapObsToCreatedObsResponseDto(observation Observation) CreatedObservationResponseDto {
	res := CreatedObservationResponseDto{
		ID:           observation.ID.String(),
		ExperimentID: observation.ExperimentID.String(),
		ObservedAt:   observation.ObservedAt.Format(time.RFC3339),
		Title:        observation.Title,
		Notes:        observation.Notes,
		CreatedBy:    observation.CreatedBy.String(),
		CreatedAt:    observation.CreatedAt.Format(time.RFC3339),
	}

	if observation.ProcedureStepID != nil {
		res.ProcedureStepID = stringutil.StringPtr(observation.ProcedureStepID.String())
	}

	if observation.PreviousObservationID != nil {
		res.PreviousObservationID = stringutil.StringPtr(observation.PreviousObservationID.String())
	}

	return res
}

func MapCreateDtoToInput(dto CreateObservationDto) (CreateObservationInput, error) {
	input := CreateObservationInput{
		Title:      dto.Title,
		Notes:      dto.Notes,
		ObservedAt: dto.ObservedAt,
	}

	if dto.PreviousObservationID != nil {
		id, err := uuid.Parse(*dto.PreviousObservationID)
		if err != nil {
			return CreateObservationInput{}, err
		}
		input.PreviousObservationID = &id
	}

	return input, nil
}
