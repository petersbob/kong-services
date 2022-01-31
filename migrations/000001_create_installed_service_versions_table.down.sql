BEGIN;

DROP TABLE IF EXISTS installed_service_versions;

DROP INDEX IF EXISTS  idx_installed_service_version_service_type;

COMMIT;