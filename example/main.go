package main

import (
	"fmt"
	"log"
	"strings"

	sqlformatter "github.com/BruceDu521/sql-formatter"
)

func main() {
	// Create formatter
	formatter := sqlformatter.NewFormatter()

	// Test SQL statement
	testSQL := "select u.id, p.product_name, u.name from users u join products p on u.id = p.user_id where u.age > 25 and p.category = 'electronics' group by u.id order by p.price desc limit 10"

	fmt.Println("Original SQL:")
	fmt.Println(testSQL)
	fmt.Println("\nFormatted SQL:")

	// Format SQL
	formatted, err := formatter.Format(testSQL)
	if err != nil {
		log.Fatalf("Formatting failed: %v", err)
	}

	fmt.Println(formatted)

	// Test other SQL statements
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("More examples:")

	examples := []string{
		"SELECT * FROM users WHERE id = 1",
		"SELECT COUNT(*) as total, category FROM products GROUP BY category HAVING COUNT(*) > 5",
		"UPDATE users SET name = 'John' WHERE id = 1",
		"INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')",
		"DELETE FROM users WHERE age < 18",
	}

	for i, sql := range examples {
		fmt.Printf("\nExample %d:\n", i+1)
		fmt.Printf("Original: %s\n", sql)

		formatted, err := formatter.Format(sql)
		if err != nil {
			fmt.Printf("Formatting failed: %v\n", err)
			continue
		}

		fmt.Printf("Formatted:\n%s\n", formatted)
	}
}
