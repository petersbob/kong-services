package main

func NewBusinessService(repo repository) BusinessService {
	return BusinessService{
		repo: repo,
	}
}

var currentServices = map[ServiceTypeCode]ServiceType{
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
	return currentServices[typeCode]
}

func (s BusinessService) GetServices() ([]Service, error) {

	services := []Service{}

	for key, value := range currentServices {

		versionsInUse, err := s.repo.GetVersionsInUseByServiceType(key)
		if err != nil {
			return nil, err
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

func (s BusinessService) GetService(typeCode ServiceTypeCode) (Service, error) {

	serviceType := getServiceType(typeCode)

	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return Service{}, err
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

func (s BusinessService) GetServiceVersions(typeCode ServiceTypeCode) ([]ServiceVersion, error) {

	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return nil, err
	}

	return versionsInUse, nil
}

func (s BusinessService) GetServiceVersion(typeCode ServiceTypeCode, versionNumber uint) (ServiceVersion, error) {
	versionsInUse, err := s.repo.GetVersionsInUseByServiceType(typeCode)
	if err != nil {
		return ServiceVersion{}, err
	}

	for i := range versionsInUse {
		if versionsInUse[i].VersionNumber == versionNumber {
			return versionsInUse[i], nil
		}
	}

	return ServiceVersion{}, nil
}
