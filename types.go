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
	GetVersionsInUseByServiceType(typeCode ServiceTypeCode) ([]ServiceVersion, error)
}

type BusinessService struct {
	repo repository
}

// ServiceType
type ServiceType struct {
	TypeCode          ServiceTypeCode `json:"type_code"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	VersionsAvailable []uint          `json:"versions_available"`
}

type Service struct {
	ServiceType
	VersionsInUse []uint `json:"versions_in_use"`
}
