# Product Requirements Document: Browser-Only Postgres Client (MVP)

## Overview
A web-based SQL client for PostgreSQL that connects to remote Postgres databases through a lightweight local proxy. The web app is 100% frontend (static files), and users run a tiny Go binary locally that bridges browser WebSocket connections to Postgres TCP connections.

**Key Innovation:** Pure frontend web app with zero backend infrastructure. Users download a single ~5MB Go binary for their OS that runs locally, enabling secure connections to any remote Postgres server without exposing credentials to any third-party service.

## Goals
- Provide a lightweight, accessible Postgres client that runs in any modern browser
- Enable developers to write and execute SQL queries with intelligent autocomplete
- Display query results in a clear, readable format
- Maintain high code quality standards with proper tooling and testing

## Non-Goals (MVP)
- Complex result filtering or transformation
- Query visualization/charting
- Database schema modifications through UI
- Multi-user collaboration features
- Query performance analysis/EXPLAIN visualization

## Tech Stack

### Web Application (Frontend)
- **Framework**: SvelteKit (static adapter for pure static site)
- **UI Components**: shadcn-svelte (Svelte port of shadcn/ui)
- **Styling**: TailwindCSS (used by shadcn-svelte)
- **Build**: Vite (via SvelteKit)
- **Deployment**: Static hosting (Vercel, Netlify, GitHub Pages, etc.)

