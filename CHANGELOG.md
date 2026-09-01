## 0.1.1 - 2026-09-01

- Fix CSV and TSV country sources to return their real ISO 3166-1 alpha-3
  codes instead of incorrectly mirroring the alpha-2 code. The bundled
  production CSV now carries the verified ISO-3 value for every country.
- Preserve the legacy two-column source format as a compatibility fallback.

## 0.1.0 - 2026-08-29

First fleet-registered release. The service itself already existed; this
version onboards it onto the `mesh-0exec` container fleet.

### Added
- `service.yaml` / `deploy.yaml` — fleet registry + deploy descriptors
  (`id: country-iso-matcher`, `category: geo`, host port `18315`).
- `GET /version` → `{"version":"<ver>"}`, the fleet-canonical version probe
  that `fleet-runner deploy`'s smoke gate compares against the pushed tag.
- `src/internal/version` — single source of truth for service id + version,
  now also feeding the `build_info` Prometheus metric.
- `CHANGELOG.md`, plus `.gitignore` rules for pre-built host binaries.

### Changed
- `GET /health` now returns the fleet-canonical
  `{"status":"ok","service":"country-iso-matcher","version":"<ver>"}`
  (was `{"status":"healthy", ...}` with no version).
- `Dockerfile` — non-root `app` user, `tini` as PID 1, `HEALTHCHECK`,
  `-trimpath -ldflags="-s -w"`, Go 1.25, default port `18315`.
- `docker-compose.yml` — canonical fleet shape: pulls
  `ghcr.io/baditaflorin/country-iso-matcher` instead of building locally,
  and drops the `./data` / `./web` bind mounts (they don't exist on the
  dockerhost, so the container would have started with no country data).
- `GUI_ENABLED` now defaults to **false**. The `/admin` GUI also mounts
  `/api/config` (returns the full config including `database.password`)
  and `/api/config/save` (writes caller-supplied config to disk), neither
  of which has its own authorization check.
- `metrics.SetBuildInfo` reports the real version instead of a hardcoded
  `"1.0.0"`.

### Fixed
- `src/internal/service/country_service_test.go` no longer compiled — it
  referenced the pre-ISO3 `domain.Country{Code,Name}` / `CountryResponse.ISOCode`
  shape. Updated to `ISO2`/`ISO3`/`Names` and `ISO2Code`/`ISO3Code`.
- `benchmarks/country_lookup_test.go` no longer compiled — it imported
  `github.com/baditaflorin/country-iso-matcher/internal/...`, a module path
  this repo never used. The lookup benchmark moved to
  `src/internal/service/country_lookup_bench_test.go` (Go forbids importing
  `src/internal/...` from outside `src/`); the normalizer benchmark stayed
  in `benchmarks/`.

### Removed
- Committed pre-built `bin/server` host binary (13.9 MB).

### Known issues
- `iso3Code` is **not real** under the default `DATA_SOURCE=csv`:
  `data/countries.csv` has only `code,name`, so `csv_loader.go` falls back
  to `ISO3: code` and the API returns e.g. `{"iso2Code":"DE","iso3Code":"DE"}`.
  The `json` data source (`data/countries/*.json`) does carry correct
  alpha-3 codes and multilingual names, but only covers 11 countries, so
  CSV remains the deployed source. Fixing this needs an alpha-3 column
  backfilled into the CSV corpus.
