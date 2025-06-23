package sqlformatter

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter()
	if formatter.IndentSize != 2 {
		t.Errorf("Expected IndentSize to be 2, got %d", formatter.IndentSize)
	}
	if !formatter.KeywordUpper {
		t.Errorf("Expected KeywordUpper to be true")
	}
}

func TestSelectFormatting(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Simple SELECT",
			input: "SELECT * FROM users WHERE id = 1",
			expected: `SELECT
  *
FROM
  users
WHERE
  id = 1`,
		},
		{
			name:  "Multi-column SELECT",
			input: "SELECT id, name, email FROM users",
			expected: `SELECT
  id,
  name,
  email
FROM
  users`,
		},
		{
			name:  "Complex SELECT with JOIN",
			input: "select u.id, p.product_name, u.name from users u join products p on u.id = p.user_id where u.age > 25 and p.category = 'electronics' group by u.id order by p.price desc limit 10",
			expected: `SELECT
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
  10`,
		},
		{
			name:  "SELECT with aggregate functions",
			input: "SELECT COUNT(*) as total, category FROM products GROUP BY category HAVING COUNT(*) > 5",
			expected: `SELECT
  COUNT(*) as total,
  category
FROM
  products
GROUP BY
  category
HAVING
  COUNT(*) > 5`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(tt.input)
			if err != nil {
				t.Fatalf("Formatting failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Formatting result mismatch\nExpected:\n%s\nActual:\n%s", tt.expected, result)
			}
		})
	}
}

func TestInsertFormatting(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Simple INSERT",
			input: "INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')",
			expected: `INSERT INTO users
  (name, email)
VALUES
  ('Alice', 'alice@example.com')`,
		},
		{
			name:  "Multi-column INSERT",
			input: "INSERT INTO products (name, price, category, description) VALUES ('Laptop', 999.99, 'electronics', 'High-performance laptop')",
			expected: `INSERT INTO products
  (name, price, category, description)
VALUES
  ('Laptop', 999.99, 'electronics', 'High-performance laptop')`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(tt.input)
			if err != nil {
				t.Fatalf("Formatting failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Formatting result mismatch\nExpected:\n%s\nActual:\n%s", tt.expected, result)
			}
		})
	}
}

func TestUpdateFormatting(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Simple UPDATE",
			input: "UPDATE users SET name = 'John' WHERE id = 1",
			expected: `UPDATE users
SET
  name = 'John'
WHERE
  id = 1`,
		},
		{
			name:  "Multi-field UPDATE",
			input: "UPDATE users SET name = 'John', email = 'john@example.com' WHERE id = 1",
			expected: `UPDATE users
SET
  name = 'John',
  email = 'john@example.com'
WHERE
  id = 1`,
		},
		{
			name:  "UPDATE without WHERE",
			input: "UPDATE users SET active = true",
			expected: `UPDATE users
SET
  active = true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(tt.input)
			if err != nil {
				t.Fatalf("Formatting failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Formatting result mismatch\nExpected:\n%s\nActual:\n%s", tt.expected, result)
			}
		})
	}
}

func TestDeleteFormatting(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "DELETE with WHERE",
			input: "DELETE FROM users WHERE age < 18",
			expected: `DELETE FROM users
WHERE
  age < 18`,
		},
		{
			name:  "DELETE with complex WHERE condition",
			input: "DELETE FROM orders WHERE status = 'cancelled' AND created_at < '2023-01-01'",
			expected: `DELETE FROM orders
WHERE
  status = 'cancelled' AND created_at < '2023-01-01'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(tt.input)
			if err != nil {
				t.Fatalf("Formatting failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Formatting result mismatch\nExpected:\n%s\nActual:\n%s", tt.expected, result)
			}
		})
	}
}

func TestFormatterOptions(t *testing.T) {
	t.Run("Custom indent size", func(t *testing.T) {
		formatter := NewFormatter()
		formatter.IndentSize = 4

		input := "SELECT id, name FROM users"
		result, err := formatter.Format(input)
		if err != nil {
			t.Fatalf("Formatting failed: %v", err)
		}

		// Check if 4 spaces indentation is used
		lines := strings.Split(result, "\n")
		if len(lines) < 2 {
			t.Fatal("Expected at least 2 lines of output")
		}
		if !strings.HasPrefix(lines[1], "    ") { // 4 spaces
			t.Errorf("Expected 4-space indentation, got indentation: '%s'", lines[1][:4])
		}
	})

	t.Run("Lowercase keywords", func(t *testing.T) {
		formatter := NewFormatter()
		formatter.KeywordUpper = false

		input := "SELECT * FROM users"
		result, err := formatter.Format(input)
		if err != nil {
			t.Fatalf("Formatting failed: %v", err)
		}

		if !strings.Contains(result, "select") || !strings.Contains(result, "from") {
			t.Errorf("Expected lowercase keywords, actual result: %s", result)
		}
	})
}

func TestErrorCases(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty SQL",
			input: "",
		},
		{
			name:  "Whitespace-only SQL",
			input: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := formatter.Format(tt.input)
			if err == nil {
				t.Errorf("Expected error to occur, but no error happened")
			}
		})
	}
}

func TestSplitColumns(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple column split",
			input:    "id, name, email",
			expected: []string{"id", " name", " email"},
		},
		{
			name:     "Column split with functions",
			input:    "id, COUNT(*), MAX(age)",
			expected: []string{"id", " COUNT(*)", " MAX(age)"},
		},
		{
			name:     "Column split with nested parentheses",
			input:    "id, CASE WHEN (age > 18) THEN 'adult' ELSE 'minor' END, name",
			expected: []string{"id", " CASE WHEN (age > 18) THEN 'adult' ELSE 'minor' END", " name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.splitColumns(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d parts, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Part %d mismatch, expected '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}

func TestKeywordFormatting(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		name     string
		input    string
		upper    bool
		expected string
	}{
		{
			name:     "Uppercase keyword",
			input:    "select",
			upper:    true,
			expected: "SELECT",
		},
		{
			name:     "Lowercase keyword",
			input:    "SELECT",
			upper:    false,
			expected: "select",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter.KeywordUpper = tt.upper
			result := formatter.keyword(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSelectFormatting(b *testing.B) {
	formatter := NewFormatter()
	sql := "select u.id, p.product_name, u.name from users u join products p on u.id = p.user_id where u.age > 25 and p.category = 'electronics' group by u.id order by p.price desc limit 10"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(sql)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInsertFormatting(b *testing.B) {
	formatter := NewFormatter()
	sql := "INSERT INTO products (name, price, category, description) VALUES ('Laptop', 999.99, 'electronics', 'High-performance laptop')"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.Format(sql)
		if err != nil {
			b.Fatal(err)
		}
	}
}
