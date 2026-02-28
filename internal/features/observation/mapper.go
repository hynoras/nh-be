package observation

import (
	"nh-be/internal/utils/stringutil"
	"time"
)

func MapObservationToDto(observation ObservationMetadata) ObservationsResponseDto {
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

func MapObservationsToDto(observation []ObservationMetadata) []ObservationsResponseDto {
	obs := make([]ObservationsResponseDto, 0, len(observation))
	for _, observation := range observation {
		obs = append(obs, MapObservationToDto(observation))
	}
	return obs
}
