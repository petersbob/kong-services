package main

type businessService struct {
	repo repository
}

func NewBusinessService(repo repository) businessService {
	return businessService{
		repo: repo,
	}
}

// CurrentServiceTypes is a map of all of the current service types that are available to create.
var currentServiceTypes = map[ServiceTypeCode]ServiceType{
	ServiceTypeCodeDatabase: {
		TypeCode:          ServiceTypeCodeDatabase,
		Name:              "Database Service",
		Description:       "A service for running databases",
		VersionsAvailable: []uint{1, 2, 3},
	},
	ServiceTypeCodeReporting: {
		TypeCode:          ServiceTypeCodeReporting,
		Name:              "Reporting Service",
		Description:       "A service for running reports",
		VersionsAvailable: []uint{55, 7},
	},
	ServiceTypeCodeCurrencyConversion: {
		TypeCode:          ServiceTypeCodeCurrencyConversion,
		Name:              "Currency Conversion",
		Description:       "A service for doing currency convertions",
		VersionsAvailable: []uint{1, 2, 3, 4, 5, 6, 7, 8},
	},
	ServiceTypeCodeTranslation: {
		TypeCode:          ServiceTypeCodeTranslation,
		Name:              "Translation service",
		Description:       "A service for doing language translations",
		VersionsAvailable: []uint{12, 14},
	},
	ServiceTypeCodeNotifications: {
		TypeCode:          ServiceTypeCodeNotifications,
		Name:              "Notifications service",
		Description:       "A service sending notifications",
		VersionsAvailable: []uint{1},
	},
}

func getServiceType(typeCode ServiceTypeCode) ServiceType {
	return currentServiceTypes[typeCode]
}

// GetServices find alls the current services that are in use by the user
func (s businessService) GetServices() ([]Service, error) {
	services := []Service{}

	for key, value := range currentServiceTypes {
		versionsInUse, err := s.repo.GetVersionsInUseByServiceType(key)
		if err != nil {
			return nil, err
		}

		if len(versionsInUse) == 0 {
			continue
		}

		versionNumbersInUse := []uint{}

		for i := range versionsInUse {
			versionNumbersInUse = append(versionNumbersInUse, versionsInUse[i].VersionNumber)
		}

		service := Service{
			ServiceType:   value,
			VersionsInUse: versionNumbersInUse,
		}

		services = append(services, service)
	}

	return services, nil
}

// GetService find the details on a specific in use service
func (s businessService) GetService(typeCode ServiceTypeCode) (Service, error) {
	serviceType := getServiceType(typeCode)

	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return Service{}, err
	}

	if len(versionsInUse) == 0 {
		return Service{}, errorNotFound
	}

	versionNumbersInUse := []uint{}

	for i := range versionsInUse {
		versionNumbersInUse = append(versionNumbersInUse, versionsInUse[i].VersionNumber)
	}

	service := Service{
		ServiceType:   serviceType,
		VersionsInUse: versionNumbersInUse,
	}

	return service, nil
}

// GetServiceVersions finds all of the versions of a service in use
func (s businessService) GetServiceVersions(typeCode ServiceTypeCode) ([]ServiceVersion, error) {
	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return nil, err
	}

	if len(versionsInUse) == 0 {
		return nil, errorNotFound
	}

	return versionsInUse, nil
}

// GetServiceVersion find info on a specific verion of a service in use
func (s businessService) GetServiceVersion(typeCode ServiceTypeCode, versionNumber uint) (ServiceVersion, error) {
	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return ServiceVersion{}, err
	}

	for i := range versionsInUse {
		if versionsInUse[i].VersionNumber == versionNumber {
			return versionsInUse[i], nil
		}
	}

	return ServiceVersion{}, errorNotFound
}
