package main

import "time"

type repository interface {
	GetVersionsInUseByServiceType(typeCode ServiceTypeCode) ([]ServiceVersion, error)
}

// TestingRepo is an implementation of the repository interface
// but meant to only be used during testing. This should be replaced by a real
// persistence layer later.
type testingRepo struct {
}

func NewTestingRepo() testingRepo {
	return testingRepo{}
}

var currentVersionsByService = map[ServiceTypeCode][]ServiceVersion{
	ServiceTypeCodeDatabase: {
		{
			ServiceType:   ServiceTypeCodeDatabase,
			VersionNumber: 1,
			CreatedAt:     time.Now(),
		},
		{
			ServiceType:   ServiceTypeCodeDatabase,
			VersionNumber: 3,
			CreatedAt:     time.Now(),
		},
	},
	ServiceTypeCodeReporting: nil,
	ServiceTypeCodeCurrencyConversion: {
		{
			ServiceType:   ServiceTypeCodeCurrencyConversion,
			VersionNumber: 5,
			CreatedAt:     time.Now(),
		},
	},
	ServiceTypeCodeTranslation: nil,
	ServiceTypeCodeNotifications: {
		{
			ServiceType:   ServiceTypeCodeNotifications,
			VersionNumber: 1,
			CreatedAt:     time.Now(),
		},
	},
}

func (r testingRepo) GetVersionsInUseByServiceType(typeCode ServiceTypeCode) ([]ServiceVersion, error) {
	return currentVersionsByService[typeCode], nil
}
