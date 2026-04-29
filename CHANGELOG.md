# Changelog

## [1.1.0] - 2026-04-29

### Added
- Landing page with hacker terminal theme, live SSE price feed, and interactive API explorer
- SSE hub pattern for live streaming with single upstream fetch and fan-out broadcast
- OpenAPI 3.1 spec served at /openapi.json
- Dynamic version display via /healthz API
- Security scanning in CI pipeline with govulncheck and gosec

### Changed
- READMEs updated with live demo URL, SSE endpoints, and OpenAPI documentation
- GitHub Actions pinned to immutable SHA hashes
- Docker base images pinned to SHA256 digests
- Docker container now runs as non-root user
- Docker Compose with memory (256M) and CPU (1.0) resource limits
- Expanded .dockerignore with sensitive file patterns

### Fixed
- HTTP server timeouts to prevent Slowloris attacks
- Ticker input validation and URL encoding for safe upstream requests
- Parameter bounds on count/offset query parameters
- Null Google response handling in decoder
- XSS via innerHTML in SSE live feed grid
- Verbose error messages replaced with generic client-facing responses
- Upstream response body capped at 10 MB via io.LimitReader
- SSE connection limit (max 50) to prevent resource exhaustion
- Security headers: X-Frame-Options, CSP, Referrer-Policy, X-Content-Type-Options
- ASCII logo rendering with pre tag to prevent formatter corruption

## [1.0.0] - 2026-04-26

### Added
- Google Finance REST API with 13 endpoints covering quotes, company info, charts, news, financials, and market data
- Support for stocks, crypto, forex, ETFs, and indices across all exchanges including Borsa Istanbul
- Docker support with multi-stage Dockerfile and docker-compose.yml

### Changed
- Docker usage section added to both Turkish and English READMEs
- Version tracking via VERSION file

### Fixed
- CI pipeline with GitHub Actions for build validation and Docker healthcheck
