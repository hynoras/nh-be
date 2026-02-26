package observation

type MeasurementDataType string

const (
	MeasurementDataTypeNumeric MeasurementDataType = "numeric"
	MeasurementDataTypeText    MeasurementDataType = "text"
	MeasurementDataTypeBoolean MeasurementDataType = "boolean"
)
