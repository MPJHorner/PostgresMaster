# GitHub Actions Workflow Setup

## ⚠️ Important: Manual Workflow File Creation Required

Due to GitHub App permission restrictions, the workflow file `.github/workflows/test.yml` could not be pushed automatically. You need to create it manually via the GitHub UI.

## 📋 Instructions

### Option 1: Create via GitHub UI (Recommended)

1. Go to your repository on GitHub: https://github.com/MPJHorner/PostgresMaster
2. Navigate to the **Actions** tab
3. Click **"New workflow"** or **"set up a workflow yourself"**
4. Name the file: `test.yml`
5. Copy and paste the complete workflow content below
6. Commit the file to the `main` branch or your current branch

### Option 2: Create via GitHub Web Interface

1. Go to: https://github.com/MPJHorner/PostgresMaster/new/main?filename=.github/workflows/test.yml
2. Paste the workflow content below
3. Commit directly to your branch

### Option 3: Local Creation with Proper Git Credentials

If you have direct push access (not via GitHub App):

```bash
cd /home/user/PostgresMaster
# The file already exists locally at .github/workflows/test.yml
git add .github/workflows/test.yml
git commit -m "Add GitHub Actions CI/CD workflow"
git push origin claude/implement-prd-section-011CV2rrR6LEDD7cgzrj2uyv
```

## 📄 Complete Workflow Content

The workflow file is already created locally at:
```
/home/user/PostgresMaster/.github/workflows/test.yml
```

Or copy this content:

```yaml
name: Tests

on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main

jobs:
  test-proxy:
    name: Go Proxy Tests
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache-dependency-path: proxy/go.sum

      - name: Install dependencies
        working-directory: ./proxy
        run: go mod download

      - name: Run unit tests
        working-directory: ./proxy
        run: go test ./... -v -race -coverprofile=coverage.out

      - name: Display unit test coverage
        working-directory: ./proxy
        run: |
          echo "## Unit Test Coverage" >> $GITHUB_STEP_SUMMARY
          go tool cover -func=coverage.out | tail -1 >> $GITHUB_STEP_SUMMARY

      - name: Run integration tests
        working-directory: ./proxy
        env:
          TEST_POSTGRES_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
        run: |
          echo "Running integration tests with TEST_POSTGRES_URL..."
          go test ./... -v -race -coverprofile=coverage-integration.out

      - name: Display integration test coverage
        working-directory: ./proxy
        run: |
          echo "## Integration Test Coverage" >> $GITHUB_STEP_SUMMARY
          go tool cover -func=coverage-integration.out | tail -1 >> $GITHUB_STEP_SUMMARY

      - name: Check coverage threshold (80%)
        working-directory: ./proxy
        run: |
          coverage=$(go tool cover -func=coverage-integration.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Total coverage: ${coverage}%"
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "❌ Coverage ${coverage}% is below 80% threshold"
            exit 1
          else
            echo "✅ Coverage ${coverage}% meets 80% threshold"
          fi

      - name: Upload coverage reports
        uses: actions/upload-artifact@v4
        with:
          name: coverage-reports
          path: |
            proxy/coverage.out
            proxy/coverage-integration.out
          retention-days: 7

  lint-proxy:
    name: Go Proxy Linting
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache-dependency-path: proxy/go.sum

      - name: Install golangci-lint
        run: |
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

      - name: Run golangci-lint
        working-directory: ./proxy
        run: $(go env GOPATH)/bin/golangci-lint run --timeout=5m

  build-proxy:
    name: Build Proxy Binaries
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache-dependency-path: proxy/go.sum

      - name: Build for all platforms
        working-directory: ./proxy
        run: make build-all

      - name: Check binary sizes
        working-directory: ./proxy
        run: |
          echo "## Binary Sizes" >> $GITHUB_STEP_SUMMARY
          echo "\`\`\`" >> $GITHUB_STEP_SUMMARY
          ls -lh bin/ >> $GITHUB_STEP_SUMMARY
          echo "\`\`\`" >> $GITHUB_STEP_SUMMARY

          # Check that binaries are under 15MB (with some tolerance over the 10MB target)
          for binary in bin/*; do
            size=$(stat -f%z "$binary" 2>/dev/null || stat -c%s "$binary")
            size_mb=$((size / 1024 / 1024))
            echo "Binary $binary: ${size_mb}MB"
            if [ $size_mb -gt 15 ]; then
              echo "❌ Binary $binary (${size_mb}MB) exceeds 15MB limit"
              exit 1
            fi
          done
          echo "✅ All binaries are within size limits"

      - name: Upload binaries
        uses: actions/upload-artifact@v4
        with:
          name: proxy-binaries
          path: proxy/bin/*
          retention-days: 7
```

## ✅ What This Workflow Does

### Job 1: test-proxy (Go Proxy Tests)
- ✅ Runs on all PRs to `main` branch
- ✅ Sets up Postgres 16-alpine service container
- ✅ Executes unit tests with race detection
- ✅ Runs integration tests with real Postgres database
- ✅ **Enforces 80% coverage threshold** (fails CI if below)
- ✅ Uploads coverage reports as artifacts

### Job 2: lint-proxy (Code Quality)
- ✅ Installs and runs golangci-lint v1.55.2
- ✅ Comprehensive linting checks
- ✅ Fails on any linting errors

### Job 3: build-proxy (Cross-Platform Builds)
- ✅ Builds for all 5 platforms (Windows, macOS Intel/ARM, Linux x64/ARM)
- ✅ Validates binary sizes (must be < 15MB)
- ✅ Uploads binaries as artifacts

## 🧪 Testing the Workflow

Once the workflow is created:

1. Create a test PR or push to main
2. Go to the **Actions** tab
3. Watch all three jobs run in parallel
4. Verify all tests pass and coverage meets 80%

## 📊 Expected Results

After the workflow runs successfully, you'll see:
- ✅ All unit tests passing
- ✅ All integration tests passing with real Postgres
- ✅ Code coverage ≥ 80%
- ✅ All linting checks passing
- ✅ All platform binaries building successfully
- ✅ Binary sizes within limits

## 🔧 Troubleshooting

If the workflow fails:
- Check the Actions tab for detailed logs
- Unit tests should pass immediately
- Integration tests require the Postgres service to be healthy
- Linting may require code fixes
- Build issues may indicate platform-specific problems

## 📝 Notes

- The workflow file cannot be pushed via GitHub App without `workflows` permission
- This is a GitHub security feature to prevent unauthorized workflow modifications
- Once created via GitHub UI, future updates can be made through normal PRs
- The workflow is production-ready and follows GitHub Actions best practices

## ✅ Next Steps

1. Create the workflow file using one of the methods above
2. Verify the Actions tab shows the workflow
3. Create a test PR to trigger the workflow
4. Once verified, delete this `WORKFLOW_SETUP.md` file
