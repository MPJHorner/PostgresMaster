-- =====================================================
-- Sample Test Queries for PostgresMaster
-- =====================================================
-- Use these queries to test the query execution and
-- results display functionality of the application.
-- =====================================================

-- -----------------------------------------------------
-- 1. SIMPLE SELECT QUERIES
-- -----------------------------------------------------

-- Basic SELECT with a single value
SELECT 1 as num;

-- SELECT with simple arithmetic
SELECT 42 as answer, 'Hello World' as greeting;

-- SELECT current database info
SELECT current_database() as database, current_user as user;


-- -----------------------------------------------------
-- 2. SELECT WITH MULTIPLE COLUMNS AND TYPES
-- -----------------------------------------------------

-- Test all common PostgreSQL data types
SELECT
  1 as int_col,
  'text value' as text_col,
  true as bool_col,
  false as bool_false_col,
  3.14159 as float_col,
  123.456::numeric(10,3) as numeric_col,
  NOW() as timestamp_col,
  CURRENT_DATE as date_col,
  CURRENT_TIME as time_col,
  '1 year 2 months 3 days'::interval as interval_col;


-- -----------------------------------------------------
-- 3. SELECT WITH NULL VALUES
-- -----------------------------------------------------

-- Test NULL handling in results display
SELECT
  1 as has_value,
  NULL as null_value,
  'text' as another_value,
  NULL::integer as typed_null,
  CASE WHEN 1=2 THEN 'yes' ELSE NULL END as conditional_null;


-- -----------------------------------------------------
-- 4. SELECT WITH TIMESTAMPS
-- -----------------------------------------------------

-- Test various timestamp functions
SELECT
  NOW() as now,
  CURRENT_TIMESTAMP as current_timestamp,
  CURRENT_DATE as current_date,
  CURRENT_TIME as current_time,
  NOW() - INTERVAL '1 day' as yesterday,
  NOW() + INTERVAL '1 week' as next_week,
  EXTRACT(YEAR FROM NOW()) as current_year,
  EXTRACT(MONTH FROM NOW()) as current_month,
  EXTRACT(DAY FROM NOW()) as current_day,
  AGE(TIMESTAMP '2000-01-01') as age_since_2000;


-- -----------------------------------------------------
-- 5. SELECT WITH JSON/JSONB
-- -----------------------------------------------------

-- Test JSON data type handling
SELECT
  '{"name": "John", "age": 30}'::json as json_col,
  '{"name": "Jane", "age": 25, "active": true}'::jsonb as jsonb_col,
  JSON_BUILD_OBJECT('key', 'value', 'number', 42) as json_build,
  JSON_AGG(x) as json_array
FROM (SELECT generate_series(1, 5) as x) t;

-- Test JSONB operators
SELECT
  '{"a": 1, "b": {"c": 2}}'::jsonb as original,
  ('{"a": 1, "b": {"c": 2}}'::jsonb)->'b' as nested_access,
  ('{"a": 1, "b": {"c": 2}}'::jsonb)->>'a' as text_value;


-- -----------------------------------------------------
-- 6. SELECT LARGE RESULT SET
-- -----------------------------------------------------

-- Test with 1000 rows to verify scrolling and performance
SELECT
  n as number,
  'Row ' || n as description,
  n * 2 as doubled,
  n % 2 = 0 as is_even,
  NOW() + (n || ' seconds')::interval as offset_time
FROM generate_series(1, 1000) as n;

-- Test with 100 rows and multiple columns
SELECT
  n as id,
  'User ' || n as username,
  'user' || n || '@example.com' as email,
  CASE WHEN n % 3 = 0 THEN 'admin' WHEN n % 3 = 1 THEN 'user' ELSE 'guest' END as role,
  random() > 0.5 as is_active,
  NOW() - (n || ' days')::interval as created_at
FROM generate_series(1, 100) as n;


-- -----------------------------------------------------
-- 7. AGGREGATE QUERIES
-- -----------------------------------------------------

-- Test various aggregate functions
SELECT
  COUNT(*) as total_count,
  SUM(n) as sum,
  AVG(n) as average,
  MIN(n) as minimum,
  MAX(n) as maximum,
  STDDEV(n) as std_deviation,
  VARIANCE(n) as variance
FROM generate_series(1, 100) as n;

-- Test aggregation with GROUP BY
SELECT
  n % 10 as group_key,
  COUNT(*) as count,
  SUM(n) as sum,
  AVG(n)::numeric(10,2) as avg
FROM generate_series(1, 100) as n
GROUP BY n % 10
ORDER BY group_key;


-- -----------------------------------------------------
-- 8. JOIN QUERIES (using CTEs for testing without tables)
-- -----------------------------------------------------

