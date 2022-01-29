package main

type ServiceTypeCode uint

const (
	ServiceTypeCodeDatabase           ServiceTypeCode = 1
	ServiceTypeCodeReporting          ServiceTypeCode = 2
	ServiceTypeCodeCurrencyConversion ServiceTypeCode = 3
	ServiceTypeCodeTranslation        ServiceTypeCode = 4
	ServiceTypeCodeNotifications      ServiceTypeCode = 5
)

type repository interface {
	GetVersionsByServiceType(typeCode ServiceTypeCode) ([]ServiceVersion, error)
}

type service struct {
	repo repository
}

// ServiceType
type ServiceType struct {
	TypeCode          ServiceTypeCode
	Name              string
	Description       string
	VersionsAvailable []uint
}

type Service struct {
	Type     ServiceType
	Versions []ServiceVersion
}
