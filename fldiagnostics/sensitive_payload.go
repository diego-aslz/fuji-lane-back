package fldiagnostics

// SensitivePayload contains sensitive information that should be hidden from logging mechanisms
type SensitivePayload interface {
	FilterSensitiveInformation() SensitivePayload
}
