package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type LicenseRootResolver interface {
	ValidateLicense(context.Context, ValidateLicenseArgs) (LicenseValidation, error)

	InternalListForOrganization(ctx context.Context, id string) ([]LicenseResolver, error)
}

type LicenseResolver interface {
	ID() graphql.ID
	Seats() int32
	UsedSeats() int32
	ExpiresAt() int32
	LicenseKey() string
}

type ValidateLicenseArgs struct {
	Input ValidateLicenseInput
}

type ValidateLicenseInput struct {
	Key           string
	Version       string
	BootedAt      int32
	UserCount     int32
	CodebaseCount int32
}

type LicenseValidation interface {
	ID() graphql.ID
	Status() LicenseValidationStatus
	Message() *string
}

type LicenseValidationStatus string

const (
	LicenseValidationStatusUnknown LicenseValidationStatus = "Unknown"
	LicenseValidationStatusOk      LicenseValidationStatus = "Ok"
	LicenseValidationStatusInvalid LicenseValidationStatus = "Invalid"
	LicenseValidationStatusExpired LicenseValidationStatus = "Expired"
)