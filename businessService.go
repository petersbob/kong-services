package main

import (
	"strconv"
	"strings"
)

type businessService struct {
	repo repository
}

type servicesFilter struct {
	search string
}

func NewBusinessService(repo repository) businessService {
	return businessService{
		repo: repo,
	}
}

func getServiceType(typeCode ServiceTypeCode) ServiceType {
	return currentServiceTypes[typeCode]
}

func filterServiceTypes(filter servicesFilter) []ServiceType {
	serviceTypes := []ServiceType{}

	// filter based on search field
	for _, value := range currentServiceTypes {
		typeCodeString := strconv.FormatUint(uint64(value.TypeCode), 10)

		if strings.Contains(value.Name, filter.search) ||
			strings.Contains(value.Description, filter.search) ||
			strings.Contains(typeCodeString, filter.search) {

			serviceTypes = append(serviceTypes, value)
		}
	}

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
