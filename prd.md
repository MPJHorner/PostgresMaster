# Product Requirements Document: Browser-Only Postgres Client (MVP)

## Overview
A browser-based SQL client for PostgreSQL that runs entirely in the browser with no backend infrastructure required. Built with modern web technologies focusing on developer experience and code quality.

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

### Framework & UI
- **Frontend Framework**: SvelteKit
- **UI Components**: shadcn-svelte (Svelte port of shadcn/ui)
- **Styling**: TailwindCSS (used by shadcn-svelte)

### SQL Editor
- **Editor**: Monaco Editor (VS Code's editor)
- **SQL Support**: @monaco-editor/react or svelte-monaco equivalent
- **Autocomplete**: SQL language service with Postgres-specific extensions
- **Schema Introspection**: Query `pg_catalog` and `information_schema` for metadata

### Database Connection
**[NEEDS CLARIFICATION]** - See Open Questions below

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
- Connection form with:
  - Host/Port
  - Database name
  - Username/Password
  - SSL options (if applicable)
- Test connection capability
- **[NEEDS CLARIFICATION]** - Connection persistence strategy

## User Flow

1. **Initial Load**
   - User opens application in browser
   - Presented with connection form

2. **Connect to Database**
   - User enters connection details
   - App validates and establishes connection
   - Schema introspection runs automatically

3. **Write Query**
   - User types in SQL editor
   - Autocomplete suggestions appear
   - User completes query

4. **Execute Query**
   - User clicks "Run" or presses Ctrl+Enter
   - Query executes
   - Results display in table below editor

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

## Open Questions

### Critical Questions (Blockers)

1. **Database Connection Method**
   - Option A: Use PGlite (Postgres compiled to WASM) - runs entirely in browser but limited to local database
   - Option B: Direct WebSocket connection to Postgres (requires pg_websocket extension or proxy)
   - Option C: Connect via a lightweight proxy service (conflicts with "no backend" requirement)
   - **Question**: Which connection method should we use? Or is the intent to work with a local WASM database?

2. **Connection Persistence**
   - Should connection credentials be saved? (LocalStorage, IndexedDB)
   - How to handle sensitive credentials in browser storage?
   - Should we support multiple saved connections?

3. **Query History**
   - Should we save query history?
   - Local storage or in-memory only?

### Nice-to-Have Clarifications

4. **Multiple Query Editors**
   - Single query editor or tabbed interface?
   - For MVP, assuming single editor is sufficient?

5. **Result Set Size Limits**
   - Should we paginate large result sets?
   - Maximum rows to display?

6. **Export Functionality**
   - Export results to CSV/JSON?
   - Consider for v2?

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
