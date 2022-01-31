package main

import (
	"sort"
	"strconv"
	"strings"
)

func NewBusinessService(repo repository) businessService {
	return businessService{
		repo: repo,
	}
}

func getServiceType(typeCode ServiceTypeCode) ServiceType {
	return currentServiceTypes[typeCode]
}

// filterServiceTypesBySearchValue filters the service types by trying to find the substring "searchString"
// in their description, name, or type code.
func filterServiceTypesBySearchValue(searchString string, servicesTypes []ServiceType) []ServiceType {
	filteredServiceTypes := []ServiceType{}

	for i := range servicesTypes {
		typeCodeString := strconv.FormatUint(uint64(servicesTypes[i].TypeCode), 10)

		if strings.Contains(servicesTypes[i].Description, searchString) ||
			strings.Contains(servicesTypes[i].Name, searchString) ||
			strings.Contains(typeCodeString, searchString) {

			filteredServiceTypes = append(filteredServiceTypes, servicesTypes[i])
		}

	}

	return filteredServiceTypes
}

// sortServiceTypes sorts a list of service types by name, description, or type code.
// If no sort string is provided or is unrecognized, the service types are sorted by type code.
func sortServiceTypes(sortString string, serviceTypes []ServiceType) []ServiceType {
	switch sortString {
	case "name":
		sort.Slice(serviceTypes, func(i, j int) bool { return serviceTypes[i].Name < serviceTypes[j].Name })
	case "description":
		sort.Slice(serviceTypes, func(i, j int) bool { return serviceTypes[i].Description < serviceTypes[j].Description })
	case "typeCode":
		fallthrough
	default:
		sort.Slice(serviceTypes, func(i, j int) bool { return serviceTypes[i].TypeCode < serviceTypes[j].TypeCode })
	}

	return serviceTypes
}

// getServicesPage finds a subset of services in a list of services given a page and the number of
// services on each page. If the page is less than 1, it defaults to 1. If page size is less than
// 1, all the services are returned.
func getServicesPage(page int, pageSize int, services []Service) []Service {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		return services
	}

	servicesLength := len(services)

	startingIndex := (pageSize * page) - 1

	servicesPage := []Service{}

	for i := startingIndex; i < startingIndex+pageSize && i < servicesLength; i++ {
		servicesPage = append(servicesPage, services[i])
	}

	return servicesPage
}

// getInUseServices receives a list of service types and returns a list of services that are both
// of types included in the service types list and also in use by the user. "in use" is defined
// as the user currently having a version of the service type installed.
func (s businessService) getInUseServices(serviceTypes []ServiceType) ([]Service, error) {
	inUseServices := []Service{}

	for i := range serviceTypes {

		versionsInUse, err := s.repo.GetVersionsInstalledByServiceType(serviceTypes[i].TypeCode)
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
			ServiceType:   serviceTypes[i],
			VersionsInUse: versionNumbersInUse,
		}

		inUseServices = append(inUseServices, service)
	}

	return inUseServices, nil
}

// GetServices find alls the current services that are in use by the user
func (s businessService) GetServices(filter ServicesFilter) ([]Service, error) {
	// create any array of the current service types
	serviceTypes := []ServiceType{}
	for _, value := range currentServiceTypes {
		serviceTypes = append(serviceTypes, value)
	}

	// filter the service types by search query
	searchFilteredServiceTypes := filterServiceTypesBySearchValue(filter.Search, serviceTypes)

	// sort the service types
	sortedServiceTypes := sortServiceTypes(filter.Sort, searchFilteredServiceTypes)

	// get a list of only the in use services (all service types where the user has a version installed)
	inUseServices, err := s.getInUseServices(sortedServiceTypes)
	if err != nil {
		return nil, err
	}

	// get the subpage of the results
	servicesPage := getServicesPage(filter.Page, filter.PageSize, inUseServices)

	return servicesPage, nil
}

// GetService find the details on a specific in use service
func (s businessService) GetService(typeCode ServiceTypeCode) (Service, error) {
	serviceType := getServiceType(typeCode)

	versionsInUse, err := s.repo.GetVersionsInstalledByServiceType(typeCode)
	if err != nil {
		return Service{}, err
	}

	if len(versionsInUse) == 0 {
		return Service{}, ErrorNotFound
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
func (s businessService) GetServiceVersions(typeCode ServiceTypeCode) ([]InstalledServiceVersion, error) {
	versionsInUse, err := s.repo.GetVersionsInstalledByServiceType(typeCode)
	if err != nil {
		return nil, err
	}

	if len(versionsInUse) == 0 {
		return nil, ErrorNotFound
	}

	return versionsInUse, nil
}

// GetServiceVersion find info on a specific verion of a service in use
func (s businessService) GetServiceVersion(typeCode ServiceTypeCode, versionNumber uint) (InstalledServiceVersion, error) {
	versionsInUse, err := s.repo.GetVersionsInstalledByServiceType(typeCode)
	if err != nil {
		return InstalledServiceVersion{}, err
	}

	for i := range versionsInUse {
		if versionsInUse[i].VersionNumber == versionNumber {
			return versionsInUse[i], nil
		}
	}

	return InstalledServiceVersion{}, ErrorNotFound
}
