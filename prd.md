# Product Requirements Document: Browser-Only Postgres Client (MVP)

## Overview
A browser extension SQL client for PostgreSQL that connects directly to remote Postgres databases with no backend infrastructure required. Built as a Chrome/Firefox extension using modern web technologies, focusing on developer experience and code quality.

**Key Innovation:** Runs entirely in the browser as an extension, leveraging extension APIs to make direct connections to remote Postgres servers without any proxy or backend service.

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

### Platform
- **Type**: Browser Extension (Chrome/Firefox)
- **Manifest**: Manifest V3 (Chrome), compatible with Firefox WebExtensions
- **Build System**: Vite with extension plugin

### Framework & UI
- **Frontend Framework**: Svelte (not SvelteKit - extensions don't need SSR)
- **UI Components**: shadcn-svelte (Svelte port of shadcn/ui)
- **Styling**: TailwindCSS (used by shadcn-svelte)

### SQL Editor
- **Editor**: Monaco Editor (VS Code's editor)
- **SQL Support**: Monaco SQL language support
- **Autocomplete**: Custom SQL language service with Postgres-specific keywords + schema-aware completions
- **Schema Introspection**: Query `pg_catalog` and `information_schema` for metadata

### Database Connection
- **Client Library**: `postgres.js` or `pg` (Node.js Postgres client)
- **Connection Method**: Direct TCP connection via extension background service worker
- **Authentication**: Standard Postgres authentication (password, md5, scram-sha-256)
- **SSL/TLS**: Support for SSL connections
- **Storage**: IndexedDB for connection credentials (encrypted), query history

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

### 4. Connection Management
- Connection form with fields:
  - Host (hostname/IP)
  - Port (default 5432)
  - Database name
  - Username
  - Password
  - SSL mode (disable, prefer, require)
- Save connections (encrypted in IndexedDB)
- Quick connect to saved connections
- Test connection button
- Connection status indicator

## User Flow

1. **Initial Setup**
   - User installs extension from Chrome Web Store / Firefox Add-ons
   - User clicks extension icon to open popup or side panel
   - Presented with connection form

2. **Connect to Database**
   - User enters connection details (host, port, database, username, password)
   - Optionally saves connection for future use
   - Clicks "Connect"
   - Extension background service establishes TCP connection to Postgres
   - Schema introspection runs automatically
   - Success: Editor becomes active

3. **Write Query**
   - User types in SQL editor
   - Autocomplete suggests:
     - SQL keywords (SELECT, FROM, WHERE, JOIN, etc.)
     - Table names from connected database
     - Column names for referenced tables
     - Postgres-specific functions and syntax
   - User completes query

4. **Execute Query**
   - User clicks "Run" or presses Ctrl+Enter
   - Query sent to remote Postgres via background service
   - Results streamed back to extension
   - Results display in table below editor
   - Execution time and row count shown

5. **Subsequent Sessions**
   - User reopens extension
   - Saved connections available for quick connect
   - Previous session state restored (last query, if applicable)

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
src/
├── lib/
│   ├── components/
│   │   ├── ui/          # shadcn-svelte components
│   │   ├── Editor.svelte
│   │   ├── Results.svelte
│   │   └── Connection.svelte
│   ├── services/
│   │   ├── database.ts  # Database connection logic
│   │   └── introspection.ts
│   ├── stores/
│   │   └── connection.ts
│   └── utils/
├── routes/
│   └── +page.svelte
└── tests/
```

## Implementation Decisions

### Confirmed for MVP
- ✅ PGlite (Postgres WASM) - runs entirely in browser
- ✅ Single query editor (no tabs)
- ✅ Auto-persist database to IndexedDB
- ✅ Display all query results (no pagination for MVP)
- ✅ Query history: in-memory only (cleared on refresh)

### Deferred to Post-MVP
- Import/export database files
- Export results to CSV/JSON
- Query history persistence
- Multiple query tabs
- Result pagination for large datasets

## Success Metrics (MVP)
- User can connect to a Postgres database from browser
- User can write a query with autocomplete assistance
- User can execute query and see results
- 80%+ test coverage
- Zero critical linting errors
- All tests passing

## Timeline (Estimated)
- Phase 1: Project setup, tooling, basic UI (shadcn-svelte) - 2 days
- Phase 2: Database connection + introspection - 3 days
- Phase 3: SQL editor with Monaco + autocomplete - 3 days
- Phase 4: Results display - 1 day
- Phase 5: Testing + polish - 2 days

**Total: ~11 days for MVP**

## Future Enhancements (Post-MVP)
- Query history and saved queries
- Result filtering and sorting
- Export to CSV/JSON
- Multiple query tabs
- Dark/light theme toggle
- Query snippets/templates
- Keyboard shortcuts panel
