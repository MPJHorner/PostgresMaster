# Product Requirements Document: Browser-Only Postgres Client MVP (Detailed)

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [System Architecture](#system-architecture)
3. [Component Specifications](#component-specifications)
4. [Development Phases](#development-phases)
5. [Testing Strategy](#testing-strategy)
6. [Success Criteria](#success-criteria)
7. [Technical Specifications](#technical-specifications)

---

## Executive Summary

### Vision
A modern, browser-based PostgreSQL client that requires zero backend infrastructure. Users download a tiny (~5MB) Go binary that runs locally, enabling secure connections to any PostgreSQL server through a clean web interface.

### Problem Statement
Existing database clients either:
- Require installation (pgAdmin, DBeaver, DataGrip)
- Need complex setup (backend services, Docker containers)
- Have limited functionality (online-only tools)
- Expose credentials to third-party services

### Solution
**Two-component architecture:**
1. **Web App** (100% frontend): Beautiful Svelte-based SQL editor with intelligent autocomplete
2. **Local Proxy** (Go binary): Lightweight bridge between browser WebSocket and Postgres TCP

### Key Benefits
- ✅ **Zero backend infrastructure** - Just static file hosting
- ✅ **Works with any Postgres** - No service lock-in
- ✅ **Secure** - Credentials never leave user's machine
- ✅ **Cross-platform** - Windows, macOS (Intel/ARM), Linux
- ✅ **Fast** - Go binary < 10MB, web app < 2s load time
- ✅ **Modern UX** - VS Code-like editor with shadcn-svelte UI

---

## System Architecture

### High-Level Architecture

```mermaid
graph TB
    User[User] -->|1. Downloads| Proxy[Go Proxy Binary]
    User -->|2. Visits| WebApp[Web App<br/>yourapp.com]

    Proxy -->|Connects to| PG[PostgreSQL<br/>Remote Server]
    Proxy -->|Generates| Secret[Secret Token]
    Proxy -->|Outputs URL| URL[http://localhost:8080?secret=xxx]

    User -->|3. Opens URL| Browser[Browser]
    Browser -->|WebSocket<br/>+ Secret Auth| Proxy

    Browser -->|4. Sends SQL| Proxy
    Proxy -->|Executes| PG
    PG -->|Returns Results| Proxy
    Proxy -->|Streams Back| Browser

    style WebApp fill:#a8d5ba
    style Proxy fill:#ffd6a5
    style PG fill:#c9c9ff
```

### Component Interaction Flow

```mermaid
sequenceDiagram
    participant U as User
    participant CLI as Proxy CLI
    participant WS as WebSocket Server
    participant PG as Postgres
    participant B as Browser
    participant WA as Web App

    U->>CLI: Run: ./postgres-proxy "postgres://..."
    CLI->>PG: Connect (with retry)
    PG-->>CLI: Connection OK
    CLI->>CLI: Generate Secret
    CLI->>WS: Start WebSocket Server :8080
    CLI->>U: Output: http://localhost:8080?secret=abc123

    U->>B: Open URL
    B->>WA: Load Web App
    WA->>WA: Parse ?secret param
    WA->>WS: WebSocket Connect (secret auth)
    WS->>WS: Validate Secret
    WS-->>WA: Connected

    WA->>WS: {"type": "introspect"}
    WS->>PG: Query pg_catalog
    PG-->>WS: Schema Data
    WS-->>WA: {"type": "schema", "payload": {...}}
    WA->>WA: Setup Autocomplete

    U->>WA: Type SQL Query
    WA->>WA: Show Autocomplete
    U->>WA: Press Ctrl+Enter

    WA->>WS: {"type": "query", "payload": {"sql": "..."}}
    WS->>PG: Execute Query
    PG-->>WS: Result Rows
    WS-->>WA: {"type": "result", "payload": {...}}
    WA->>WA: Render Results Table
```

### Data Flow Architecture

```mermaid
flowchart LR
    subgraph Browser
        UI[Web UI]
        WS_Client[WebSocket Client]
        Monaco[Monaco Editor]
        Autocomplete[Autocomplete Engine]
    end

    subgraph Go_Proxy[Go Proxy]
        WS_Server[WebSocket Server]
        Auth[Auth Validator]
        Router[Message Router]
        PG_Client[Postgres Client]
    end

    subgraph Remote
        Postgres[(PostgreSQL)]
    end

    UI -->|User Input| Monaco
    Monaco -->|SQL Text| Autocomplete
    Autocomplete -->|Suggestions| Monaco

    UI -->|Execute Query| WS_Client
    WS_Client -->|JSON Message| WS_Server
    WS_Server -->|Validate| Auth
    Auth -->|Route| Router
    Router -->|Execute| PG_Client
    PG_Client -->|TCP/SSL| Postgres

    Postgres -->|Result Rows| PG_Client
    PG_Client -->|Format| Router
    Router -->|JSON Response| WS_Server
    WS_Server -->|Stream| WS_Client
    WS_Client -->|Display| UI
```

---

## Component Specifications

### 1. Go Proxy Binary

#### 1.1 Responsibilities
- Parse Postgres connection string or prompt interactively
- Establish and maintain connection to remote Postgres
- Generate cryptographically secure random secret
- Start WebSocket server on localhost:8080
- Authenticate browser connections using secret
- Route messages between browser and Postgres
- Handle query execution and result streaming
- Graceful error handling and shutdown

#### 1.2 Directory Structure
```
proxy/
├── cmd/
│   └── proxy/
│       └── main.go                 # Entry point, CLI
├── pkg/
│   ├── server/
│   │   ├── websocket.go            # WebSocket server
│   │   └── handler.go              # Message handlers
│   ├── postgres/
│   │   ├── client.go               # Postgres connection
│   │   ├── query.go                # Query execution
│   │   └── introspect.go           # Schema introspection
│   ├── protocol/
│   │   ├── messages.go             # Message types
│   │   └── serializer.go           # JSON (de)serialization
│   └── auth/
│       └── secret.go               # Secret generation/validation
├── go.mod
├── go.sum
├── Makefile                        # Cross-compilation
└── README.md
```

#### 1.3 CLI Interface

**Syntax:**
```bash
postgres-proxy [CONNECTION_STRING]
```

**Option A: Connection String**
```bash
./postgres-proxy "postgresql://username:password@host:5432/database?sslmode=require"
```

**Option B: Interactive Mode**
```bash
./postgres-proxy

PostgreSQL Connection Setup
---------------------------
Host: db.example.com
Port [5432]: 5432
Database: myapp_production
Username: admin
Password: ********
SSL Mode [prefer]: require

Connecting to postgres://db.example.com:5432/myapp_production...
✓ Connected successfully!

WebSocket server running on ws://localhost:8080
→ Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6g7h8i9j0

Press Ctrl+C to stop
```

**Connection Retry Logic:**
```
Attempt 1: Immediate
Attempt 2: Wait 2s
Attempt 3: Wait 4s
Attempt 4: Wait 8s
Fail: Exit with error
```

#### 1.4 Message Protocol

**Message Structure:**
```go
type Message struct {
    ID      string          `json:"id"`      // UUID for request/response matching
    Type    string          `json:"type"`    // Message type
    Payload json.RawMessage `json:"payload"` // Type-specific payload
}
```

**Client → Server Messages:**

1. **Query Request**
```json
{
  "id": "uuid-1234",
  "type": "query",
  "payload": {
    "sql": "SELECT * FROM users WHERE id = $1",
    "params": [123],
    "timeout": 30000
  }
}
```

2. **Introspection Request**
```json
{
  "id": "uuid-5678",
  "type": "introspect",
  "payload": {}
}
```

3. **Ping (Keep-alive)**
```json
{
  "id": "uuid-9999",
  "type": "ping",
  "payload": {}
}
```

**Server → Client Messages:**

1. **Query Result**
```json
{
  "id": "uuid-1234",
  "type": "result",
  "payload": {
    "columns": [
      {"name": "id", "dataType": "integer", "nullable": false},
      {"name": "email", "dataType": "text", "nullable": false},
      {"name": "created_at", "dataType": "timestamp", "nullable": true}
    ],
    "rows": [
      {"id": 123, "email": "user@example.com", "created_at": "2024-01-15T10:30:00Z"}
    ],
    "rowCount": 1,
    "executionTime": 45,
    "affectedRows": 0
  }
}
```

2. **Schema Result**
```json
{
  "id": "uuid-5678",
  "type": "schema",
  "payload": {
    "tables": [
      {
        "schema": "public",
        "name": "users",
        "type": "table",
        "columns": [
          {"name": "id", "dataType": "integer", "nullable": false, "isPrimaryKey": true},
          {"name": "email", "dataType": "text", "nullable": false}
        ]
      }
    ],
    "functions": [
      {"schema": "public", "name": "get_user_count", "returnType": "integer"}
    ]
  }
}
```

3. **Error**
```json
{
  "id": "uuid-1234",
  "type": "error",
  "payload": {
    "code": "42P01",
    "message": "relation \"userz\" does not exist",
    "detail": "Hint: Did you mean \"users\"?",
    "position": 15
  }
}
```

4. **Pong**
```json
{
  "id": "uuid-9999",
  "type": "pong",
  "payload": {}
}
```

#### 1.5 Authentication

**Secret Generation:**
```go
func generateSecret() (string, error) {
    bytes := make([]byte, 32) // 32 bytes = 256 bits
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil // Returns 64-char hex string
}
```

**Secret Validation:**
- On WebSocket upgrade, check query parameter: `?secret=xxx`
- On every message, validate secret (stored in connection context)
- Invalid secret → Close connection with code 1008 (Policy Violation)

#### 1.6 Postgres Interaction

**Connection Pool:**
```go
// Use pgx connection pool
config, _ := pgxpool.ParseConfig(connString)
config.MaxConns = 5
config.MinConns = 1
pool, _ := pgxpool.NewWithConfig(context.Background(), config)
```

**Query Execution:**
```go
func (c *Client) ExecuteQuery(ctx context.Context, sql string, params []interface{}) (*QueryResult, error) {
    start := time.Now()

    rows, err := c.pool.Query(ctx, sql, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Read column descriptions
    fieldDescriptions := rows.FieldDescriptions()
    columns := make([]Column, len(fieldDescriptions))
    for i, fd := range fieldDescriptions {
        columns[i] = Column{
            Name:     string(fd.Name),
            DataType: c.getDataTypeName(fd.DataTypeOID),
        }
    }

    // Read rows
    results := []map[string]interface{}{}
    for rows.Next() {
        values, _ := rows.Values()
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col.Name] = values[i]
        }
        results = append(results, row)
    }

    executionTime := time.Since(start).Milliseconds()

    return &QueryResult{
        Columns:       columns,
        Rows:          results,
        RowCount:      len(results),
        ExecutionTime: executionTime,
    }, nil
}
```

**Schema Introspection:**
```sql
-- Get all tables
SELECT
    n.nspname AS schema,
    c.relname AS name,
    CASE c.relkind
        WHEN 'r' THEN 'table'
        WHEN 'v' THEN 'view'
        WHEN 'm' THEN 'materialized view'
    END AS type
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('r', 'v', 'm')
  AND n.nspname NOT IN ('pg_catalog', 'information_schema')
ORDER BY schema, name;

-- Get columns for a table
SELECT
    a.attname AS name,
    format_type(a.atttypid, a.atttypmod) AS data_type,
    a.attnotnull AS not_null,
    (SELECT COUNT(*) FROM pg_index i
     WHERE i.indrelid = a.attrelid
     AND a.attnum = ANY(i.indkey)
     AND i.indisprimary) > 0 AS is_primary_key
FROM pg_attribute a
WHERE a.attrelid = $1::regclass
  AND a.attnum > 0
  AND NOT a.attisdropped
ORDER BY a.attnum;

-- Get functions
SELECT
    n.nspname AS schema,
    p.proname AS name,
    pg_get_function_result(p.oid) AS return_type
FROM pg_proc p
JOIN pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
ORDER BY schema, name;
```

#### 1.7 Build System

**Makefile:**
```makefile
BINARY_NAME=postgres-proxy
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION} -s -w"

.PHONY: all clean test lint build-all

all: build-all

build-all: \
	build-windows-amd64 \
	build-darwin-amd64 \
	build-darwin-arm64 \
	build-linux-amd64 \
	build-linux-arm64

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-windows-amd64.exe cmd/proxy/main.go

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-amd64 cmd/proxy/main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-darwin-arm64 cmd/proxy/main.go

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-amd64 cmd/proxy/main.go

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}-linux-arm64 cmd/proxy/main.go

test:
	go test ./... -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run --timeout=5m

clean:
	rm -rf bin/ coverage.out coverage.html

.PHONY: dev
dev:
	go run cmd/proxy/main.go
```

**GitHub Actions (for releases):**
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build binaries
        run: cd proxy && make build-all

      - name: Calculate checksums
        run: cd proxy/bin && sha256sum * > checksums.txt

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: proxy/bin/*
```

#### 1.8 Testing

**Unit Tests:**
```go
// pkg/auth/secret_test.go
func TestGenerateSecret(t *testing.T) {
    secret, err := generateSecret()
    assert.NoError(t, err)
    assert.Len(t, secret, 64) // 32 bytes hex = 64 chars
}

// pkg/protocol/messages_test.go
func TestMessageSerialization(t *testing.T) {
    msg := Message{
        ID:   "test-123",
        Type: "query",
        Payload: json.RawMessage(`{"sql":"SELECT 1"}`),
    }

    data, err := json.Marshal(msg)
    assert.NoError(t, err)

    var decoded Message
    err = json.Unmarshal(data, &decoded)
    assert.NoError(t, err)
    assert.Equal(t, msg.ID, decoded.ID)
}
```

**Integration Tests:**
```go
// pkg/postgres/client_test.go (requires Docker)
func TestQueryExecution(t *testing.T) {
    // Start test Postgres container
    pool, resource := startTestPostgres(t)
    defer pool.Purge(resource)

    // Connect
    client, err := NewClient(testConnString)
    require.NoError(t, err)
    defer client.Close()

    // Execute query
    result, err := client.ExecuteQuery(context.Background(), "SELECT 1 as num", nil)
    require.NoError(t, err)

    assert.Equal(t, 1, result.RowCount)
    assert.Equal(t, "num", result.Columns[0].Name)
    assert.Equal(t, 1, result.Rows[0]["num"])
}
```

---

### 2. Web Application

#### 2.1 Responsibilities
- Provide landing page with proxy download links
- Parse `?secret` parameter from URL
- Connect to local proxy via WebSocket
- Display SQL editor with syntax highlighting
- Provide intelligent autocomplete (keywords + schema)
- Execute queries and display results
- Handle errors gracefully
- Maintain query history (in-memory)

#### 2.2 Directory Structure
```
web/
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/                     # shadcn-svelte components
│   │   │   │   ├── button.svelte
│   │   │   │   ├── input.svelte
│   │   │   │   ├── table.svelte
│   │   │   │   ├── alert.svelte
│   │   │   │   ├── card.svelte
│   │   │   │   └── ...
│   │   │   ├── LandingPage.svelte      # No secret
│   │   │   ├── ConnectionStatus.svelte # Connection indicator
│   │   │   ├── Editor.svelte           # Monaco wrapper
│   │   │   ├── QueryPanel.svelte       # Editor + Execute button
│   │   │   ├── Results.svelte          # Results table
│   │   │   ├── ErrorDisplay.svelte     # Error messages
│   │   │   └── SchemaSidebar.svelte    # Schema browser (optional)
│   │   ├── services/
│   │   │   ├── websocket.ts            # WebSocket client
│   │   │   ├── protocol.ts             # Message types (matches Go)
│   │   │   └── introspection.ts        # Schema parsing
│   │   ├── stores/
│   │   │   ├── connection.ts           # Connection state
│   │   │   ├── schema.ts               # Schema data
│   │   │   └── history.ts              # Query history
│   │   └── utils/
│   │       ├── autocomplete.ts         # Autocomplete logic
│   │       └── format.ts               # Value formatting
│   ├── routes/
│   │   ├── +layout.svelte              # Root layout
│   │   └── +page.svelte                # Main page
│   └── app.html                        # HTML template
├── static/
│   └── favicon.png
├── tests/
│   ├── services/
│   │   └── websocket.test.ts
│   └── components/
│       └── Editor.test.ts
├── package.json
├── svelte.config.js
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

#### 2.3 WebSocket Client Service

**File: `src/lib/services/websocket.ts`**

```typescript
import { v4 as uuidv4 } from 'uuid';
import type { ClientMessage, ServerMessage, QueryResult, SchemaInfo } from './protocol';

export class PostgresProxyClient {
  private ws: WebSocket | null = null;
  private secret: string;
  private pendingRequests: Map<string, {
    resolve: (value: any) => void;
    reject: (error: Error) => void;
  }> = new Map();

  constructor(secret: string) {
    this.secret = secret;
  }

  async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(`ws://localhost:8080?secret=${this.secret}`);

        this.ws.onopen = () => {
          console.log('Connected to proxy');
          resolve();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(new Error('Failed to connect to proxy'));
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data);
        };

        this.ws.onclose = () => {
          console.log('Disconnected from proxy');
          // Reject all pending requests
          this.pendingRequests.forEach(({ reject }) => {
            reject(new Error('Connection closed'));
          });
          this.pendingRequests.clear();
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  private handleMessage(data: string) {
    try {
      const message: ServerMessage = JSON.parse(data);
      const pending = this.pendingRequests.get(message.id);

      if (!pending) {
        console.warn('Received message for unknown request:', message.id);
        return;
      }

      if (message.type === 'error') {
        pending.reject(new Error(message.payload.message));
      } else {
        pending.resolve(message.payload);
      }

      this.pendingRequests.delete(message.id);
    } catch (error) {
      console.error('Failed to parse message:', error);
    }
  }

  private send<T>(type: string, payload: any): Promise<T> {
    return new Promise((resolve, reject) => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        reject(new Error('Not connected'));
        return;
      }

      const id = uuidv4();
      const message: ClientMessage = { id, type, payload };

      this.pendingRequests.set(id, { resolve, reject });
      this.ws.send(JSON.stringify(message));

      // Timeout after 60 seconds
      setTimeout(() => {
        if (this.pendingRequests.has(id)) {
          this.pendingRequests.delete(id);
          reject(new Error('Request timeout'));
        }
      }, 60000);
    });
  }

  async executeQuery(sql: string, params: any[] = []): Promise<QueryResult> {
    return this.send<QueryResult>('query', { sql, params });
  }

  async introspectSchema(): Promise<SchemaInfo> {
    return this.send<SchemaInfo>('introspect', {});
  }

  close() {
    this.ws?.close();
    this.ws = null;
  }
}
```

#### 2.4 Protocol Types

**File: `src/lib/services/protocol.ts`**

```typescript
// Matches Go message types

export interface ClientMessage {
  id: string;
  type: 'query' | 'introspect' | 'ping';
  payload: any;
}

export interface ServerMessage {
  id: string;
  type: 'result' | 'schema' | 'error' | 'pong';
  payload: any;
}

export interface QueryResult {
  columns: ColumnInfo[];
  rows: Record<string, any>[];
  rowCount: number;
  executionTime: number;
  affectedRows?: number;
}

export interface ColumnInfo {
  name: string;
  dataType: string;
  nullable: boolean;
  isPrimaryKey?: boolean;
}

export interface SchemaInfo {
  tables: TableInfo[];
  functions?: FunctionInfo[];
}

export interface TableInfo {
  schema: string;
  name: string;
  type: 'table' | 'view' | 'materialized view';
  columns: ColumnInfo[];
}

export interface FunctionInfo {
  schema: string;
  name: string;
  returnType: string;
}

export interface ErrorPayload {
  code?: string;
  message: string;
  detail?: string;
  position?: number;
}
```

#### 2.5 Connection Store

**File: `src/lib/stores/connection.ts`**

```typescript
import { writable, derived } from 'svelte/store';
import type { PostgresProxyClient } from '$lib/services/websocket';

interface ConnectionState {
  status: 'disconnected' | 'connecting' | 'connected' | 'error';
  client: PostgresProxyClient | null;
  error: string | null;
  secret: string | null;
}

export const connectionStore = writable<ConnectionState>({
  status: 'disconnected',
  client: null,
  error: null,
  secret: null
});

export const isConnected = derived(
  connectionStore,
  $conn => $conn.status === 'connected'
);

export const client = derived(
  connectionStore,
  $conn => $conn.client
);
```

#### 2.6 Schema Store

**File: `src/lib/stores/schema.ts`**

```typescript
import { writable, derived } from 'svelte/store';
import type { SchemaInfo, TableInfo } from '$lib/services/protocol';

export const schemaStore = writable<SchemaInfo | null>(null);

export const tables = derived(
  schemaStore,
  $schema => $schema?.tables || []
);

export const tableNames = derived(
  tables,
  $tables => $tables.map(t => t.name)
);

export const allColumns = derived(
  tables,
  $tables => {
    const columns: Array<{ table: string; column: string; type: string }> = [];
    $tables.forEach(table => {
      table.columns.forEach(col => {
        columns.push({
          table: table.name,
          column: col.name,
          type: col.dataType
        });
      });
    });
    return columns;
  }
);
```

#### 2.7 Main Page Component

**File: `src/routes/+page.svelte`**

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { connectionStore } from '$lib/stores/connection';
  import { schemaStore } from '$lib/stores/schema';
  import { PostgresProxyClient } from '$lib/services/websocket';

  import LandingPage from '$lib/components/LandingPage.svelte';
  import ConnectionStatus from '$lib/components/ConnectionStatus.svelte';
  import QueryPanel from '$lib/components/QueryPanel.svelte';
  import { Alert } from '$lib/components/ui/alert';

  let secret: string | null = null;

  onMount(async () => {
    // Parse secret from URL
    const params = new URLSearchParams(window.location.search);
    secret = params.get('secret');

    if (secret) {
      await connectToProxy(secret);
    }
  });

  async function connectToProxy(secretValue: string) {
    connectionStore.update(s => ({ ...s, status: 'connecting', secret: secretValue }));

    try {
      const client = new PostgresProxyClient(secretValue);
      await client.connect();

      connectionStore.update(s => ({ ...s, status: 'connected', client }));

      // Introspect schema
      const schema = await client.introspectSchema();
      schemaStore.set(schema);
    } catch (error) {
      connectionStore.update(s => ({
        ...s,
        status: 'error',
        error: error instanceof Error ? error.message : 'Connection failed'
      }));
    }
  }
</script>

<div class="min-h-screen bg-background">
  {#if !secret}
    <LandingPage />
  {:else}
    <div class="container mx-auto p-6">
      <ConnectionStatus />

      {#if $connectionStore.status === 'connecting'}
        <div class="flex items-center justify-center h-64">
          <p class="text-muted-foreground">Connecting to proxy...</p>
        </div>
      {:else if $connectionStore.status === 'error'}
        <Alert variant="destructive">
          {$connectionStore.error}
        </Alert>
      {:else if $connectionStore.status === 'connected'}
        <QueryPanel client={$connectionStore.client} />
      {/if}
    </div>
  {/if}
</div>
```

#### 2.8 Monaco Editor Component

**File: `src/lib/components/Editor.svelte`**

```svelte
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import * as monaco from 'monaco-editor';
  import { setupAutocomplete } from '$lib/utils/autocomplete';
  import { schemaStore } from '$lib/stores/schema';

  export let value = '';
  export let onChange: (value: string) => void = () => {};
  export let onExecute: () => void = () => {};

  let editorContainer: HTMLDivElement;
  let editor: monaco.editor.IStandaloneCodeEditor;
  let completionDisposable: monaco.IDisposable | null = null;

  onMount(() => {
    // Create editor
    editor = monaco.editor.create(editorContainer, {
      value: value,
      language: 'sql',
      theme: 'vs-dark',
      fontSize: 14,
      lineNumbers: 'on',
      minimap: { enabled: false },
      automaticLayout: true,
      scrollBeyondLastLine: false,
      wordWrap: 'on',
      formatOnPaste: true,
      formatOnType: true
    });

    // Setup autocomplete
    completionDisposable = setupAutocomplete(monaco, $schemaStore);

    // Listen for changes
    editor.onDidChangeModelContent(() => {
      onChange(editor.getValue());
    });

    // Ctrl+Enter to execute
    editor.addCommand(
      monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
      onExecute
    );

    // Update autocomplete when schema changes
    const unsubscribe = schemaStore.subscribe(schema => {
      if (completionDisposable) {
        completionDisposable.dispose();
      }
      completionDisposable = setupAutocomplete(monaco, schema);
    });

    return () => {
      unsubscribe();
      completionDisposable?.dispose();
      editor?.dispose();
    };
  });

  onDestroy(() => {
    editor?.dispose();
  });
</script>

<div bind:this={editorContainer} class="h-96 border rounded-md overflow-hidden" />
```

#### 2.9 Autocomplete Implementation

**File: `src/lib/utils/autocomplete.ts`**

```typescript
import type * as monaco from 'monaco-editor';
import type { SchemaInfo } from '$lib/services/protocol';

const POSTGRES_KEYWORDS = [
  'SELECT', 'FROM', 'WHERE', 'JOIN', 'INNER', 'LEFT', 'RIGHT', 'OUTER', 'CROSS',
  'ON', 'USING', 'GROUP BY', 'ORDER BY', 'HAVING', 'LIMIT', 'OFFSET',
  'INSERT', 'INTO', 'VALUES', 'UPDATE', 'SET', 'DELETE', 'TRUNCATE',
  'CREATE', 'TABLE', 'DROP', 'ALTER', 'INDEX', 'VIEW', 'SEQUENCE',
  'AS', 'DISTINCT', 'DISTINCT ON', 'ALL',
  'CASE', 'WHEN', 'THEN', 'ELSE', 'END',
  'AND', 'OR', 'NOT', 'IN', 'EXISTS', 'BETWEEN', 'LIKE', 'ILIKE', 'SIMILAR TO',
  'IS', 'NULL', 'TRUE', 'FALSE',
  'UNION', 'UNION ALL', 'INTERSECT', 'EXCEPT',
  'WITH', 'RECURSIVE', 'CTE',
  'RETURNING', 'CONFLICT', 'DO', 'NOTHING', 'ON CONFLICT',
  'SERIAL', 'BIGSERIAL', 'TEXT', 'VARCHAR', 'INTEGER', 'BIGINT', 'NUMERIC',
  'BOOLEAN', 'DATE', 'TIMESTAMP', 'TIMESTAMPTZ', 'INTERVAL',
  'JSONB', 'JSON', 'UUID', 'ARRAY', 'HSTORE',
  'PRIMARY KEY', 'FOREIGN KEY', 'REFERENCES', 'UNIQUE', 'CHECK', 'DEFAULT',
  'NOT NULL', 'CASCADE', 'RESTRICT', 'NO ACTION',
];

const POSTGRES_FUNCTIONS = [
  'COUNT', 'SUM', 'AVG', 'MIN', 'MAX',
  'NOW', 'CURRENT_TIMESTAMP', 'CURRENT_DATE', 'CURRENT_TIME',
  'EXTRACT', 'DATE_TRUNC', 'AGE',
  'COALESCE', 'NULLIF', 'GREATEST', 'LEAST',
  'CONCAT', 'CONCAT_WS', 'LOWER', 'UPPER', 'TRIM', 'LENGTH',
  'SUBSTRING', 'POSITION', 'REPLACE',
  'CAST', 'CONVERT',
  'ARRAY_AGG', 'STRING_AGG', 'JSON_AGG', 'JSONB_AGG',
  'ROW_NUMBER', 'RANK', 'DENSE_RANK', 'LAG', 'LEAD',
];

export function setupAutocomplete(
  monaco: typeof import('monaco-editor'),
  schema: SchemaInfo | null
): monaco.IDisposable {
  return monaco.languages.registerCompletionItemProvider('sql', {
    provideCompletionItems: (model, position) => {
      const word = model.getWordUntilPosition(position);
      const range = {
        startLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endLineNumber: position.lineNumber,
        endColumn: word.endColumn
      };

      const suggestions: monaco.languages.CompletionItem[] = [];

      // Add keywords
      POSTGRES_KEYWORDS.forEach(keyword => {
        suggestions.push({
          label: keyword,
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: keyword,
          range: range,
          sortText: '1_' + keyword // Sort keywords first
        });
      });

      // Add functions
      POSTGRES_FUNCTIONS.forEach(func => {
        suggestions.push({
          label: func,
          kind: monaco.languages.CompletionItemKind.Function,
          insertText: func + '($1)',
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range: range,
          sortText: '2_' + func
        });
      });

      // Add table names
      if (schema) {
        schema.tables.forEach(table => {
          suggestions.push({
            label: table.name,
            kind: monaco.languages.CompletionItemKind.Class,
            insertText: table.name,
            detail: `Table (${table.schema})`,
            documentation: table.columns.map(c => `${c.name}: ${c.dataType}`).join('\n'),
            range: range,
            sortText: '3_' + table.name
          });

          // Add columns for this table
          table.columns.forEach(col => {
            suggestions.push({
              label: `${table.name}.${col.name}`,
              kind: monaco.languages.CompletionItemKind.Field,
              insertText: col.name,
              detail: `Column: ${col.dataType}`,
              documentation: `From table: ${table.name}`,
              range: range,
              sortText: '4_' + col.name
            });
          });
        });

        // Add functions from schema
        schema.functions?.forEach(func => {
          suggestions.push({
            label: func.name,
            kind: monaco.languages.CompletionItemKind.Function,
            insertText: func.name + '($1)',
            insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            detail: `Returns: ${func.returnType}`,
            documentation: `Schema: ${func.schema}`,
            range: range,
            sortText: '2_' + func.name
          });
        });
      }

      return { suggestions };
    }
  });
}
```

#### 2.10 Query Panel Component

**File: `src/lib/components/QueryPanel.svelte`**

```svelte
<script lang="ts">
  import Editor from './Editor.svelte';
  import Results from './Results.svelte';
  import ErrorDisplay from './ErrorDisplay.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator';
  import type { PostgresProxyClient } from '$lib/services/websocket';
  import type { QueryResult } from '$lib/services/protocol';

  export let client: PostgresProxyClient;

  let sql = 'SELECT version();';
  let results: QueryResult | null = null;
  let error: string | null = null;
  let loading = false;

  async function executeQuery() {
    if (!sql.trim()) {
      error = 'Please enter a query';
      return;
    }

    loading = true;
    error = null;
    results = null;

    try {
      results = await client.executeQuery(sql);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Query execution failed';
    } finally {
      loading = false;
    }
  }

  function handleEditorChange(value: string) {
    sql = value;
  }
</script>

<div class="space-y-6">
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">SQL Editor</h2>
      <Button on:click={executeQuery} disabled={loading}>
        {loading ? 'Executing...' : 'Run Query'}
        <kbd class="ml-2 text-xs">Ctrl+Enter</kbd>
      </Button>
    </div>

    <Editor
      value={sql}
      onChange={handleEditorChange}
      onExecute={executeQuery}
    />
  </div>

  <Separator />

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="text-muted-foreground">Executing query...</div>
    </div>
  {:else if error}
    <ErrorDisplay {error} />
  {:else if results}
    <Results data={results} />
  {/if}
</div>
```

#### 2.11 Results Component

**File: `src/lib/components/Results.svelte`**

```svelte
<script lang="ts">
  import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow
  } from '$lib/components/ui/table';
  import { Badge } from '$lib/components/ui/badge';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
  import type { QueryResult } from '$lib/services/protocol';
  import { formatValue } from '$lib/utils/format';

  export let data: QueryResult;

  $: columns = data.columns;
  $: rows = data.rows;
  $: hasResults = rows.length > 0;
</script>

<Card>
  <CardHeader>
    <CardTitle>
      <div class="flex items-center justify-between">
        <span>Results</span>
        <div class="flex gap-4 text-sm font-normal text-muted-foreground">
          <span>{data.rowCount} {data.rowCount === 1 ? 'row' : 'rows'}</span>
          <span>{data.executionTime}ms</span>
        </div>
      </div>
    </CardTitle>
  </CardHeader>
  <CardContent>
    {#if hasResults}
      <div class="border rounded-md overflow-auto max-h-[500px]">
        <Table>
          <TableHeader class="sticky top-0 bg-background">
            <TableRow>
              {#each columns as column}
                <TableHead>
                  <div class="flex flex-col gap-1">
                    <span class="font-semibold">{column.name}</span>
                    <Badge variant="outline" class="text-xs w-fit">
                      {column.dataType}
                    </Badge>
                  </div>
                </TableHead>
              {/each}
            </TableRow>
          </TableHeader>
          <TableBody>
            {#each rows as row}
              <TableRow>
                {#each columns as column}
                  <TableCell class="font-mono text-sm">
                    {formatValue(row[column.name])}
                  </TableCell>
                {/each}
              </TableRow>
            {/each}
          </TableBody>
        </Table>
      </div>
    {:else}
      <div class="text-center py-8 text-muted-foreground">
        Query executed successfully. No rows returned.
      </div>
    {/if}
  </CardContent>
</Card>
```

#### 2.12 Value Formatting Utility

**File: `src/lib/utils/format.ts`**

```typescript
export function formatValue(value: any): string {
  if (value === null || value === undefined) {
    return 'NULL';
  }

  if (typeof value === 'boolean') {
    return value ? 'true' : 'false';
  }

  if (typeof value === 'number') {
    return value.toString();
  }

  if (typeof value === 'string') {
    return value;
  }

  if (value instanceof Date) {
    return value.toISOString();
  }

  if (Array.isArray(value)) {
    return JSON.stringify(value);
  }

  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2);
  }

  return String(value);
}
```

#### 2.13 Landing Page

**File: `src/lib/components/LandingPage.svelte`**

```svelte
<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Download, Database, Zap, Lock } from 'lucide-svelte';

  const downloads = [
    { os: 'Windows (x64)', file: 'postgres-proxy-windows-amd64.exe', icon: '🪟' },
    { os: 'macOS (Intel)', file: 'postgres-proxy-darwin-amd64', icon: '' },
    { os: 'macOS (Apple Silicon)', file: 'postgres-proxy-darwin-arm64', icon: '' },
    { os: 'Linux (x64)', file: 'postgres-proxy-linux-amd64', icon: '🐧' }
  ];

  const features = [
    { icon: Database, title: 'Any Postgres', description: 'Connect to any PostgreSQL server - cloud or self-hosted' },
    { icon: Zap, title: 'Lightning Fast', description: 'Lightweight proxy (~5MB) with instant query execution' },
    { icon: Lock, title: 'Secure', description: 'Credentials never leave your machine' }
  ];
</script>

<div class="min-h-screen bg-gradient-to-b from-background to-muted">
  <!-- Hero Section -->
  <div class="container mx-auto px-6 py-20">
    <div class="text-center space-y-6 max-w-3xl mx-auto">
      <h1 class="text-5xl font-bold tracking-tight">
        PostgreSQL Client<br />for the Browser
      </h1>
      <p class="text-xl text-muted-foreground">
        Modern SQL client with zero backend. Download a tiny proxy, connect to any Postgres, and start querying.
      </p>
    </div>

    <!-- Features -->
    <div class="grid md:grid-cols-3 gap-6 mt-16">
      {#each features as feature}
        <Card>
          <CardHeader>
            <svelte:component this={feature.icon} class="h-10 w-10 mb-2 text-primary" />
            <CardTitle>{feature.title}</CardTitle>
            <CardDescription>{feature.description}</CardDescription>
          </CardHeader>
        </Card>
      {/each}
    </div>

    <!-- Download Section -->
    <div class="mt-20">
      <Card class="max-w-3xl mx-auto">
        <CardHeader>
          <CardTitle class="text-center text-2xl">
            Get Started in 3 Steps
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-8">
          <!-- Step 1: Download -->
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <div class="flex items-center justify-center w-8 h-8 rounded-full bg-primary text-primary-foreground font-bold">
                1
              </div>
              <h3 class="text-lg font-semibold">Download Proxy for Your OS</h3>
            </div>
            <div class="grid sm:grid-cols-2 gap-3 ml-11">
              {#each downloads as dl}
                <Button variant="outline" class="justify-start" asChild>
                  <a href="/downloads/{dl.file}" download>
                    <Download class="mr-2 h-4 w-4" />
                    <span class="mr-2">{dl.icon}</span>
                    {dl.os}
                  </a>
                </Button>
              {/each}
            </div>
          </div>

          <!-- Step 2: Run -->
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <div class="flex items-center justify-center w-8 h-8 rounded-full bg-primary text-primary-foreground font-bold">
                2
              </div>
              <h3 class="text-lg font-semibold">Run Proxy with Your Connection</h3>
            </div>
            <div class="ml-11 space-y-2">
              <code class="block p-4 bg-muted rounded-md text-sm">
                ./postgres-proxy "postgres://user:pass@host:5432/db"
              </code>
              <p class="text-sm text-muted-foreground">
                Or run without arguments for interactive mode
              </p>
            </div>
          </div>

          <!-- Step 3: Connect -->
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <div class="flex items-center justify-center w-8 h-8 rounded-full bg-primary text-primary-foreground font-bold">
                3
              </div>
              <h3 class="text-lg font-semibold">Open URL in Browser</h3>
            </div>
            <div class="ml-11">
              <p class="text-sm text-muted-foreground">
                Proxy will output a URL with secret token. Open it and start querying!
              </p>
              <code class="block p-4 bg-muted rounded-md text-sm mt-2">
                → Open: http://localhost:8080?secret=abc123...
              </code>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- FAQ / Tech Stack -->
    <div class="mt-16 text-center text-sm text-muted-foreground">
      <p>Built with SvelteKit • Powered by Go • Open Source</p>
    </div>
  </div>
</div>
```

---

## Development Phases

### Phase 1: Go Proxy Foundation (Days 1-3)

**Goal**: Working proxy that connects to Postgres and starts WebSocket server

**Deliverables:**
1. ✅ CLI that parses connection string or prompts interactively
2. ✅ Postgres connection with retry logic
3. ✅ Secret generation
4. ✅ WebSocket server on localhost:8080
5. ✅ Basic message routing (echo test)
6. ✅ Unit tests (80% coverage)
7. ✅ Cross-platform build system

**Testing Checklist:**
- [ ] Run proxy with connection string - connects successfully
- [ ] Run proxy in interactive mode - prompts for details
- [ ] Invalid credentials - retries then fails gracefully
- [ ] Secret is 64 characters (32 bytes hex)
- [ ] WebSocket server accepts connections on :8080
- [ ] Can send/receive JSON messages
- [ ] Builds for all platforms (Windows, macOS, Linux)

**Acceptance Criteria:**
```bash
# Success output:
./postgres-proxy "postgres://user:pass@localhost:5432/db"

Connecting to postgres://localhost:5432/db...
✓ Connected successfully!

WebSocket server running on ws://localhost:8080
→ Open in browser: http://localhost:8080?secret=a1b2c3...f6g7

Press Ctrl+C to stop
```

---

### Phase 2: Query Execution (Days 4-5)

**Goal**: Proxy executes queries and returns results

**Deliverables:**
1. ✅ Query execution handler
2. ✅ Result serialization (rows + metadata)
3. ✅ Error handling and formatting
4. ✅ Schema introspection queries
5. ✅ Integration tests with real Postgres

**Testing Checklist:**
- [ ] Execute `SELECT 1` - returns correct result
- [ ] Execute `SELECT * FROM pg_tables` - returns multiple rows
- [ ] Execute invalid SQL - returns error message
- [ ] Introspect schema - returns tables and columns
- [ ] Query timeout works (30s default)
- [ ] Handle NULL values correctly
- [ ] Handle various data types (int, text, timestamp, jsonb, etc.)

**Test Script:**
```javascript
// Manual test via WebSocket client
const ws = new WebSocket('ws://localhost:8080?secret=abc123');

ws.onopen = () => {
  // Test 1: Simple query
  ws.send(JSON.stringify({
    id: '1',
    type: 'query',
    payload: { sql: 'SELECT 1 as num, \'hello\' as text' }
  }));

  // Test 2: Schema introspection
  ws.send(JSON.stringify({
    id: '2',
    type: 'introspect',
    payload: {}
  }));
};

ws.onmessage = (event) => {
  console.log('Received:', JSON.parse(event.data));
};
```

**Acceptance Criteria:**
- Query execution returns results in < 1s for simple queries
- Errors include SQL state codes (42P01, etc.)
- Schema introspection returns all user tables/columns
- No memory leaks after 100 queries

---

### Phase 3: Web App Foundation (Days 6-7)

**Goal**: Web app connects to proxy and displays connection status

**Deliverables:**
1. ✅ SvelteKit app with static adapter
2. ✅ shadcn-svelte integration
3. ✅ Landing page with download links
4. ✅ URL secret parsing
5. ✅ WebSocket client service
6. ✅ Connection flow (landing → connecting → connected)

**Testing Checklist:**
- [ ] Visit `/` without secret - shows landing page
- [ ] Visit `/?secret=invalid` - shows error
- [ ] Visit `/?secret=valid` (proxy running) - connects successfully
- [ ] Connection status displays correctly
- [ ] Page builds as static site (`npm run build`)
- [ ] Built site works when served (`npm run preview`)

**Test Commands:**
```bash
# Terminal 1: Start proxy
./postgres-proxy "postgres://localhost:5432/postgres"
# Note the secret from output

# Terminal 2: Start web app
cd web
npm run dev

# Browser:
# Visit http://localhost:5173
# Click download link (fake for now)
# Visit http://localhost:5173?secret=<secret-from-proxy>
# Should see "Connected" status
```

**Acceptance Criteria:**
- Landing page loads in < 1s
- Connection to proxy succeeds within 2s
- Error messages are clear and actionable
- UI is responsive on mobile and desktop

---

### Phase 4: SQL Editor & Autocomplete (Days 8-10)

**Goal**: Full-featured SQL editor with intelligent autocomplete

**Deliverables:**
1. ✅ Monaco editor integration
2. ✅ SQL syntax highlighting
3. ✅ Keyword autocomplete
4. ✅ Schema-aware autocomplete (tables, columns)
5. ✅ Ctrl+Enter to execute
6. ✅ Query execution flow

**Testing Checklist:**
- [ ] Editor loads and displays SQL
- [ ] Syntax highlighting works
- [ ] Type `SEL` - suggests `SELECT`
- [ ] Type `SELECT * FROM ` - suggests table names
- [ ] Type `SELECT users.` - suggests columns from users table
- [ ] Ctrl+Enter executes query
- [ ] Execute button works
- [ ] Loading state shows during execution

**Interactive Test:**
```sql
-- Test autocomplete at each step:

SELECT    -- Should suggest: *, columns, functions
FROM      -- Should suggest: table names
WHERE     -- Should suggest: columns, operators
```

**Acceptance Criteria:**
- Autocomplete appears within 100ms
- Suggestions are relevant to context
- Editor supports multi-line queries
- No lag when typing (smooth 60fps)

---

### Phase 5: Results Display & Error Handling (Days 11-12)

**Goal**: Beautiful results table and clear error messages

**Deliverables:**
1. ✅ Results table component
2. ✅ Column headers with data types
3. ✅ Row count and execution time
4. ✅ Value formatting (NULL, JSON, timestamps, etc.)
5. ✅ Error display component
6. ✅ Scrolling for large result sets

**Testing Checklist:**
- [ ] Query with 0 rows - shows "No results"
- [ ] Query with 1 row - displays correctly
- [ ] Query with 100 rows - scrolls smoothly
- [ ] NULL values display as "NULL"
- [ ] JSON values are formatted
- [ ] Timestamps are readable
- [ ] Execution time is accurate
- [ ] Error messages are clear

**Test Queries:**
```sql
-- Test various result types

-- Empty result
SELECT * FROM pg_tables WHERE tablename = 'nonexistent';

-- Single row
SELECT version();

-- Multiple rows with various types
SELECT
  1 as int_col,
  'text' as text_col,
  true as bool_col,
  NULL as null_col,
  NOW() as timestamp_col,
  '{"key": "value"}'::jsonb as json_col;

-- Large result
SELECT * FROM generate_series(1, 1000);

-- Error
SELECT * FROM nonexistent_table;
```

**Acceptance Criteria:**
- Results table renders in < 500ms for 100 rows
- Scrolling is smooth (60fps)
- All data types display correctly
- Errors show SQL position when available

---

### Phase 6: Polish & Documentation (Days 13-14)

**Goal**: Production-ready MVP with full documentation

**Deliverables:**
1. ✅ Code linting (ESLint, golangci-lint)
2. ✅ Code formatting (Prettier, gofmt)
3. ✅ Test coverage > 80%
4. ✅ README files
5. ✅ User documentation
6. ✅ Developer documentation
7. ✅ Deployment (web app to Vercel/Netlify)
8. ✅ Release (binaries to GitHub)

**Testing Checklist:**
- [ ] All Go tests pass: `go test ./...`
- [ ] All JS tests pass: `npm test`
- [ ] Go coverage > 80%: `go test -cover ./...`
- [ ] JS coverage > 80%: `npm run test:coverage`
- [ ] No linting errors: `golangci-lint run && npm run lint`
- [ ] Code is formatted: `gofmt && npm run format`
- [ ] Binaries work on all platforms
- [ ] Web app deploys successfully
- [ ] Documentation is clear and complete

**Documentation Checklist:**
- [ ] README.md (root) - Project overview, quick start
- [ ] proxy/README.md - Building, testing, contributing
- [ ] web/README.md - Development setup, testing
- [ ] INSTALLATION.md - User installation guide
- [ ] ARCHITECTURE.md - Technical deep-dive
- [ ] CONTRIBUTING.md - How to contribute
- [ ] LICENSE - Open source license

**Deployment Steps:**
```bash
# Build proxy binaries
cd proxy
make build-all

# Create GitHub release
gh release create v0.1.0 \
  --title "MVP Release" \
  --notes "Initial release" \
  bin/*

# Deploy web app
cd web
npm run build
vercel deploy --prod
```

**Acceptance Criteria:**
- Binaries downloadable from GitHub releases
- Web app accessible at public URL
- All documentation proofread and accurate
- No critical bugs
- Passes manual E2E test on all platforms

---

## Testing Strategy

### Unit Testing

**Go Proxy:**
```bash
# Run tests
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# With race detection
go test ./... -race
```

**Key test files:**
- `pkg/auth/secret_test.go` - Secret generation/validation
- `pkg/protocol/messages_test.go` - Message serialization
- `pkg/postgres/client_test.go` - Query execution (mocked)
- `pkg/server/websocket_test.go` - WebSocket handling (mocked)

**Web App:**
```bash
# Run tests
npm test

# With coverage
npm run test:coverage

# Watch mode
npm run test:watch
```

**Key test files:**
- `src/lib/services/websocket.test.ts` - WebSocket client
- `src/lib/utils/autocomplete.test.ts` - Autocomplete logic
- `src/lib/utils/format.test.ts` - Value formatting
- `src/lib/components/Editor.test.ts` - Editor component
- `src/lib/components/Results.test.ts` - Results rendering

### Integration Testing

**Proxy → Postgres:**
```go
// Use testcontainers-go for real Postgres
func TestIntegration(t *testing.T) {
    ctx := context.Background()

    // Start Postgres container
    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    defer postgresContainer.Terminate(ctx)

    connString, _ := postgresContainer.ConnectionString(ctx)

    // Connect and test
    client, err := NewClient(connString)
    require.NoError(t, err)

    result, err := client.ExecuteQuery(ctx, "SELECT 1 as num", nil)
    require.NoError(t, err)
    assert.Equal(t, 1, result.RowCount)
}
```

**Web App → Proxy:**
```typescript
// Use Playwright for E2E tests
import { test, expect } from '@playwright/test';

test('full query flow', async ({ page }) => {
  // Start mock proxy in test
  const mockProxy = startMockProxy();

  // Visit app with secret
  await page.goto('http://localhost:5173?secret=test123');

  // Wait for connection
  await expect(page.locator('text=Connected')).toBeVisible();

  // Type query
  await page.locator('.monaco-editor textarea').fill('SELECT 1');

  // Execute
  await page.keyboard.press('Control+Enter');

  // Check results
  await expect(page.locator('text=1 row')).toBeVisible();

  mockProxy.close();
});
```

### Manual Testing Checklist

**Before Each Release:**

1. **Proxy Binary Testing**
   - [ ] Windows x64: Run, connect to Postgres, generate URL
   - [ ] macOS Intel: Same
   - [ ] macOS ARM: Same
   - [ ] Linux x64: Same
   - [ ] Interactive mode works on all platforms
   - [ ] Connection string mode works on all platforms
   - [ ] Retry logic works (test with wrong password)

2. **Web App Testing**
   - [ ] Landing page loads on Chrome, Firefox, Safari
   - [ ] Download links work
   - [ ] Connection flow works
   - [ ] Query execution works
   - [ ] Autocomplete works
   - [ ] Results display correctly
   - [ ] Error handling works
   - [ ] Mobile responsive
   - [ ] Works offline after initial load

3. **Security Testing**
   - [ ] Invalid secret rejected
   - [ ] No secret = landing page only
   - [ ] Credentials not exposed in browser DevTools
   - [ ] XSS protection (test with malicious SQL)
   - [ ] CSRF protection (secret per session)

4. **Performance Testing**
   - [ ] Query with 1000 rows completes in < 2s
   - [ ] Autocomplete responds in < 100ms
   - [ ] No memory leaks after 100 queries
   - [ ] WebSocket reconnects after disconnect
   - [ ] Proxy uses < 50MB RAM under load

---

## Success Criteria

### MVP Launch Criteria

**Must Have:**
- ✅ Proxy binary < 10MB for all platforms
- ✅ Web app loads in < 2s
- ✅ Connects to any Postgres 12+
- ✅ Query execution with results display
- ✅ Autocomplete for SQL keywords
- ✅ Autocomplete for schema (tables, columns)
- ✅ Error messages are clear
- ✅ 80% test coverage (both codebases)
- ✅ Zero critical bugs
- ✅ Documentation complete

**Nice to Have (Stretch Goals):**
- ⭐ Query history persistence
- ⭐ Export results to CSV
- ⭐ Schema sidebar with tree view
- ⭐ Dark/light theme toggle
- ⭐ Multiple query tabs
- ⭐ Query templates/snippets

### Post-Launch Success Metrics

**Week 1:**
- 100+ downloads
- < 5 critical bug reports
- Positive feedback on HackerNews/Reddit

**Month 1:**
- 1,000+ downloads
- 10+ GitHub stars
- Feature requests in GitHub issues
- Community contributions (PRs)

**Month 3:**
- 5,000+ downloads
- 50+ GitHub stars
- v2 planning based on user feedback

---

## Technical Specifications

### Supported Postgres Versions
- PostgreSQL 12+
- All standard Postgres distributions (vanilla, AWS RDS, Google Cloud SQL, Azure, etc.)
- Compatible extensions: PostGIS, pgvector, etc.

### Browser Support
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

### Operating Systems (Proxy)
- Windows 10/11 (x64)
- macOS 10.15+ (Intel and Apple Silicon)
- Linux (Ubuntu 20.04+, Debian 10+, Fedora 33+, etc.)

### Network Requirements
- Proxy requires outbound connection to Postgres (port 5432 or custom)
- Web app requires connection to localhost:8080
- No inbound ports need to be opened
- Works behind corporate proxies (HTTP_PROXY env var supported)

### Security Considerations
- Secrets are cryptographically random (256 bits)
- Secrets are single-use per proxy session
- WebSocket connections are local only (localhost:8080)
- No credentials stored in browser
- No telemetry or analytics (100% private)

### Performance Targets
- Proxy startup: < 1s
- Query execution overhead: < 10ms
- Autocomplete latency: < 100ms
- Results rendering (100 rows): < 500ms
- Memory usage (proxy): < 50MB
- Memory usage (web app): < 100MB

### Data Limits
- Max query size: 10MB
- Max result set: 10,000 rows (for MVP, no pagination)
- Max WebSocket message: 50MB
- Query timeout: 30s (configurable)

---

## Open Questions

### 1. Microsoft SQL Autocomplete Package

You mentioned wanting to use "the Microsoft package that does autocomplete help on writing Postgres SQL". I need clarification:

**Options:**
- A) Monaco Editor's built-in SQL support (already included above)
- B) `sql-language-server` (LSP for SQL)
- C) Azure Data Studio's extension architecture
- D) Something else?

**Please specify** which package/library you're referring to, and I'll integrate it.

---

### 2. Linter/Formatter Configuration

**Preferences for:**
- Go: `golangci-lint` with default settings OK?
- JavaScript: Airbnb style guide or Standard.js?
- Prettier: Default config or custom?

---

### 3. Deployment

**Web App hosting:**
- Vercel (recommended - free, easy)
- Netlify
- GitHub Pages
- Cloudflare Pages
- Other?

**Binary distribution:**
- GitHub Releases (recommended - free, integrated)
- Custom download server
- CDN?

---

### 4. License

**Open source license:**
- MIT (most permissive)
- Apache 2.0
- GPL v3
- Other?

---

## Next Steps

1. **Answer open questions** above
2. **Review and approve** this PRD
3. **Create GitHub repository**
4. **Begin Phase 1**: Go proxy foundation

Once approved, I can:
- Generate initial project structure
- Set up build systems
- Create skeleton code for all components
- Set up testing frameworks

Ready to start building? 🚀
