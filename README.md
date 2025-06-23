# SQL Formatter

A SQL formatter tool written in Go that can be used both as a library and as a CLI command-line tool.

[中文说明](README_zh.md)

## Features

- 🔧 Format compressed SQL statements into readable multi-line format
- 📚 Can be integrated into your projects as a Go library
- 💻 Provides command-line tool with multiple input methods
- ⚙️ Configurable formatting options (indent size, keyword case, etc.)
- 🚀 Supports complex SQL statements including JOIN, subqueries, etc.

## Supported SQL Statements

- SELECT (including JOIN, WHERE, GROUP BY, HAVING, ORDER BY, LIMIT, etc.)
- INSERT
- UPDATE  
- DELETE

## Installation

### As a Library

```bash
go get github.com/BruceDu521/sql-formatter
```

### Compile CLI Tool

```bash
go build -o sqlformatter cmd/main.go
```

Or download pre-built binaries from [Releases](https://github.com/BruceDu521/sql-formatter/releases).

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "log"
    
    sqlformatter "github.com/BruceDu521/sql-formatter"
)

func main() {
    // Create formatter
    formatter := sqlformatter.NewFormatter()
    
    // Configuration options
    formatter.IndentSize = 2      // Number of spaces for indentation
    formatter.KeywordUpper = true // Use uppercase for keywords
    
    // Format SQL
    sql := "select u.id, u.name from users u where u.age > 25"
    formatted, err := formatter.Format(sql)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(formatted)
}
```

### CLI Command Line Tool

#### Basic Usage

```bash
# Direct SQL input
./sqlformatter -sql "select * from users where id = 1"

# Read from file
./sqlformatter -input input.sql

# Output to file
./sqlformatter -input input.sql -output formatted.sql

# Use pipe
echo "select * from users" | ./sqlformatter
```

#### Command Line Options

```
Options:
  -sql string      SQL statement to format
  -input string    Input SQL file
  -output string   Output file
  -indent int      Number of spaces for indentation (default: 2)
  -uppercase       Use uppercase for keywords (default: true)
  -help            Show help information
```

## Examples

### SELECT Statement

**Input:**
```sql
select u.id, p.product_name, u.name from users u join products p on u.id = p.user_id where u.age > 25 and p.category = 'electronics' group by u.id order by p.price desc limit 10
```

**Output:**
```sql
SELECT
  u.id,
  p.product_name,
  u.name
FROM
  users u
  JOIN products p on u.id = p.user_id
WHERE
  u.age > 25 and p.category = 'electronics'
GROUP BY
  u.id
ORDER BY
  p.price desc
LIMIT
  10
```

### INSERT Statement

**Input:**
```sql
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')
```

**Output:**
```sql
INSERT INTO users
  (name, email)
VALUES
  ('Alice', 'alice@example.com')
```

### UPDATE Statement

**Input:**
```sql
UPDATE users SET name = 'John', email = 'john@example.com' WHERE id = 1
```

**Output:**
```sql
UPDATE users
SET
  name = 'John',
  email = 'john@example.com'
WHERE
  id = 1
```

### DELETE Statement

**Input:**
```sql
DELETE FROM users WHERE age < 18
```

**Output:**
```sql
DELETE FROM users
WHERE
  age < 18
```

## Project Structure

```
sql-formatter/
├── formatter.go        # Core formatting logic
├── cmd/
│   └── main.go        # CLI tool
├── example/
│   └── main.go        # Usage examples
├── formatter_test.go  # Unit tests
├── go.mod             # Go module file
├── README.md          # Project documentation (English)
├── README_zh.md       # Project documentation (Chinese)
├── test.sql           # Test SQL file
└── formatted.sql      # Formatted output example
```

## Running Examples

```bash
# Run library usage example
go run example/main.go

# Run CLI tool
go run cmd/main.go -help
go run cmd/main.go -sql "select * from users"
go run cmd/main.go -input test.sql

# Run tests
go test -v

# Run benchmarks
go test -bench=.
```

## Development Roadmap

- [ ] Support more SQL statement types (CREATE TABLE, ALTER TABLE, etc.)
- [ ] Add more formatting options
- [ ] Support different database dialects

## Contributing

Issues and Pull Requests are welcome!

## License

MIT License 