-- Test INNER JOIN with generated data
WITH users AS (
  SELECT n as user_id, 'User ' || n as username
  FROM generate_series(1, 10) as n
),
orders AS (
  SELECT n as order_id, ((n-1) % 10) + 1 as user_id, n * 10 as amount
  FROM generate_series(1, 25) as n
)
SELECT
  u.user_id,
  u.username,
  COUNT(o.order_id) as order_count,
  COALESCE(SUM(o.amount), 0) as total_amount
FROM users u
LEFT JOIN orders o ON u.user_id = o.user_id
GROUP BY u.user_id, u.username
ORDER BY u.user_id;


-- -----------------------------------------------------
-- 9. WINDOW FUNCTIONS
-- -----------------------------------------------------

-- Test window functions
SELECT
  n as value,
  ROW_NUMBER() OVER (ORDER BY n) as row_num,
  RANK() OVER (ORDER BY n % 5) as rank,
  DENSE_RANK() OVER (ORDER BY n % 5) as dense_rank,
  LAG(n, 1) OVER (ORDER BY n) as previous_value,
  LEAD(n, 1) OVER (ORDER BY n) as next_value,
  SUM(n) OVER (ORDER BY n) as running_total,
  AVG(n) OVER (ORDER BY n ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg
FROM generate_series(1, 20) as n;


-- -----------------------------------------------------
-- 10. STRING FUNCTIONS
-- -----------------------------------------------------

-- Test various string manipulation functions
SELECT
  'PostgreSQL' as original,
  UPPER('PostgreSQL') as uppercase,
  LOWER('PostgreSQL') as lowercase,
  LENGTH('PostgreSQL') as length,
  SUBSTRING('PostgreSQL', 1, 4) as substring,
  CONCAT('Postgre', 'SQL') as concatenated,
  REPLACE('PostgreSQL', 'SQL', 'Master') as replaced,
  TRIM('  space  ') as trimmed,
  REVERSE('PostgreSQL') as reversed,
  REPEAT('*', 5) as repeated;


-- -----------------------------------------------------
-- 11. ARRAY DATA TYPE
-- -----------------------------------------------------

-- Test array handling
SELECT
  ARRAY[1, 2, 3, 4, 5] as int_array,
  ARRAY['a', 'b', 'c'] as text_array,
  ARRAY[1, 2, 3] || ARRAY[4, 5, 6] as concatenated_array,
  ARRAY[1, 2, 3, 4, 5][2:4] as array_slice,
  CARDINALITY(ARRAY[1, 2, 3, 4, 5]) as array_length,
  ARRAY_AGG(n) as aggregated_array,
  STRING_AGG(n::text, ', ') as string_aggregated
FROM generate_series(1, 10) as n;


-- -----------------------------------------------------
-- 12. CONDITIONAL LOGIC
-- -----------------------------------------------------

-- Test CASE expressions
SELECT
  n as number,
  CASE
    WHEN n % 15 = 0 THEN 'FizzBuzz'
    WHEN n % 3 = 0 THEN 'Fizz'
    WHEN n % 5 = 0 THEN 'Buzz'
    ELSE n::text
  END as fizzbuzz,
  COALESCE(NULLIF(n % 3, 0), -1) as nullif_example
FROM generate_series(1, 30) as n;


-- -----------------------------------------------------
-- 13. SUBQUERIES
-- -----------------------------------------------------

-- Test subqueries
SELECT
  main.n as number,
  (SELECT COUNT(*) FROM generate_series(1, main.n)) as count_up_to,
  main.n IN (SELECT generate_series(1, 10)) as is_in_first_ten
FROM (SELECT generate_series(1, 20) as n) main
WHERE main.n <= 15;


-- -----------------------------------------------------
-- 14. COMMON TABLE EXPRESSIONS (CTEs)
-- -----------------------------------------------------

-- Test recursive CTE
WITH RECURSIVE fibonacci(n, fib_n, fib_n_plus_1) AS (
  SELECT 1, 0::bigint, 1::bigint
  UNION ALL
  SELECT n + 1, fib_n_plus_1, fib_n + fib_n_plus_1
  FROM fibonacci
  WHERE n < 20
)
SELECT n, fib_n as fibonacci_number
FROM fibonacci;


-- -----------------------------------------------------
-- 15. INSERT QUERY (into temporary table)
-- -----------------------------------------------------

-- Create a temporary table and insert data
CREATE TEMP TABLE IF NOT EXISTS test_users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Insert a single row
INSERT INTO test_users (username, email)
VALUES ('alice', 'alice@example.com');

-- Insert multiple rows
INSERT INTO test_users (username, email)
VALUES
  ('bob', 'bob@example.com'),
  ('charlie', 'charlie@example.com'),
  ('diana', 'diana@example.com');

-- Insert with RETURNING clause
INSERT INTO test_users (username, email)
VALUES ('eve', 'eve@example.com')
RETURNING *;


-- -----------------------------------------------------
-- 16. UPDATE QUERY
-- -----------------------------------------------------

-- Update single row
UPDATE test_users
SET email = 'alice.updated@example.com'
WHERE username = 'alice';

-- Update with RETURNING clause
UPDATE test_users
SET email = LOWER(email),
    created_at = NOW()
WHERE username IN ('bob', 'charlie')
RETURNING *;


-- -----------------------------------------------------
-- 17. DELETE QUERY
-- -----------------------------------------------------

-- Delete with condition
DELETE FROM test_users
WHERE username = 'eve';

-- Delete with RETURNING clause
DELETE FROM test_users
WHERE id = (SELECT MIN(id) FROM test_users)
RETURNING *;


-- -----------------------------------------------------
-- 18. SELECT FROM TEMP TABLE
-- -----------------------------------------------------

-- Verify the data after INSERT/UPDATE/DELETE
SELECT * FROM test_users ORDER BY id;

-- Count records
SELECT COUNT(*) as total_users FROM test_users;


-- -----------------------------------------------------
-- 19. QUERY WITH SYNTAX ERROR
-- -----------------------------------------------------

-- Uncomment to test error handling:
-- SELCT * FROM test_users;  -- Typo: SELCT instead of SELECT


-- -----------------------------------------------------
-- 20. QUERY ON NON-EXISTENT TABLE
-- -----------------------------------------------------

-- Uncomment to test error handling:
-- SELECT * FROM nonexistent_table_xyz;


-- -----------------------------------------------------
-- 21. QUERY WITH MISSING COLUMN
-- -----------------------------------------------------

-- Uncomment to test error handling:
-- SELECT id, username, nonexistent_column FROM test_users;


-- -----------------------------------------------------
-- 22. INVALID SQL SYNTAX
-- -----------------------------------------------------

-- Uncomment to test error handling:
-- SELECT * FROM WHERE id = 1;


-- -----------------------------------------------------
-- 23. DIVISION BY ZERO ERROR
-- -----------------------------------------------------

-- Uncomment to test runtime error handling:
-- SELECT 10 / 0 as division_by_zero;


-- -----------------------------------------------------
-- 24. COMPLEX QUERY WITH MULTIPLE FEATURES
-- -----------------------------------------------------

-- A comprehensive query combining many features
WITH ranked_data AS (
  SELECT
    n as id,
    'Item ' || n as name,
    (n % 5) + 1 as category_id,
    random() * 100 as price,
    ROW_NUMBER() OVER (PARTITION BY (n % 5) + 1 ORDER BY random()) as rank_in_category
  FROM generate_series(1, 50) as n
),
category_stats AS (
  SELECT
    category_id,
    COUNT(*) as item_count,
    AVG(price)::numeric(10,2) as avg_price,
    MIN(price)::numeric(10,2) as min_price,
    MAX(price)::numeric(10,2) as max_price
  FROM ranked_data
  GROUP BY category_id
)
SELECT
  rd.id,
  rd.name,
  rd.category_id,
  rd.price::numeric(10,2),
  rd.rank_in_category,
  cs.avg_price as category_avg_price,
  CASE
    WHEN rd.price > cs.avg_price THEN 'Above Average'
    WHEN rd.price < cs.avg_price THEN 'Below Average'
    ELSE 'Average'
  END as price_category,
  JSON_BUILD_OBJECT(
    'id', rd.id,
    'name', rd.name,
    'price', rd.price::numeric(10,2),
    'category', rd.category_id
  ) as json_representation
FROM ranked_data rd
JOIN category_stats cs ON rd.category_id = cs.category_id
WHERE rd.rank_in_category <= 3
ORDER BY rd.category_id, rd.rank_in_category;


-- -----------------------------------------------------
-- 25. VACUUM AND ANALYZE (maintenance queries)
-- -----------------------------------------------------

-- These would normally be administrative queries
-- VACUUM test_users;
-- ANALYZE test_users;


-- -----------------------------------------------------
-- 26. EXPLAIN QUERY (for query planning)
-- -----------------------------------------------------

-- Test EXPLAIN to see query execution plan
EXPLAIN SELECT * FROM generate_series(1, 100) WHERE generate_series % 2 = 0;

-- Test EXPLAIN ANALYZE for actual execution statistics
-- Note: This actually runs the query
EXPLAIN ANALYZE
SELECT COUNT(*)
FROM generate_series(1, 10000) as n
WHERE n % 7 = 0;


-- =====================================================
-- END OF TEST QUERIES
-- =====================================================
--
-- USAGE NOTES:
-- - Most queries are ready to run as-is
-- - Some queries with errors are commented out
-- - Uncomment error queries to test error handling
-- - The temp table queries should be run in sequence
-- - Use these to test various UI features:
--   * Results display
--   * Column type formatting
--   * NULL value handling
--   * Large result sets
--   * Error messages
--   * Execution time tracking
-- =====================================================