### SQL Editor
- **Editor**: Monaco Editor (VS Code's editor)
- **SQL Support**: Monaco SQL language support
- **Autocomplete**: Custom SQL language service with Postgres-specific keywords + schema-aware completions
- **Schema Introspection**: Query `pg_catalog` and `information_schema` for metadata

### Local Proxy (Go Binary)
- **Language**: Go 1.21+
- **Size**: ~5-10MB single binary
- **Platforms**: Windows, macOS (Intel + ARM), Linux (amd64 + arm64)
- **Communication**: WebSocket server (browser ↔ proxy)
- **Database**: Postgres wire protocol (proxy ↔ database)
- **Libraries**:
  - `gorilla/websocket` for WebSocket server
  - `jackc/pgx` for Postgres protocol
  - `spf13/cobra` for CLI (optional)
- **Features**:
  - CORS-enabled WebSocket server on localhost:8080
  - Postgres connection pooling
  - SSL/TLS support for Postgres connections
  - Query timeout handling
  - Graceful shutdown

### Browser Connection Flow
```
Web App (Browser)
    ↓ WebSocket (ws://localhost:8080?secret=xxx)
Go Proxy (Local)
    ↓ Postgres Protocol (TCP + SSL)
Remote Postgres Server
```

### WebSocket Protocol (Browser ↔ Proxy)
Messages are JSON-encoded with the following structure:

**Client → Proxy:**
```json
{
  "type": "query" | "introspect",
  "payload": {
    "sql": "SELECT * FROM users",  // for query type
    "timeout": 30000                // optional, milliseconds
  }
}
```

**Proxy → Client:**
```json
{
  "type": "result" | "error" | "schema",
  "payload": {
    "rows": [...],              // for result type
    "columns": [...],           // column metadata
    "rowCount": 10,
    "executionTime": 45,        // milliseconds
    "error": "...",             // for error type
    "schema": {...}             // for schema type
  }
}
```

### Data Storage
- **Query History**: Browser LocalStorage (in-memory for MVP, cleared on refresh)
- **Last Query**: Browser sessionStorage (persists on refresh during session)
- **Connection Info**: NOT stored in browser (proxy handles all connection details)
- **Credentials**: NEVER touch the browser - only handled by proxy

### Code Quality
- **Formatter**: Prettier
- **Linter**: ESLint with popular config (eslint-config-airbnb-base or standard)
- **Type Safety**: TypeScript
- **Testing**: Vitest + Svelte Testing Library
- **Coverage Target**: 80% minimum

## Core Features

### 1. SQL Query Editor
- Monaco-based editor with SQL syntax highlighting
- Intelligent autocomplete for:
  - SQL keywords (SELECT, FROM, WHERE, etc.)
  - Database objects (tables, columns, functions)
  - Postgres-specific syntax
- Multi-line query support
- Execute query button (Ctrl+Enter shortcut)

### 2. Schema Introspection
- Automatically fetch database schema on connection
- Cache schema information for autocomplete
- Include:
  - Table names
  - Column names and types
  - Function names
- Manual refresh capability

### 3. Query Results Display
- Tabular display of query results
- Column headers with data types
- Row count display
- Basic result set information (execution time, rows affected)
- Error display for failed queries

### 4. Proxy Authentication
- Secret-based authentication to local proxy
- Secret passed via URL parameter
- WebSocket messages include secret in headers
- Invalid secret → connection refused
- Proxy validates secret on every request

## User Flow

### First-Time Setup
1. **Download Proxy**
   - User visits `yourapp.com`
   - Landing page prompts: "Download proxy for your OS"
   - User downloads appropriate binary:
     - `postgres-proxy-windows.exe` (Windows)
     - `postgres-proxy-darwin-amd64` (macOS Intel)
     - `postgres-proxy-darwin-arm64` (macOS Apple Silicon)
     - `postgres-proxy-linux-amd64` (Linux)
   - Instructions shown on page

2. **Connect via Proxy**

   **Option A: Connection String**
   ```bash
   ./postgres-proxy "postgres://user:pass@host:5432/dbname"
   ```

   **Option B: Interactive Mode**
   ```bash
   ./postgres-proxy
   Host: db.example.com
   Port [5432]:
   Database: mydb
   Username: admin
   Password: ****
   SSL Mode [prefer]: require

   Connecting to postgres://db.example.com:5432/mydb...
   Retrying... (attempt 2)
   ✓ Connected!

   → Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6
   ```

3. **Open Web App**
   - User clicks or copies URL
   - Browser opens with secret parameter
   - Web app auto-connects to proxy using secret
   - Schema introspection runs automatically
   - Editor becomes active - ready to query!

### Regular Usage
4. **Write Query**
   - User types in SQL editor
   - Autocomplete suggests:
     - SQL keywords (SELECT, FROM, WHERE, JOIN, etc.)
     - Table names from connected database
     - Column names for referenced tables
     - Postgres-specific functions and syntax
   - User completes query

5. **Execute Query**
   - User clicks "Run" or presses Ctrl+Enter
   - Query → WebSocket (with secret auth) → Proxy → Postgres
   - Results streamed back through same path
   - Results display in table below editor
   - Execution time and row count shown

### Subsequent Sessions
6. **Reconnect**
   - User runs proxy again with connection details
   - Proxy generates new secret and URL
   - User opens new URL in browser
   - Continue querying

## Development Standards

### Code Style
- ESLint with Airbnb base configuration (adapted for Svelte)
- Prettier for formatting
- Pre-commit hooks with Husky + lint-staged

### Testing Strategy
- **Unit Tests**: Component logic, utilities
- **Integration Tests**: Query execution, schema introspection
- **E2E Tests**: Critical user flows (connect → query → results)
- Minimum 80% code coverage

### Project Structure
```
postgres-client/
├── web/                    # Frontend web app
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   ├── ui/          # shadcn-svelte components
│   │   │   │   ├── Editor.svelte
│   │   │   │   ├── Results.svelte
│   │   │   │   └── Connection.svelte
│   │   │   ├── services/
│   │   │   │   ├── websocket.ts  # WebSocket client
│   │   │   │   ├── protocol.ts   # Message protocol
│   │   │   │   └── introspection.ts
│   │   │   ├── stores/
│   │   │   │   └── connection.ts
│   │   │   └── utils/
│   │   ├── routes/
│   │   │   └── +page.svelte
│   │   └── tests/
│   ├── package.json
│   ├── svelte.config.js
│   └── vite.config.ts
│
├── proxy/                  # Go proxy binary
│   ├── cmd/
│   │   └── proxy/
│   │       └── main.go     # Entry point
│   ├── pkg/
│   │   ├── server/         # WebSocket server
│   │   ├── postgres/       # Postgres client
│   │   ├── protocol/       # Message protocol
│   │   └── config/         # Configuration
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile            # Build for all platforms
│   └── README.md
│
└── docs/
    ├── installation.md
    └── architecture.md
```

## Implementation Decisions

### Confirmed for MVP
- ✅ Go proxy binary (WebSocket ↔ Postgres bridge)
- ✅ SvelteKit web app (static build)
- ✅ Single query editor (no tabs)
- ✅ Display all query results (no pagination for MVP)
- ✅ Connection saving in LocalStorage (optional)
- ✅ Query history: in-memory only (cleared on refresh)
- ✅ Cross-platform binaries (Windows, macOS, Linux)

### Deferred to Post-MVP
- Auto-update mechanism for proxy binary
- Export results to CSV/JSON
- Query history persistence across sessions
- Multiple query tabs
- Result pagination for large datasets
- Query visualization/charting
- Proxy configuration options (custom port, HTTPS)

## Success Metrics (MVP)
- User can download and run proxy binary for their OS
- User can connect to any remote Postgres database through web app
- User can write queries with intelligent autocomplete
- User can execute queries and see formatted results
- Proxy binary < 10MB
- Web app loads in < 2 seconds
- 80%+ test coverage (both web app and proxy)
- Zero critical linting errors
- All tests passing

## Timeline (Estimated)

### Phase 1: Project Setup & Go Proxy (3 days)
- Initialize monorepo structure
- Go proxy: WebSocket server + Postgres client
- Message protocol definition
- Build system for cross-platform binaries
- Basic integration tests

### Phase 2: Web App Foundation (2 days)
- SvelteKit setup with static adapter
- shadcn-svelte integration
- WebSocket client service
- Connection form component
- Test connection flow

### Phase 3: SQL Editor & Autocomplete (3 days)
- Monaco editor integration
- SQL syntax highlighting
- Postgres keyword autocomplete
- Schema introspection service
- Schema-aware autocomplete (tables, columns)

### Phase 4: Query Execution & Results (2 days)
- Query execution flow
- Results table component
- Error handling and display
- Execution metadata (time, row count)

### Phase 5: Testing & Polish (2 days)
- E2E tests (web app + proxy)
- Code quality checks (linting, formatting)
- Documentation (README, installation guide)
- Binary release workflow

**Total: ~12 days for MVP**

## Future Enhancements (Post-MVP)
- Query history and saved queries
- Result filtering and sorting
- Export to CSV/JSON
- Multiple query tabs
- Dark/light theme toggle
- Query snippets/templates
- Keyboard shortcuts panel
# Master Implementation Plan - Browser-Only Postgres Client MVP

> **Single source of truth for building the MVP. Check off items as completed.**

---

## 🎯 Project Overview

**Goal**: Build a browser-based Postgres client that requires zero backend infrastructure.

**Architecture**:
- Web App (100% frontend): SvelteKit + shadcn-svelte + Monaco Editor
- Go Proxy (~5MB): WebSocket ↔ Postgres TCP bridge

**Timeline**: 12 days • **Target Coverage**: 80%+ • **License**: AGPL-3.0

---

## 📋 Phase 1: Go Proxy Foundation (Days 1-3)

### 1.1 Project Structure Setup

- [ ] Create monorepo root structure
- [ ] Create `proxy/` directory
- [ ] Create `web/` directory
- [ ] Create `docs/` directory
- [ ] Add root `.gitignore` (Go + Node.js)
- [ ] Add root `README.md` with project overview
- [ ] Add `LICENSE` file (AGPL-3.0)

### 1.2 Go Proxy - Initial Setup

- [ ] Create `proxy/go.mod`: `go mod init github.com/youruser/postgres-client/proxy`
- [ ] Create directory structure:
  ```
  proxy/
  ├── cmd/proxy/main.go
  ├── pkg/
  │   ├── auth/secret.go
  │   ├── postgres/client.go
  │   ├── protocol/messages.go
  │   └── server/websocket.go
  ├── go.mod
  ├── go.sum
  ├── Makefile
  └── README.md
  ```
- [ ] Install dependencies:
  - [ ] `go get github.com/gorilla/websocket`
  - [ ] `go get github.com/jackc/pgx/v5`
  - [ ] `go get github.com/jackc/pgx/v5/pgxpool`

### 1.3 Secret Generation & Validation

**File**: `proxy/pkg/auth/secret.go`

- [ ] Implement `GenerateSecret()` function (32 bytes, hex encoded = 64 chars)
- [ ] Implement `ValidateSecret(secret string)` function
- [ ] Add unit tests:
  - [ ] Test secret length is 64 characters
  - [ ] Test secret is hex encoded
  - [ ] Test validation accepts valid secret
  - [ ] Test validation rejects invalid secret

### 1.4 Message Protocol Types

**File**: `proxy/pkg/protocol/messages.go`

- [ ] Define `Message` struct with ID, Type, Payload
- [ ] Define `ClientMessage` types: query, introspect, ping
- [ ] Define `ServerMessage` types: result, error, schema, pong
- [ ] Define payload structs:
  - [ ] `QueryPayload` (sql, params, timeout)
  - [ ] `ResultPayload` (rows, columns, rowCount, executionTime)
  - [ ] `ErrorPayload` (code, message, detail, position)
  - [ ] `SchemaPayload` (tables, functions)
- [ ] Add JSON serialization tags
- [ ] Add unit tests:
  - [ ] Test message serialization to JSON
  - [ ] Test message deserialization from JSON

### 1.5 Connection String Parser & Interactive Mode

**File**: `proxy/cmd/proxy/main.go`

- [ ] Implement connection string parsing from CLI args
- [ ] Implement interactive mode with prompts:
  - [ ] Host prompt
  - [ ] Port prompt (default: 5432)
  - [ ] Database prompt
  - [ ] Username prompt
  - [ ] Password prompt (hidden input)
  - [ ] SSL mode prompt (default: prefer)
- [ ] Build connection string from interactive inputs
- [ ] Add validation for required fields
- [ ] Test both modes work correctly

### 1.6 Postgres Connection with Retry Logic

**File**: `proxy/pkg/postgres/client.go`

- [ ] Implement `NewClient(connString)` function
- [ ] Implement connection with pgxpool
- [ ] Implement retry logic with exponential backoff:
  - [ ] Attempt 1: immediate
  - [ ] Attempt 2: wait 2s
  - [ ] Attempt 3: wait 4s
  - [ ] Attempt 4: wait 8s
  - [ ] Fail with clear error message
- [ ] Add connection pool config (MaxConns: 5, MinConns: 1)
- [ ] Implement `Close()` method
- [ ] Add graceful shutdown on SIGINT/SIGTERM
- [ ] Add unit tests (mock Postgres):
  - [ ] Test successful connection
  - [ ] Test retry on connection failure
  - [ ] Test failure after max retries

### 1.7 Query Execution

**File**: `proxy/pkg/postgres/client.go`

- [ ] Implement `ExecuteQuery(ctx, sql, params)` function
- [ ] Parse query results from pgx.Rows
- [ ] Convert column types to JSON-friendly format
- [ ] Handle NULL values correctly
- [ ] Measure execution time
- [ ] Return structured QueryResult
- [ ] Add error handling for:
  - [ ] Syntax errors
  - [ ] Permission errors
  - [ ] Timeout errors
- [ ] Add unit tests:
  - [ ] Test simple SELECT query
  - [ ] Test query with parameters
  - [ ] Test query with NULL values
  - [ ] Test query timeout

### 1.8 Schema Introspection

**File**: `proxy/pkg/postgres/client.go`

- [ ] Implement `IntrospectSchema(ctx)` function
- [ ] Query for tables:
  ```sql
  SELECT n.nspname, c.relname, c.relkind
  FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE c.relkind IN ('r','v','m') AND n.nspname NOT IN ('pg_catalog','information_schema')
  ```
- [ ] Query for columns per table:
  ```sql
  SELECT a.attname, format_type(a.atttypid, a.atttypmod), a.attnotnull
  FROM pg_attribute a WHERE a.attrelid = $1::regclass AND a.attnum > 0
  ```
- [ ] Query for functions:
  ```sql
  SELECT n.nspname, p.proname, pg_get_function_result(p.oid)
  FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
  WHERE n.nspname NOT IN ('pg_catalog','information_schema')
  ```
- [ ] Return structured SchemaInfo
- [ ] Add unit tests with test database

### 1.9 WebSocket Server

**File**: `proxy/pkg/server/websocket.go`

- [ ] Implement WebSocket upgrader with CORS
- [ ] Implement secret validation on connection
- [ ] Implement connection handler
- [ ] Implement message router (switch on message type)
- [ ] Route "query" messages to Postgres client
- [ ] Route "introspect" messages to schema introspection
- [ ] Route "ping" messages to pong response
- [ ] Stream results back to client
- [ ] Handle connection errors gracefully
- [ ] Add request/response correlation (message IDs)
- [ ] Add unit tests:
  - [ ] Test WebSocket upgrade
  - [ ] Test secret validation (valid/invalid)
  - [ ] Test message routing

### 1.10 Main Entry Point

**File**: `proxy/cmd/proxy/main.go`

- [ ] Wire up all components in main()
- [ ] Parse CLI args or start interactive mode
- [ ] Connect to Postgres with retry
- [ ] Generate secret
- [ ] Start WebSocket server on :8080
- [ ] Print success message with URL
- [ ] Handle graceful shutdown
- [ ] Add logging for debugging
- [ ] Test full flow end-to-end

### 1.11 Build System

**File**: `proxy/Makefile`

- [ ] Add `build-all` target for all platforms
- [ ] Add individual build targets:
  - [ ] `build-windows-amd64` → `postgres-proxy-windows-amd64.exe`
  - [ ] `build-darwin-amd64` → `postgres-proxy-darwin-amd64`
  - [ ] `build-darwin-arm64` → `postgres-proxy-darwin-arm64`
  - [ ] `build-linux-amd64` → `postgres-proxy-linux-amd64`
  - [ ] `build-linux-arm64` → `postgres-proxy-linux-arm64`
- [ ] Add `test` target: `go test ./... -v -race -coverprofile=coverage.out`
- [ ] Add `lint` target: `golangci-lint run`
- [ ] Add `clean` target: remove bin/ directory
- [ ] Add `dev` target: `go run cmd/proxy/main.go`
- [ ] Test building for all platforms
- [ ] Verify binary sizes are < 10MB

### 1.12 Go Testing & Coverage

- [ ] Write unit tests for all packages
- [ ] Run tests: `go test ./... -v`
- [ ] Generate coverage report: `go test ./... -coverprofile=coverage.out`
- [ ] Open coverage HTML: `go tool cover -html=coverage.out`
- [ ] Verify 80%+ coverage
- [ ] Fix any failing tests
- [ ] Document test commands in proxy/README.md

### 1.13 Integration Testing (Proxy → Postgres)

**File**: `proxy/pkg/postgres/integration_test.go`

- [ ] Set up test Postgres container (Docker)
- [ ] Test full query execution flow
- [ ] Test schema introspection
- [ ] Test connection retry logic
- [ ] Test error handling
- [ ] Clean up containers after tests

### 1.14 Phase 1 - Manual Testing

- [ ] Build proxy for your OS
- [ ] Start test Postgres: `docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=test postgres`
- [ ] Run proxy with connection string
- [ ] Verify proxy connects successfully
- [ ] Verify secret is generated (64 chars)
- [ ] Verify URL is printed
- [ ] Test WebSocket connection with wscat or custom script
- [ ] Send test query message
- [ ] Verify results are returned
- [ ] Test schema introspection
- [ ] Verify error handling
- [ ] Stop proxy gracefully with Ctrl+C

**✅ Phase 1 Complete - Go proxy working end-to-end**

---

## 📋 Phase 2: Web App Foundation (Days 4-5)

### 2.1 SvelteKit Initialization

- [ ] Navigate to project root
- [ ] Create web directory: `mkdir -p web && cd web`
- [ ] Initialize SvelteKit: `npm create svelte@latest .`
  - [ ] Choose: Skeleton project
  - [ ] TypeScript: Yes
  - [ ] ESLint: Yes
  - [ ] Prettier: Yes
  - [ ] Playwright: No
  - [ ] Vitest: Yes
- [ ] Install dependencies: `npm install`
- [ ] Test dev server: `npm run dev`

### 2.2 Static Adapter Configuration

- [ ] Install adapter: `npm install -D @sveltejs/adapter-static`
- [ ] Configure `svelte.config.js`:
  ```js
  import adapter from '@sveltejs/adapter-static';

  export default {
    kit: {
      adapter: adapter({
        pages: 'build',
        assets: 'build',
        fallback: 'index.html'
      })
    }
  };
  ```
- [ ] Test build: `npm run build`
- [ ] Test preview: `npm run preview`
- [ ] Verify static files in `build/`

### 2.3 TailwindCSS Setup

- [ ] Install: `npm install -D tailwindcss postcss autoprefixer`
- [ ] Initialize: `npx tailwindcss init -p`
- [ ] Configure `tailwind.config.js`:
  ```js
  export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: { extend: {} },
    plugins: []
  };
  ```
- [ ] Create `src/app.css`:
  ```css
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
  ```
- [ ] Import in `src/routes/+layout.svelte`:
  ```svelte
  <script>
    import '../app.css';
  </script>
  <slot />
  ```
- [ ] Test Tailwind classes work

### 2.4 shadcn-svelte Integration

- [ ] Install: `npx shadcn-svelte@latest init`
- [ ] Follow prompts for configuration
- [ ] Add components:
  - [ ] `npx shadcn-svelte@latest add button`
  - [ ] `npx shadcn-svelte@latest add input`
  - [ ] `npx shadcn-svelte@latest add card`
  - [ ] `npx shadcn-svelte@latest add table`
  - [ ] `npx shadcn-svelte@latest add alert`
  - [ ] `npx shadcn-svelte@latest add badge`
  - [ ] `npx shadcn-svelte@latest add separator`
- [ ] Verify components render correctly
- [ ] Create test page using components

### 2.5 WebSocket Protocol Types (TypeScript)

**File**: `web/src/lib/services/protocol.ts`

- [ ] Define `ClientMessage` interface
- [ ] Define `ServerMessage` interface
- [ ] Define `QueryResult` interface
- [ ] Define `ColumnInfo` interface
- [ ] Define `SchemaInfo` interface
- [ ] Define `TableInfo` interface
- [ ] Define `ErrorPayload` interface
- [ ] Match Go message types exactly
- [ ] Export all types

### 2.6 WebSocket Client Service

**File**: `web/src/lib/services/websocket.ts`

- [ ] Install uuid: `npm install uuid && npm install -D @types/uuid`
- [ ] Create `PostgresProxyClient` class
- [ ] Implement constructor with secret
- [ ] Implement `connect()` method
- [ ] Implement WebSocket message handling
- [ ] Implement request/response correlation with message IDs
- [ ] Implement `executeQuery(sql, params)` method
- [ ] Implement `introspectSchema()` method
- [ ] Implement `close()` method
- [ ] Add error handling and reconnection logic
- [ ] Add TypeScript types for all methods
- [ ] Add JSDoc comments

### 2.7 Connection Store

**File**: `web/src/lib/stores/connection.ts`

- [ ] Create `ConnectionState` interface
- [ ] Create `connectionStore` writable store
- [ ] Create `isConnected` derived store
- [ ] Create `client` derived store
- [ ] Add helper functions for state updates
- [ ] Export all stores

### 2.8 Schema Store

**File**: `web/src/lib/stores/schema.ts`

- [ ] Create `schemaStore` writable store
- [ ] Create `tables` derived store
- [ ] Create `tableNames` derived store
- [ ] Create `allColumns` derived store (for autocomplete)
- [ ] Export all stores

### 2.9 Landing Page Component

**File**: `web/src/lib/components/LandingPage.svelte`

- [ ] Create component structure
- [ ] Add hero section with title and description
- [ ] Add features section (3 cards):
  - [ ] "Any Postgres" card
  - [ ] "Lightning Fast" card
  - [ ] "Secure" card
- [ ] Add "Get Started" section with 3 steps
- [ ] Add download buttons for all platforms:
  - [ ] Windows (x64)
  - [ ] macOS (Intel)
  - [ ] macOS (Apple Silicon)
  - [ ] Linux (x64)
- [ ] Add code examples for usage
- [ ] Style with shadcn components
- [ ] Make responsive for mobile/desktop
- [ ] Add footer with tech stack info

### 2.10 Connection Status Component

**File**: `web/src/lib/components/ConnectionStatus.svelte`

- [ ] Create component with connection state display
- [ ] Show "Connected" badge when connected
- [ ] Show "Connecting..." when connecting
- [ ] Show "Disconnected" with retry when disconnected
- [ ] Show error message on connection error
- [ ] Style with shadcn Badge component
- [ ] Make it prominent in the UI

### 2.11 Main Page Component

**File**: `web/src/routes/+page.svelte`

- [ ] Parse `?secret` parameter from URL on mount
- [ ] Show `LandingPage` if no secret
- [ ] Show connecting state if secret present
- [ ] Implement connection logic with proxy
- [ ] Show `ConnectionStatus` component
- [ ] Show error if connection fails
- [ ] On successful connection:
  - [ ] Introspect schema
  - [ ] Store in schema store
  - [ ] Show editor UI (placeholder for now)
- [ ] Add proper error boundaries

### 2.12 Phase 2 - Manual Testing

- [ ] Start proxy with test database
- [ ] Copy secret from proxy output
- [ ] Run web app: `npm run dev`
- [ ] Visit `http://localhost:5173` without secret
  - [ ] Verify landing page shows
  - [ ] Verify download buttons present
  - [ ] Verify instructions clear
- [ ] Visit `http://localhost:5173?secret=invalidsecret`
  - [ ] Verify connection fails with error
- [ ] Visit `http://localhost:5173?secret=<valid-secret>`
  - [ ] Verify "Connecting..." shows
  - [ ] Verify connection succeeds
  - [ ] Verify "Connected" badge shows
  - [ ] Verify no errors in console
- [ ] Test on mobile viewport
- [ ] Build and test production: `npm run build && npm run preview`

**✅ Phase 2 Complete - Web app connects to proxy**

---

## 📋 Phase 3: SQL Editor & Autocomplete (Days 6-8)

### 3.1 Monaco Editor Installation

- [ ] Install Monaco: `npm install monaco-editor`
- [ ] Install Vite plugin: `npm install -D vite-plugin-monaco-editor`
- [ ] Configure `vite.config.ts`:
  ```ts
  import monacoEditorPlugin from 'vite-plugin-monaco-editor';

  export default {
    plugins: [
      sveltekit(),
      monacoEditorPlugin({ languageWorkers: ['editorWorkerService'] })
    ]
  };
  ```
- [ ] Test Monaco imports work

### 3.2 Editor Component

**File**: `web/src/lib/components/Editor.svelte`

- [ ] Create component with editor container div
- [ ] Import Monaco Editor
- [ ] Initialize editor in onMount:
  - [ ] Language: 'sql'
  - [ ] Theme: 'vs-dark'
  - [ ] Options: minimap disabled, auto layout, word wrap
- [ ] Bind value prop (two-way binding)
- [ ] Add onChange callback prop
- [ ] Add onExecute callback prop
- [ ] Register Ctrl+Enter command for execution
- [ ] Dispose editor onDestroy
- [ ] Style editor container (height, border, rounded)
- [ ] Test editor renders and accepts input

### 3.3 SQL Keywords List

**File**: `web/src/lib/utils/autocomplete.ts`

- [ ] Create comprehensive Postgres keywords array:
  - [ ] SELECT, FROM, WHERE, JOIN, etc.
  - [ ] INSERT, UPDATE, DELETE, etc.
  - [ ] CREATE, DROP, ALTER, etc.
  - [ ] Postgres-specific: RETURNING, CONFLICT, JSONB, UUID, etc.
  - [ ] Data types: INTEGER, TEXT, TIMESTAMP, etc.
  - [ ] Constraints: PRIMARY KEY, FOREIGN KEY, etc.
- [ ] Create functions array:
  - [ ] COUNT, SUM, AVG, MIN, MAX
  - [ ] NOW, CURRENT_TIMESTAMP, etc.
  - [ ] String functions: CONCAT, LOWER, UPPER, etc.
  - [ ] JSON functions: JSON_AGG, JSONB_AGG, etc.
  - [ ] Window functions: ROW_NUMBER, RANK, etc.

### 3.4 Basic Keyword Autocomplete

**File**: `web/src/lib/utils/autocomplete.ts`

- [ ] Create `setupAutocomplete(monaco, schema)` function
- [ ] Register completion item provider for 'sql'
- [ ] Implement `provideCompletionItems` callback
- [ ] Get word at cursor position
- [ ] Calculate range for replacement
- [ ] Generate keyword suggestions:
  - [ ] Map keywords to CompletionItem
  - [ ] Set kind: Keyword
  - [ ] Set sortText to prioritize keywords
- [ ] Generate function suggestions:
  - [ ] Map functions to CompletionItem
  - [ ] Set kind: Function
  - [ ] Add snippet with $1 placeholder
- [ ] Return suggestions
- [ ] Test autocomplete triggers on typing

### 3.5 Schema-Aware Autocomplete

**File**: `web/src/lib/utils/autocomplete.ts`

- [ ] Add table name suggestions:
  - [ ] Map tables from schema to CompletionItem
  - [ ] Set kind: Class
  - [ ] Add detail with schema name
  - [ ] Add documentation with column list
- [ ] Add column name suggestions:
  - [ ] Map columns from all tables
  - [ ] Format as `table.column`
  - [ ] Set kind: Field
  - [ ] Add detail with data type
- [ ] Add function suggestions from schema:
  - [ ] Map functions to CompletionItem
  - [ ] Add return type as detail
- [ ] Sort suggestions: keywords → functions → tables → columns
- [ ] Test autocomplete shows schema objects

### 3.6 Context-Aware Autocomplete (Advanced)

**File**: `web/src/lib/utils/autocomplete.ts`

- [ ] Parse SQL context at cursor position
- [ ] Detect if inside FROM clause → suggest tables
- [ ] Detect if after table name + dot → suggest columns for that table
- [ ] Detect if inside WHERE clause → suggest columns
- [ ] Add smarter filtering based on context
- [ ] Test context detection works

### 3.7 Query Panel Component

**File**: `web/src/lib/components/QueryPanel.svelte`

- [ ] Import Editor component
- [ ] Create layout with editor at top
- [ ] Add "Run Query" button with Ctrl+Enter hint
- [ ] Add SQL state variable
- [ ] Implement executeQuery function (placeholder for now)
- [ ] Add loading state
- [ ] Add error display area
- [ ] Add results display area (placeholder)
- [ ] Style with shadcn components
- [ ] Wire up to connection store

### 3.8 Update Main Page with Editor

**File**: `web/src/routes/+page.svelte`

- [ ] Import QueryPanel component
- [ ] Show QueryPanel when connected
- [ ] Pass client from connection store
- [ ] Remove placeholder editor UI
- [ ] Test full flow: landing → connect → editor shows

### 3.9 Phase 3 - Manual Testing

- [ ] Start proxy and web app
- [ ] Connect to proxy
- [ ] Verify editor loads
- [ ] Type "SEL" → verify "SELECT" suggested
- [ ] Type "SELECT * FROM " → verify table names suggested
- [ ] Create a test table in Postgres
- [ ] Refresh schema
- [ ] Type "SELECT * FROM test_table." → verify columns suggested
- [ ] Type "SELECT cou" → verify "COUNT" suggested with ()
- [ ] Test Ctrl+Enter (should do nothing yet, just verify command registered)
- [ ] Test editor styling looks good
- [ ] Test autocomplete on mobile (should work)

**✅ Phase 3 Complete - SQL editor with intelligent autocomplete**

---

## 📋 Phase 4: Query Execution & Results (Days 9-10)

### 4.1 Query Execution Logic

**File**: `web/src/lib/components/QueryPanel.svelte`

- [ ] Implement executeQuery function:
  - [ ] Validate SQL is not empty
  - [ ] Set loading state
  - [ ] Clear previous results/errors
  - [ ] Call client.executeQuery(sql)
  - [ ] Handle success: store results
  - [ ] Handle error: display error
  - [ ] Clear loading state
- [ ] Add try/catch with proper error handling
- [ ] Add loading spinner during execution
- [ ] Test query execution end-to-end

### 4.2 Value Formatting Utility

**File**: `web/src/lib/utils/format.ts`

- [ ] Create `formatValue(value)` function
- [ ] Handle null/undefined → "NULL"
- [ ] Handle boolean → "true"/"false"
- [ ] Handle number → toString()
- [ ] Handle string → as-is
- [ ] Handle Date → toISOString()
- [ ] Handle Array → JSON.stringify
- [ ] Handle Object → JSON.stringify with pretty print
- [ ] Add unit tests for all data types

### 4.3 Results Table Component

**File**: `web/src/lib/components/Results.svelte`

- [ ] Import shadcn Table components
- [ ] Accept `data: QueryResult` prop
- [ ] Display metadata section:
  - [ ] Row count
  - [ ] Execution time (ms)
- [ ] Render table:
  - [ ] TableHeader with column names + types as badges
  - [ ] TableBody with rows
  - [ ] Format each cell value with formatValue()
  - [ ] Handle NULL values specially (gray text)
- [ ] Add scrolling for large result sets (max height 500px)
- [ ] Add sticky header when scrolling
- [ ] Handle empty results (show message)
- [ ] Style with Card component
- [ ] Test with various data types

### 4.4 Error Display Component

**File**: `web/src/lib/components/ErrorDisplay.svelte`

- [ ] Import shadcn Alert components
- [ ] Accept `error: string` prop
- [ ] Parse Postgres error if possible:
  - [ ] Extract error code if present
  - [ ] Extract position if present
- [ ] Display error with AlertCircle icon
- [ ] Use destructive variant
- [ ] Show error code if available
- [ ] Show hint/detail if available
- [ ] Make error message monospace font
- [ ] Test with various Postgres errors

### 4.5 Wire Up Results Display

**File**: `web/src/lib/components/QueryPanel.svelte`

- [ ] Import Results component
- [ ] Import ErrorDisplay component
- [ ] Add results state variable
- [ ] Add error state variable
- [ ] Show Results component when results available
- [ ] Show ErrorDisplay when error present
- [ ] Show loading spinner when executing
- [ ] Show empty state when no results yet
- [ ] Test full flow: type query → execute → see results

### 4.6 Query History (In-Memory)

**File**: `web/src/lib/stores/history.ts`

- [ ] Create `QueryHistoryItem` interface (sql, timestamp, success)
- [ ] Create `queryHistory` writable store
- [ ] Create `addToHistory(item)` function
- [ ] Keep last 50 queries
- [ ] Export store and functions

**File**: `web/src/lib/components/QueryPanel.svelte`

- [ ] Add queries to history on execution
- [ ] (Optional) Display recent queries in sidebar/dropdown
- [ ] (Optional) Click to load query into editor

### 4.7 Sample Queries for Testing

Create test queries file for manual testing:

**File**: `web/test-queries.sql`

- [ ] Simple SELECT: `SELECT 1 as num`
- [ ] SELECT with multiple columns and types
- [ ] SELECT with NULL values
- [ ] SELECT with timestamps
- [ ] SELECT with JSON/JSONB
- [ ] SELECT large result set: `SELECT * FROM generate_series(1, 1000)`
- [ ] INSERT query
- [ ] UPDATE query
- [ ] DELETE query
- [ ] Query with syntax error
- [ ] Query on non-existent table

### 4.8 Phase 4 - Manual Testing

- [ ] Start proxy and web app
- [ ] Connect to proxy
- [ ] Execute simple query: `SELECT 1 as num`
  - [ ] Verify result shows: 1
  - [ ] Verify row count: 1 row
  - [ ] Verify execution time shown
- [ ] Execute query with multiple types:
  ```sql
  SELECT
    1 as int_col,
    'text' as text_col,
    true as bool_col,
    NULL as null_col,
    NOW() as timestamp_col
  ```
  - [ ] Verify all types display correctly
  - [ ] Verify NULL shows as "NULL"
  - [ ] Verify timestamp formatted
- [ ] Execute large result:
  ```sql
  SELECT * FROM generate_series(1, 1000) as n
  ```
  - [ ] Verify scrolling works
  - [ ] Verify header sticky
  - [ ] Verify performance acceptable
- [ ] Execute query with error:
  ```sql
  SELECT * FROM nonexistent_table
  ```
  - [ ] Verify error displays
  - [ ] Verify error message clear
- [ ] Execute INSERT/UPDATE/DELETE
  - [ ] Verify success message
  - [ ] Verify affected rows shown
- [ ] Test on mobile viewport
- [ ] Build and test production build

**✅ Phase 4 Complete - Full query execution with results display**

---

## 📋 Phase 5: Testing & Polish (Days 11-12)

### 5.1 Web App Unit Testing

**Files**: `web/src/lib/**/*.test.ts`

- [ ] Install testing libraries (already installed in Phase 2)
- [ ] Write tests for `format.ts`:
  - [ ] Test formatValue with all data types
  - [ ] Test NULL handling
  - [ ] Test JSON formatting
- [ ] Write tests for `websocket.ts`:
  - [ ] Test message serialization
  - [ ] Test request/response correlation
  - [ ] Test error handling
- [ ] Write tests for stores:
  - [ ] Test connection store state transitions
  - [ ] Test schema store derived values
- [ ] Write component tests:
  - [ ] Test Editor component renders
  - [ ] Test Results component with mock data
  - [ ] Test ErrorDisplay component
- [ ] Run tests: `npm test`
- [ ] Generate coverage: `npm run test:coverage`
- [ ] Verify 80%+ coverage

### 5.2 Go Proxy Final Testing

- [ ] Run all tests: `go test ./... -v`
- [ ] Generate coverage: `go test ./... -coverprofile=coverage.out`
- [ ] View coverage: `go tool cover -html=coverage.out`
- [ ] Verify 80%+ coverage
- [ ] Fix any failing tests
- [ ] Test with race detector: `go test ./... -race`
- [ ] Fix any race conditions

### 5.3 Code Quality - Linting

**Go Proxy**:
- [ ] Install golangci-lint (if not installed)
- [ ] Run linter: `golangci-lint run`
- [ ] Fix all errors
- [ ] Fix all warnings
- [ ] Run gofmt: `go fmt ./...`
- [ ] Commit formatted code

**Web App**:
- [ ] Run ESLint: `npm run lint`
- [ ] Fix all errors
- [ ] Fix all warnings
- [ ] Run Prettier: `npm run format`
- [ ] Commit formatted code

### 5.4 Code Quality - Type Checking

- [ ] Run TypeScript check: `npm run check`
- [ ] Fix all type errors
- [ ] Ensure no `any` types (or minimal)
- [ ] Add JSDoc comments to public APIs

### 5.5 End-to-End Testing

**Test Scenario 1: Fresh Install**
- [ ] Build proxy for your OS
- [ ] Start fresh Postgres in Docker
- [ ] Run proxy with interactive mode
- [ ] Open URL in browser
- [ ] Verify full flow works
- [ ] Create table, insert data, query

**Test Scenario 2: Connection Errors**
- [ ] Stop Postgres
- [ ] Try to connect with proxy
- [ ] Verify retry logic works
- [ ] Verify clear error message
- [ ] Restart Postgres
- [ ] Verify proxy can connect

**Test Scenario 3: Invalid Secret**
- [ ] Start proxy
- [ ] Modify secret in URL
- [ ] Try to connect from browser
- [ ] Verify connection rejected
- [ ] Verify clear error message

**Test Scenario 4: Large Queries**
- [ ] Execute query with 10,000 rows
- [ ] Verify results display
- [ ] Verify scrolling smooth
- [ ] Verify no memory leaks

**Test Scenario 5: Cross-Browser**
- [ ] Test in Chrome
- [ ] Test in Firefox
- [ ] Test in Safari
- [ ] Test on mobile (iOS Safari / Chrome)
- [ ] Verify all features work

**Test Scenario 6: Build & Deploy**
- [ ] Build proxy for all platforms
- [ ] Verify all binaries < 10MB
- [ ] Test each binary on respective OS (if possible)
- [ ] Build web app: `npm run build`
- [ ] Verify build succeeds
- [ ] Test production build locally
- [ ] Verify all features work in production

### 5.6 Documentation - User Facing

**File**: `README.md` (root)

- [ ] Add project title and description
- [ ] Add badges (build status, license, etc.)
- [ ] Add screenshot/GIF of app in action
- [ ] Add "Features" section
- [ ] Add "Quick Start" section:
  - [ ] Download proxy
  - [ ] Run proxy
  - [ ] Open web app
- [ ] Add "Installation" link to docs
- [ ] Add "Documentation" link to docs
- [ ] Add "Contributing" link
- [ ] Add "License" section (AGPL-3.0)
- [ ] Add credits/attribution

**File**: `docs/installation.md`

- [ ] Add system requirements
- [ ] Add download instructions
- [ ] Add platform-specific instructions:
  - [ ] Windows
  - [ ] macOS
  - [ ] Linux
- [ ] Add troubleshooting section
- [ ] Add FAQ

**File**: `docs/usage.md`

- [ ] Add connection instructions
- [ ] Add query examples
- [ ] Add autocomplete guide
- [ ] Add keyboard shortcuts
- [ ] Add tips & tricks

### 5.7 Documentation - Developer Facing

**File**: `proxy/README.md`

- [ ] Add "Development" section
- [ ] Add "Building" instructions
- [ ] Add "Testing" instructions
- [ ] Add "Contributing" guidelines

**File**: `web/README.md`

- [ ] Add "Development" section
- [ ] Add "Building" instructions
- [ ] Add "Testing" instructions
- [ ] Add "Contributing" guidelines

**File**: `docs/architecture.md`

- [ ] Add architecture overview
- [ ] Add component diagrams
- [ ] Add message protocol documentation
- [ ] Add security considerations
- [ ] Add deployment guide

**File**: `CONTRIBUTING.md`

- [ ] Add code of conduct
- [ ] Add contribution guidelines
- [ ] Add PR process
- [ ] Add coding standards
- [ ] Add testing requirements

### 5.8 GitHub Actions - CI/CD

**File**: `.github/workflows/test.yml`

- [ ] Add workflow for running tests on PR
- [ ] Run Go tests
- [ ] Run JS tests
- [ ] Check test coverage
- [ ] Run linters
- [ ] Fail on errors

**File**: `.github/workflows/release.yml`

- [ ] Add workflow for releases on tag push
- [ ] Build proxy for all platforms
- [ ] Calculate checksums
- [ ] Create GitHub release
- [ ] Upload binaries
- [ ] (Optional) Deploy web app to Cloudflare Pages

### 5.9 Security Audit

- [ ] Review secret generation (cryptographically secure)
- [ ] Review secret validation
- [ ] Review CORS configuration (localhost only)
- [ ] Review error messages (no sensitive info leaked)
- [ ] Review SQL injection protection (use parameterized queries)
- [ ] Review XSS protection
- [ ] Add security documentation
- [ ] Add responsible disclosure policy

### 5.10 Performance Testing

**Proxy**:
- [ ] Test with 1000 concurrent queries
- [ ] Measure memory usage
- [ ] Measure CPU usage
- [ ] Verify no memory leaks
- [ ] Profile with pprof if needed

**Web App**:
- [ ] Test loading time (should be < 2s)
- [ ] Test autocomplete latency (should be < 100ms)
- [ ] Test results rendering with 1000 rows
- [ ] Check bundle size
- [ ] Optimize if needed (code splitting, lazy loading)

### 5.11 Accessibility Audit

- [ ] Check keyboard navigation
- [ ] Check screen reader compatibility
- [ ] Check color contrast
- [ ] Add ARIA labels where needed
- [ ] Test with accessibility tools

### 5.12 Final Polish

**Web App**:
- [ ] Add favicon
- [ ] Add meta tags (title, description)
- [ ] Add Open Graph tags
- [ ] Improve error messages
- [ ] Add loading states everywhere
- [ ] Add empty states
- [ ] Polish animations/transitions
- [ ] Test all edge cases

**Proxy**:
- [ ] Improve CLI help text
- [ ] Improve error messages
- [ ] Add version command
- [ ] Polish output formatting
- [ ] Test all edge cases

### 5.13 Pre-Launch Checklist

- [ ] All tests passing (Go + JS)
- [ ] 80%+ code coverage (Go + JS)
- [ ] Zero linting errors
- [ ] All documentation complete
- [ ] All binaries built and tested
- [ ] Web app deployed
- [ ] README has clear instructions
- [ ] LICENSE file present (AGPL-3.0)
- [ ] Security review complete
- [ ] Performance acceptable
- [ ] Accessibility acceptable
- [ ] Cross-browser tested
- [ ] Mobile tested

**✅ Phase 5 Complete - Production ready!**

---

## 📋 Phase 6: Deployment & Release (Day 13)

### 6.1 Web App Deployment to Cloudflare Pages

- [ ] Create Cloudflare account (if needed)
- [ ] Install Wrangler CLI: `npm install -g wrangler`
- [ ] Login: `wrangler login`
- [ ] Create new Pages project
- [ ] Configure build settings:
  - [ ] Build command: `npm run build`
  - [ ] Build output: `build/`
- [ ] Deploy: `wrangler pages deploy build/`
- [ ] Test deployed app
- [ ] Verify all features work
- [ ] Configure custom domain (optional)
- [ ] Set up automatic deployments from GitHub

### 6.2 Proxy Binary Release

- [ ] Build all platform binaries: `cd proxy && make build-all`
- [ ] Calculate checksums: `cd bin && sha256sum * > checksums.txt`
- [ ] Create Git tag: `git tag v0.1.0`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Create GitHub release:
  - [ ] Title: "v0.1.0 - MVP Release"
  - [ ] Write release notes
  - [ ] Upload all binaries
  - [ ] Upload checksums.txt
- [ ] Verify download links work
- [ ] Test downloaded binaries

### 6.3 Release Notes

**File**: Write release notes with:

- [ ] Overview of features
- [ ] Supported platforms
- [ ] Installation instructions
- [ ] Quick start guide
- [ ] Known issues
- [ ] Roadmap for v0.2.0
- [ ] Contributors
- [ ] Changelog

### 6.4 Update Landing Page with Links

- [ ] Update download links to point to GitHub releases
- [ ] Update documentation links
- [ ] Add "View on GitHub" link
- [ ] Deploy updated web app

### 6.5 Launch Announcement

- [ ] Post on Hacker News
- [ ] Post on Reddit (r/PostgreSQL, r/webdev)
- [ ] Post on Twitter/X
- [ ] Post on LinkedIn
- [ ] Post in relevant Discord/Slack communities
- [ ] Write blog post (optional)
- [ ] Submit to ProductHunt (optional)

### 6.6 Monitoring & Feedback

- [ ] Set up GitHub issue templates
- [ ] Monitor GitHub issues
- [ ] Monitor social media feedback
- [ ] Respond to questions
- [ ] Fix critical bugs immediately
- [ ] Document common issues

**✅ Phase 6 Complete - MVP Launched! 🚀**

---

## 📋 Post-MVP: Future Enhancements

### Planned for v0.2.0

- [ ] Query history persistence (LocalStorage)
- [ ] Export results to CSV
- [ ] Export results to JSON
- [ ] Multiple query tabs
- [ ] Dark/light theme toggle
- [ ] Query templates/snippets
- [ ] Result pagination for large datasets
- [ ] Keyboard shortcuts panel
- [ ] Schema sidebar with tree view

### Planned for v0.3.0

- [ ] Query visualization (charts/graphs)
- [ ] EXPLAIN plan visualization
- [ ] Query formatting
- [ ] SQL linting
- [ ] Saved queries
- [ ] Connection profiles
- [ ] Auto-update mechanism for proxy

### Community Requested Features

- [ ] (Track feature requests from GitHub issues)
- [ ] (Prioritize based on user feedback)
- [ ] (Plan implementation in future versions)

---

## 🎯 Success Metrics

### MVP Success Criteria

- [x] ✅ All checkboxes above completed
- [ ] 100+ downloads in first week
- [ ] < 5 critical bugs reported
- [ ] Positive feedback on launch platforms
- [ ] 10+ GitHub stars
- [ ] 80%+ test coverage maintained

### Month 1 Goals

- [ ] 1,000+ downloads
- [ ] 50+ GitHub stars
- [ ] Active community in issues/discussions
- [ ] First external contribution (PR)
- [ ] v0.2.0 planning complete

### Month 3 Goals

- [ ] 5,000+ downloads
- [ ] 100+ GitHub stars
- [ ] Multiple active contributors
- [ ] Featured on awesome lists
- [ ] v0.2.0 and v0.3.0 released

---

## 📝 Notes & Issues

*Use this section to track blockers, decisions, and important notes during implementation*

### Blockers

- None yet

### Decisions

- Architecture: Go proxy + Static web app (decided)
- License: AGPL-3.0 (decided)
- Deployment: Cloudflare Pages (decided)

### Open Questions

- None

---

**Last Updated**: 2025-11-11
**Current Phase**: Pre-Phase 1 (Planning Complete)
**Next Action**: Begin Phase 1.1 - Project Structure Setup
