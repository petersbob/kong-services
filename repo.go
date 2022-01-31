package main

import (
	"database/sql"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type repository interface {
	GetVersionsInstalledByServiceType(typeCode ServiceTypeCode) ([]InstalledServiceVersion, error)
}

type postgresRepo struct {
	db *sql.DB
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
	if err != nil && err != migrate.ErrNoChange {
		return postgresRepo{}, err
	}

	return postgresRepo{
		db: db,
	}, nil

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
