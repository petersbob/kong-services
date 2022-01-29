package main

import "time"

// TestingRepo is an implementation of the repository interface
// but meant to only be used during testing. This should be replaced by a real
// persistence layer later.
type testingRepo struct {
}

type ServiceVersion struct {
	Type      ServiceTypeCode
	Version   uint
	CreatedAt time.Time
}

var currentVersionsByService = map[ServiceTypeCode][]ServiceVersion{
	ServiceTypeCodeDatabase:           nil,
	ServiceTypeCodeReporting:          nil,
	ServiceTypeCodeCurrencyConversion: nil,
	ServiceTypeCodeTranslation:        nil,
	ServiceTypeCodeNotifications:      nil,
}

func (r testingRepo) GetVersionsByServiceType(typeCode ServiceTypeCode) ([]ServiceVersion, error) {
	return currentVersionsByService[typeCode], nil
}
