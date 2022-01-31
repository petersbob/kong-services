package main

import (
	"database/sql"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type repository interface {
	GetVersionsInstalledByServiceType(typeCode ServiceTypeCode) ([]InstalledServiceVersion, error)
}

// TestingRepo is an implementation of the repository interface
// but meant to only be used during testing. This should be replaced by a real
// persistence layer later.
type testingRepo struct {
}

type postgresRepo struct {
	db *sql.DB
}

func NewTestingRepo() testingRepo {
	return testingRepo{}
}

func NewPostgresRepo(databaseURL, migrationsPath string) (postgresRepo, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return postgresRepo{}, err
	}

	err = db.Ping()
	if err != nil {
		return postgresRepo{}, err
	}

	m, err := migrate.New(
		migrationsPath,
		databaseURL,
	)
	if err != nil {
		return postgresRepo{}, err
	}
	err = m.Up()
	if err != nil {
		return postgresRepo{}, err
	}

	return postgresRepo{
		db: db,
	}, nil

}

var currentVersionsByService = map[ServiceTypeCode][]InstalledServiceVersion{
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

func (r testingRepo) GetVersionsInstalledByServiceType(typeCode ServiceTypeCode) ([]InstalledServiceVersion, error) {
	return currentVersionsByService[typeCode], nil
}

func (r postgresRepo) GetVersionsInstalledByServiceType(typeCode ServiceTypeCode) ([]InstalledServiceVersion, error) {
	installedServiceVersions := []InstalledServiceVersion{}

	selectInstalledServiceVersions := `
		SELECT
			created_at,
			service_type,
			version_number
		FROM
			installed_service_versions
		WHERE
			service_type = $1
	`

	rows, err := r.db.Query(selectInstalledServiceVersions, typeCode)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		installedServiceVersion := InstalledServiceVersion{}
		err := rows.Scan(
			&installedServiceVersion.CreatedAt,
			&installedServiceVersion.ServiceType,
			&installedServiceVersion.VersionNumber,
		)

		if err != nil {
			return nil, err
		}

		installedServiceVersions = append(installedServiceVersions, installedServiceVersion)
	}

	return installedServiceVersions, nil
}
