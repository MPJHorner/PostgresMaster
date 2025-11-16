# Usage Guide

This guide covers everything you need to know to use the PostgreSQL Client effectively.

## Table of Contents

- [Connecting to a Database](#connecting-to-a-database)
- [Writing Queries](#writing-queries)
- [Using Autocomplete](#using-autocomplete)
- [Keyboard Shortcuts](#keyboard-shortcuts)
- [Tips & Tricks](#tips--tricks)

---

## Connecting to a Database

### Step 1: Start the Proxy

The proxy is the local bridge between your browser and your PostgreSQL database. You can start it in two ways:

#### Option A: Using a Connection String (Quick)

If you already know your connection details, use a PostgreSQL connection string:

```bash
# Basic format
./postgres-proxy "postgres://username:password@hostname:port/database"

# Example with localhost
./postgres-proxy "postgres://myuser:mypass@localhost:5432/mydb"

# Example with remote server
./postgres-proxy "postgres://admin:secret@db.example.com:5432/production"

# Example with SSL required
./postgres-proxy "postgres://user:pass@host:5432/db?sslmode=require"
```

#### Option B: Interactive Mode (Guided)

If you prefer step-by-step prompts, run the proxy without arguments:

```bash
./postgres-proxy
```

You'll be prompted for:
- **Host**: The server hostname (e.g., `localhost` or `db.example.com`)
- **Port**: The PostgreSQL port (default: 5432)
- **Database**: The database name
- **Username**: Your PostgreSQL username
- **Password**: Your password (hidden as you type)
- **SSL Mode**: Connection security (prefer, require, disable, etc.)

The proxy will attempt to connect with retry logic (up to 4 attempts with exponential backoff).

### Step 2: Open the Web App

Once connected, the proxy will display:

```
✓ Connected to postgres://db.example.com:5432/mydb

→ Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6...
```

**Important**: Copy the entire URL including the `?secret=...` parameter. This secret authenticates your browser session.

### Step 3: Start Querying

- Click the URL or paste it into your browser
- The web app will automatically connect to the proxy
- The schema will be introspected (tables, columns, functions)
- You're ready to write queries!

### Connection Status

The connection status badge appears in the top-right corner:

- 🟢 **Connected**: Active connection to database
- 🟡 **Connecting**: Attempting to connect
- 🔴 **Disconnected**: Connection lost or failed

---

## Writing Queries

### Basic Query Execution

1. Type your SQL query in the editor
2. Click **Run Query** or press `Ctrl+Enter` (Windows/Linux) / `Cmd+Enter` (macOS)
3. Results appear in the table below

### Query Examples

#### SELECT Queries

```sql
-- Simple select
SELECT * FROM users;

-- Select with conditions
SELECT id, name, email
FROM users
WHERE created_at > '2024-01-01'
  AND status = 'active';

-- Select with joins
SELECT o.id, o.total, u.name
FROM orders o
JOIN users u ON u.id = o.user_id
WHERE o.total > 100;

-- Aggregate queries
SELECT category, COUNT(*) as count, AVG(price) as avg_price
FROM products
GROUP BY category
ORDER BY avg_price DESC;

-- Window functions
SELECT
  name,
  salary,
  department,
  AVG(salary) OVER (PARTITION BY department) as dept_avg
FROM employees;
```

#### INSERT Queries

```sql
-- Single row insert
INSERT INTO users (name, email, created_at)
VALUES ('Alice', 'alice@example.com', NOW());

-- Multiple row insert
INSERT INTO products (name, price, category)
VALUES
  ('Widget A', 19.99, 'widgets'),
  ('Widget B', 29.99, 'widgets'),
  ('Gadget X', 49.99, 'gadgets');

-- Insert with RETURNING
INSERT INTO users (name, email)
VALUES ('Bob', 'bob@example.com')
RETURNING id, created_at;
```

#### UPDATE Queries

```sql
-- Simple update
UPDATE users
SET status = 'inactive'
WHERE last_login < NOW() - INTERVAL '1 year';

-- Update with RETURNING
UPDATE products
SET price = price * 1.1
WHERE category = 'widgets'
RETURNING id, name, price;

-- Update with join
UPDATE orders o
SET status = 'shipped'
FROM shipments s
WHERE s.order_id = o.id
  AND s.shipped_at IS NOT NULL;
```

#### DELETE Queries

```sql
-- Simple delete
DELETE FROM logs
WHERE created_at < NOW() - INTERVAL '30 days';

-- Delete with RETURNING
DELETE FROM sessions
WHERE expires_at < NOW()
RETURNING id, user_id;
```

#### CREATE Queries

```sql
-- Create table
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status VARCHAR(20) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index
CREATE INDEX idx_tasks_status ON tasks(status);

-- Create view
CREATE VIEW active_tasks AS
SELECT id, title, status
FROM tasks
WHERE status != 'completed';
```

#### PostgreSQL-Specific Features

```sql
-- JSON/JSONB queries
SELECT
  data->>'name' as name,
  data->'address'->>'city' as city
FROM customers
WHERE data @> '{"active": true}';

-- Array operations
SELECT *
FROM products
WHERE tags @> ARRAY['electronics', 'sale'];

-- Common Table Expressions (CTEs)
WITH monthly_sales AS (
  SELECT
    DATE_TRUNC('month', created_at) as month,
    SUM(total) as revenue
  FROM orders
  GROUP BY month
)
SELECT * FROM monthly_sales
ORDER BY month DESC;

-- Recursive queries
WITH RECURSIVE org_chart AS (
  -- Base case
  SELECT id, name, manager_id, 1 as level
  FROM employees
  WHERE manager_id IS NULL

  UNION ALL

  -- Recursive case
  SELECT e.id, e.name, e.manager_id, oc.level + 1
  FROM employees e
  JOIN org_chart oc ON e.manager_id = oc.id
)
SELECT * FROM org_chart;
```

### Understanding Results

#### Result Metadata

After each query, you'll see:

- **Row Count**: Number of rows returned
- **Execution Time**: Query execution duration in milliseconds

#### Data Types in Results

The results table displays column types as badges:

- `INTEGER`, `BIGINT`, `SMALLINT`: Numeric types
- `TEXT`, `VARCHAR`, `CHAR`: String types
- `BOOLEAN`: true/false values
- `TIMESTAMP`, `DATE`, `TIME`: Temporal types
- `JSONB`, `JSON`: JSON data
- `UUID`: UUIDs
- `ARRAY`: Array types

#### NULL Values

NULL values are displayed as `NULL` in gray text to distinguish them from empty strings.

#### Large Results

Results are displayed in a scrollable table with:
- Sticky headers (headers remain visible while scrolling)
- Maximum height of 500px (scroll for more rows)
- Full data displayed (no pagination in MVP)

---

## Using Autocomplete

The SQL editor provides intelligent autocomplete to help you write queries faster and with fewer errors.

### How to Trigger Autocomplete

Autocomplete suggestions appear automatically as you type. You can also trigger it manually:

- **Windows/Linux**: `Ctrl+Space`
- **macOS**: `Cmd+Space`

### What Gets Suggested

#### 1. SQL Keywords

Type the first few letters of a keyword:

```sql
SEL    →  SELECT
FRO    →  FROM
WHE    →  WHERE
ORD    →  ORDER BY
```

**Common keywords**: `SELECT`, `FROM`, `WHERE`, `JOIN`, `LEFT JOIN`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `ALTER`, `DROP`, `GROUP BY`, `ORDER BY`, `LIMIT`, `OFFSET`

#### 2. Database Tables

After typing `FROM` or in a `JOIN` clause:

```sql
SELECT * FROM us    →  users, user_sessions, user_profiles
```

Autocomplete shows:
- ✅ Table name
- 📋 Schema name (e.g., `public.users`)
- 📄 Column list preview

#### 3. Column Names

After typing a table name followed by a dot:

```sql
SELECT users.    →  id, name, email, created_at, updated_at
```

Or anywhere in the query:

```sql
SELECT id, na    →  name, national_id, created_at
```

Autocomplete shows:
- 🔤 Column name
- 📊 Data type
- 🔑 Constraints (if any)

#### 4. Functions

Type the function name:

```sql
SELECT COU    →  COUNT(), COALESCE()
SELECT NOW    →  NOW()
```

**Common functions**:
- **Aggregate**: `COUNT()`, `SUM()`, `AVG()`, `MIN()`, `MAX()`
- **String**: `CONCAT()`, `LOWER()`, `UPPER()`, `SUBSTRING()`, `TRIM()`
- **Date/Time**: `NOW()`, `CURRENT_TIMESTAMP`, `DATE_TRUNC()`, `EXTRACT()`
- **JSON**: `JSON_AGG()`, `JSONB_BUILD_OBJECT()`, `JSONB_SET()`
- **Window**: `ROW_NUMBER()`, `RANK()`, `DENSE_RANK()`, `LAG()`, `LEAD()`

#### 5. PostgreSQL-Specific Syntax

Autocomplete includes Postgres-specific features:

```sql
RETUR    →  RETURNING
CONFL    →  ON CONFLICT
JSONB    →  JSONB
GENER    →  GENERATED ALWAYS AS
```

### Context-Aware Suggestions

The autocomplete is smart about context:

- **After FROM**: Prioritizes table names
- **After SELECT**: Prioritizes column names and functions
- **After WHERE**: Prioritizes column names for conditions
- **After table name + dot**: Shows only columns from that table

### Navigation and Selection

- **Arrow Keys**: Navigate through suggestions
- **Enter/Tab**: Accept selected suggestion
- **Esc**: Dismiss suggestions
- **Continue typing**: Filter suggestions

### Example Workflow

```sql
-- 1. Type "SEL" → select "SELECT"
SELECT

-- 2. Type "us" → select "users.id"
SELECT users.id,

-- 3. Type "na" → select "users.name"
SELECT users.id, users.name

-- 4. Type "FRO" → select "FROM"
SELECT users.id, users.name FROM

-- 5. Type "us" → select "users"
SELECT users.id, users.name FROM users

-- 6. Type "WHE" → select "WHERE"
SELECT users.id, users.name FROM users WHERE

-- 7. Type "created" → select "created_at"
SELECT users.id, users.name FROM users WHERE created_at > '2024-01-01';
```

---

## Keyboard Shortcuts

Master these shortcuts to query faster and more efficiently.

### Query Execution

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+Enter` (Win/Linux)<br>`Cmd+Enter` (macOS) | Run Query | Execute the current SQL query |

### Editor Navigation

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+Space` (Win/Linux)<br>`Cmd+Space` (macOS) | Trigger Autocomplete | Show autocomplete suggestions |
| `Ctrl+F` (Win/Linux)<br>`Cmd+F` (macOS) | Find | Search within editor |
| `Ctrl+H` (Win/Linux)<br>`Cmd+H` (macOS) | Replace | Find and replace in editor |
| `Ctrl+/` (Win/Linux)<br>`Cmd+/` (macOS) | Toggle Comment | Comment/uncomment line or selection |
| `Ctrl+]` (Win/Linux)<br>`Cmd+]` (macOS) | Indent | Indent selected lines |
| `Ctrl+[` (Win/Linux)<br>`Cmd+[` (macOS) | Outdent | Outdent selected lines |

### Text Selection

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+A` (Win/Linux)<br>`Cmd+A` (macOS) | Select All | Select entire query |
| `Shift+Arrow Keys` | Extend Selection | Select text character by character |
| `Ctrl+Shift+Arrow` (Win/Linux)<br>`Cmd+Shift+Arrow` (macOS) | Select Word/Line | Select by word or line |
| `Alt+Click` | Multi-Cursor | Add cursor at clicked position |

### Text Editing

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+D` (Win/Linux)<br>`Cmd+D` (macOS) | Add Selection to Next Find Match | Multi-cursor editing |
| `Alt+Up/Down` | Move Line Up/Down | Reorder lines |
| `Ctrl+Z` (Win/Linux)<br>`Cmd+Z` (macOS) | Undo | Undo last change |
| `Ctrl+Y` (Win/Linux)<br>`Cmd+Shift+Z` (macOS) | Redo | Redo last undone change |

### Monaco Editor Features

The SQL editor is powered by Monaco Editor (the same editor as VS Code), so most VS Code shortcuts work!

---

## Tips & Tricks

### Performance Tips

#### 1. Use EXPLAIN to Understand Query Plans

```sql
-- See how Postgres will execute your query
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';

-- Get detailed execution statistics
EXPLAIN ANALYZE SELECT * FROM large_table WHERE created_at > NOW() - INTERVAL '1 day';
```

#### 2. Limit Large Result Sets During Development

```sql
-- Use LIMIT while developing queries
SELECT * FROM huge_table LIMIT 100;

-- Use OFFSET for pagination
SELECT * FROM huge_table LIMIT 100 OFFSET 200;
```

#### 3. Use Indexes Wisely

```sql
-- Check if a query is using indexes
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';
-- Look for "Index Scan" vs "Seq Scan"

-- Create indexes on frequently queried columns
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

### Query Development Tips

#### 1. Build Queries Incrementally

Start simple and add complexity:

```sql
-- Step 1: Basic select
SELECT * FROM orders;

-- Step 2: Add WHERE clause
SELECT * FROM orders WHERE created_at > '2024-01-01';

-- Step 3: Add JOIN
SELECT o.*, u.name
FROM orders o
JOIN users u ON u.id = o.user_id
WHERE o.created_at > '2024-01-01';

-- Step 4: Add aggregation
SELECT u.name, COUNT(*) as order_count, SUM(o.total) as revenue
FROM orders o
JOIN users u ON u.id = o.user_id
WHERE o.created_at > '2024-01-01'
GROUP BY u.name
ORDER BY revenue DESC;
```

#### 2. Use Comments to Document Complex Queries

```sql
-- Get top 10 customers by revenue in the last quarter
WITH quarterly_orders AS (
  SELECT
    user_id,
    SUM(total) as revenue
  FROM orders
  WHERE created_at >= DATE_TRUNC('quarter', NOW())
  GROUP BY user_id
)
SELECT
  u.name,
  u.email,
  qo.revenue
FROM quarterly_orders qo
JOIN users u ON u.id = qo.user_id
ORDER BY qo.revenue DESC
LIMIT 10;
```

#### 3. Test Modifications in a Transaction

```sql
-- Start a transaction to test changes safely
BEGIN;

-- Make changes
UPDATE users SET status = 'inactive' WHERE last_login < '2023-01-01';

-- Verify the changes
SELECT * FROM users WHERE status = 'inactive';

-- If satisfied, commit
COMMIT;

-- Or rollback to undo
ROLLBACK;
```

### Data Exploration Tips

#### 1. Quick Table Inspection

```sql
-- See first 10 rows
SELECT * FROM table_name LIMIT 10;

-- Count rows
SELECT COUNT(*) FROM table_name;

-- Get distinct values in a column
SELECT DISTINCT category FROM products;

-- See data distribution
SELECT category, COUNT(*) as count
FROM products
GROUP BY category
ORDER BY count DESC;
```

#### 2. Schema Exploration

```sql
-- List all tables in the current database
SELECT tablename
FROM pg_tables
WHERE schemaname = 'public';

-- See columns for a specific table
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'users';

-- Find tables containing a specific column
SELECT table_name
FROM information_schema.columns
WHERE column_name = 'email';

-- See table sizes
SELECT
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### 3. Find Duplicate Records

```sql
-- Find duplicate emails
SELECT email, COUNT(*)
FROM users
GROUP BY email
HAVING COUNT(*) > 1;

-- Find duplicate combinations
SELECT user_id, product_id, COUNT(*)
FROM orders
GROUP BY user_id, product_id
HAVING COUNT(*) > 1;
```

### Advanced PostgreSQL Features

#### 1. JSON/JSONB Operations

```sql
-- Extract JSON fields
SELECT
  id,
  metadata->>'name' as name,
  metadata->'address'->>'city' as city
FROM customers;

-- Query JSON arrays
SELECT *
FROM products
WHERE tags @> '["electronics"]';

-- Aggregate to JSON
SELECT
  category,
  JSON_AGG(JSON_BUILD_OBJECT('id', id, 'name', name)) as products
FROM products
GROUP BY category;
```

#### 2. Full Text Search

```sql
-- Basic text search
SELECT *
FROM articles
WHERE to_tsvector('english', content) @@ to_tsquery('english', 'postgres & database');

-- Ranked search results
SELECT
  title,
  ts_rank(to_tsvector('english', content), query) as rank
FROM articles,
     to_tsquery('english', 'postgres & database') query
WHERE to_tsvector('english', content) @@ query
ORDER BY rank DESC;
```

#### 3. Array Operations

```sql
-- Check if array contains value
SELECT * FROM products WHERE 'electronics' = ANY(tags);

-- Check if array overlaps
SELECT * FROM products WHERE tags && ARRAY['sale', 'clearance'];

-- Array aggregation
SELECT
  category,
  ARRAY_AGG(DISTINCT brand) as brands
FROM products
GROUP BY category;
```

### Working with Query History

#### View Recent Queries

The left sidebar shows your query history with:
- ✅ Successful queries (green badge)
- ❌ Failed queries (red badge)
- ⏱️ Execution time
- 📊 Row count

#### Reuse Previous Queries

- Click any query in the history to load it into the editor
- Modify and re-execute as needed
- History is preserved during your browser session

#### Clear History

Click the "Clear History" button in the query history panel to remove all entries.

### Troubleshooting Common Issues

#### "Connection Lost" Error

**Cause**: Proxy stopped or network issue

**Solution**:
1. Check if the proxy is still running
2. Restart the proxy if needed
3. Refresh the browser page with the same `?secret=...` URL

#### "Invalid Secret" Error

**Cause**: Wrong secret in URL or proxy restarted

**Solution**:
1. Check the proxy output for the current secret
2. Copy the full URL including `?secret=...`
3. Open in a new browser tab

#### "Query Timeout" Error

**Cause**: Query taking longer than timeout limit

**Solution**:
1. Optimize the query (add indexes, use WHERE clauses)
2. Increase timeout in proxy configuration (future feature)
3. Reduce result set size with LIMIT

#### Autocomplete Not Working

**Cause**: Schema not loaded or browser issue

**Solution**:
1. Check connection status (should show "Connected")
2. Refresh the page
3. Try manual trigger: `Ctrl+Space` or `Cmd+Space`

#### Syntax Error in Query

**Cause**: Invalid SQL syntax

**Solution**:
1. Check error message for line/position
2. Verify SQL keyword spelling
3. Check for missing commas, parentheses, or quotes
4. Consult PostgreSQL documentation for correct syntax

---

## Next Steps

- **Explore the Database**: Use schema exploration queries to understand your data
- **Save Common Queries**: Keep a text file of frequently used queries for quick copy-paste
- **Learn PostgreSQL**: Check out the [official PostgreSQL documentation](https://www.postgresql.org/docs/)
- **Provide Feedback**: Report issues or suggest features on our GitHub repository

Happy querying! 🚀
