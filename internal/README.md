# Bitcoin Investment Tracker - Internal Package

> This directory contains the core internal packages of the Bitcoin Investment Tracker application.

## Directory Structure

```
internal
├── README.md
├── calculator     # Investment calculations and profit/loss logic
├── formatter      # Currency and number formatting utilities
├── metrics        # Prometheus metrics collection
├── webserver     # HTTP server and web interface
└── websocket     # Binance WebSocket client
```

The internal packages follow clean architecture principles, with clear separation of concerns:

- **Business Logic**: `calculator/` handles core investment calculations
- **Presentation**: `formatter/` manages data presentation and formatting
- **Infrastructure**: 
  - `metrics/` for observability
  - `webserver/` for HTTP interface
  - `websocket/` for external data sources

This structure ensures:
- Independence between business rules and external dependencies
- Easy testing and maintenance
- Clear boundaries between different parts of the system
- Flexibility to change external implementations

## Package Structure

```
internal/
├── calculator/    # Investment calculations and profit/loss logic
├── formatter/     # Currency and number formatting utilities (BRL, USD)
├── metrics/       # Prometheus metrics collection and exporters
├── webserver/     # HTTP server, templates and API endpoints
└── websocket/     # Binance WebSocket client for real-time price updates
```

## Local Development

### Prerequisites
- Go 1.22 or higher
- Make (optional, for Makefile commands)
- Access to a PostgreSQL database

### Setup

1. Initialize Go modules (if not already done):
   ```bash
   go mod init bitcoin-investment-tracker
   ```

2. Install dependencies:
   ```bash
   go mod download
   go mod tidy
   ```

3. Verify dependencies:
   ```bash
   go mod verify
   ```

### Running the Application

1. Run directly with Go:
   ```bash
   go run cmd/main.go
   ```

2. Build and run binary:
   ```bash
   # Build
   go build -o bin/bitcoin-tracker cmd/main.go

   # Run
   ./bin/bitcoin-tracker
   ```

### Development Commands

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Update dependencies
go get -u ./...
```

### Environment Variables

Make sure to set these environment variables or create a `.env` file in the root directory:

```bash
# API URLs
MERCADO_BITCOIN_API_URL=https://www.mercadobitcoin.net/api/BTC/ticker/
BINANCE_USDT_API_URL=https://api.binance.com/api/v3/ticker/price?symbol=USDTBRL

# Investment Details
INVESTMENT_BRL=20000      # Your initial investment in BRL
BTC_AMOUNT=0.05912530     # Amount of BTC you own

# Update Frequency
CHECK_INTERVAL=1s         # How often to update calculations

# Optional: Metrics Port (default: 2112)
METRICS_PORT=2112         # Port for Prometheus metrics
```

## Continuous Integration & Deployment

> The application follows GitOps principles using a combination of GitHub Actions and ArgoCD.

### CI Pipeline (GitHub Actions)

The CI pipeline is triggered on:
- Push to `main` branch
- Pull requests targeting `main`

Pipeline stages:
1. **Build**
   - Builds the application
   - Performs code quality checks

2. **Container Image**
   - Builds Docker image
   - Tags with git SHA and version
   - Pushes to container registry

3. **Manifest Update**
   - Updates Kubernetes manifests with new image tag
   - Commits changes back to repository

### CD Pipeline (ArgoCD)

ArgoCD handles the continuous deployment by:
1. Monitoring the Kubernetes manifests in the repository
2. Automatically syncing changes to the cluster
3. Ensuring the desired state matches the actual state

#### ArgoCD Applications

- **Development**: Auto-sync enabled, automated pruning and self-healing
- **Staging**: Auto-sync enabled, requires manual approval for breaking changes
- **Production**: Manual sync required, follows strict change management

### Deployment Flow

1. Developer pushes code to `main`
2. GitHub Actions:
   - Builds code
   - Creates new container image
   - Updates K8s manifests with new image tag
   - Pushes changes to git

3. ArgoCD:
   - Detects manifest changes
   - Applies changes to cluster
   - Monitors deployment health
   - Performs rollback if needed

### Monitoring Deployments

- GitHub Actions: Check workflow status in `.github/workflows`
- ArgoCD Dashboard: Monitor sync status and health
- Kubernetes: Use `kubectl` to verify pod status

### Rollback Process

In case of deployment issues:
1. ArgoCD automatically detects health issues
2. For auto-sync environments: automatic rollback to last known good state
3. For manual environments: use ArgoCD UI/CLI to rollback

### Code Structure
1. Follow the established package structure
2. Keep the clean architecture principles
3. Use meaningful variable and function names
4. Add comments for complex logic

### Pending Improvements
- [ ] Add unit tests for all packages
- [ ] Add integration tests for API endpoints
- [ ] Implement error handling middleware
- [ ] Add request validation
- [ ] Improve logging with structured logs
- [ ] Add circuit breaker for external APIs
- [ ] Add cache layer for API responses
- [ ] Add rate limiting
- [ ] Add API documentation (Swagger/OpenAPI)

### Testing (To Be Implemented)
- Unit tests for business logic
- Integration tests for external APIs
- End-to-end tests for critical flows
- Performance tests for WebSocket connections
- Test coverage reports

### Best Practices
1. Write clean, maintainable code
2. Follow Go idioms and conventions
3. Handle errors appropriately
4. Add proper logging
5. Document public functions and types
6. Keep dependencies up to date
7. Use environment variables for configuration
8. Follow semantic versioning

### Pull Request Process
1. Create feature branch from `main`
2. Make your changes following the guidelines above
3. Update documentation if needed
4. Ensure CI pipeline passes
5. Request review from maintainers 
