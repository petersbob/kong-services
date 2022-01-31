BEGIN;

CREATE TABLE IF NOT EXISTS installed_service_versions(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    service_type INTEGER NOT NULL,
    version_number INTEGER NOT NULL
);

CREATE INDEX idx_installed_service_version_service_type on installed_service_versions(service_type);

COMMIT;