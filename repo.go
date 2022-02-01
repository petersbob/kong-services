package main

import (
	"database/sql"
	"math/rand"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

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

	//generateTestData(db)

	return postgresRepo{
		db: db,
	}, nil

}

// generateTestData generates some random data.
func generateTestData(db *sql.DB) error {
	insertStatement := `
		INSERT INTO installed_service_versions
			(created_at, updated_at, service_type, version_number)
		VALUES
			(NOW(), NOW(), $1, $2)
	`

	for key, value := range currentServiceTypes {

		for i := 0; i < len(value.VersionsAvailable) && i < 3; i++ {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			if r1.Intn(2) < 1 {
				continue
			}

			_, err := db.Exec(insertStatement, key, value.VersionsAvailable[i])
			if err != nil {
				return err
			}
		}

	}

	return nil
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
