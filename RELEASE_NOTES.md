# go-log v1.3.0

## Highlights

- New `EnrichHTTPMeta` helper to centralize HTTP request/response diagnostic metadata (origin, method/path/query, selected headers, timestamp, hostname, and short stack for 5xx errors).
- Size-based log rotation support (configurable). Rotation remains day-based by default; size rotation is available for tests and can be enabled as needed.
- Several correctness and defensive fixes to prevent runtime panics and race issues under load.
- Large test-suite additions: concurrency stress tests (race-checked), defensive unit tests, integration tests that exercise disk log creation and rotation, and a size-rotation test.
- Test helpers moved behind an explicit build tag (`test`) so they don't appear in normal builds.

## What’s new (detailed)

- Feature: `EnrichHTTPMeta(status int, req *http.Request, meta map[string]interface{}, callerSkip int) map[string]interface{}`
  - Centralizes and standardizes enrichment of metadata for error logging and JSON error responses. Prefer calling this instead of duplicating enrichment logic.
- Feature: Size-based rotation
  - `gMaxFileSizeBytes` controls rotation by file size when > 0. Tests use a test-only helper to set a small threshold and validate rotation behavior.
- Fix: nil-guard when closing log files — prevents panics when closing an absent file.
- Fix: defensive `Gorm.Print` – guards against variable argument shapes from GORM to avoid panics.
- Fix: `Logger.Println` now reliably logs printed values (no empty format).
- Tests: Added
  - Concurrency stress tests (run with `-race`).
  - Unit tests for `EnrichHTTPMeta` and defensive cases.
  - Integration test that writes logs to disk and validates symlink and rotated files.
  - Size rotation test that forces rotation via test helper.

## Test-only helpers (only available when compiled with build tag `test`)

- `SetMaxFileSizeBytes(n int64)` — test-only setter to override size threshold.
- `ResetForTests()` — closes open files and resets internal logger state between tests.
- These helpers are in `test_helpers_test.go` and are compiled only with `-tags=test`.

## Upgrade & migration notes

- If you use Go modules, update your module requirement and remove any local `replace` to the local copy:

  ```diff
  - require github.com/dainiauskas/go-log v1.2.9
  + require github.com/dainiauskas/go-log v1.3.0
  - replace github.com/dainiauskas/go-log => /path/to/local/go-log
  ```

- Running tests locally or in CI:
  - To run tests that rely on test helpers:
    ```bash
    go test -tags=test ./... -race -v
    ```
  - Normal test run (without test helpers):
    ```bash
    go test ./... -v
    ```

- Note: After pushing the tag the public checksum service (sum.golang.org) can take a short while to reflect the new version. If CI needs to verify immediately, you can temporarily fetch direct:
  ```bash
  GOPROXY=direct go mod tidy
  ```
  (Avoid setting `GOSUMDB=off` in CI; it bypasses checksum verification and is not recommended for long-term use.)

## Changelog (concise)

- Added: `EnrichHTTPMeta` for standardized metadata enrichment.
- Added: size-based rotation support and tests.
- Added: concurrency, integration, and defensive tests.
- Fixed: nil-close guard when closing log files.
- Fixed: defensive checks in `Gorm.Print` for varying arg shapes.
- Fixed: `Logger.Println` formatting issue.

## Suggested PR/Release description (short)

This release centralizes HTTP metadata enrichment (via `EnrichHTTPMeta`), hardens logging against panics, adds size-based rotation and a broad test suite, and hides test helpers behind a build tag. Consumers can upgrade to v1.3.0 and remove local replaces; see notes for running tests that require test helpers.
