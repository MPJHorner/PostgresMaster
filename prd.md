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
