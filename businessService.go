package main

import (
	"sort"
	"strconv"
	"strings"
)

type businessService struct {
	repo repository
}

type servicesFilter struct {
	search   string
	sort     string
	page     int
	pageSize int
}

func NewBusinessService(repo repository) businessService {
	return businessService{
		repo: repo,
	}
}

func getServiceType(typeCode ServiceTypeCode) ServiceType {
	return currentServiceTypes[typeCode]
}

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

func getServiceTypesPage(page int, pageSize int, serviceTypes []ServiceType) []ServiceType {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 1
	}

	serviceTypesLength := len(serviceTypes)

	startingIndex := (pageSize * page) - 1

	serviceTypesPage := []ServiceType{}

	for i := startingIndex; i < startingIndex+pageSize && i < serviceTypesLength; i++ {
		serviceTypesPage = append(serviceTypesPage, serviceTypes[i])
	}

	return serviceTypesPage
}

func filterServiceTypes(filter servicesFilter) []ServiceType {
	serviceTypes := []ServiceType{}

	for _, value := range currentServiceTypes {
		serviceTypes = append(serviceTypes, value)
	}

	serviceTypes = filterServiceTypesBySearchValue(filter.search, serviceTypes)

	serviceTypes = sortServiceTypes(filter.sort, serviceTypes)

	serviceTypes = getServiceTypesPage(filter.page, filter.pageSize, serviceTypes)

	return serviceTypes
}

// GetServices find alls the current services that are in use by the user
func (s businessService) GetServices(filter servicesFilter) ([]Service, error) {
	services := []Service{}

	serviceTypes := filterServiceTypes(filter)

	for i := range serviceTypes {

		versionsInUse, err := s.repo.GetVersionsInUseByServiceType(serviceTypes[i].TypeCode)
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
