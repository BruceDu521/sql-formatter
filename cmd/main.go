package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	sqlformatter "github.com/BruceDu521/sql-formatter"
)

var (
	// Command line flags
	indentSize   = flag.Int("indent", 2, "Number of spaces for indentation")
	keywordUpper = flag.Bool("uppercase", true, "Use uppercase for keywords")
	inputFile    = flag.String("input", "", "Input SQL file")
	outputFile   = flag.String("output", "", "Output file")
	sqlString    = flag.String("sql", "", "SQL statement to format")
	showHelp     = flag.Bool("help", false, "Show help information")
)

func main() {
	flag.Parse()

	if *showHelp {
		showUsage()
		return
	}

	// 创建格式化器
	formatter := sqlformatter.NewFormatter()
	formatter.IndentSize = *indentSize
	formatter.KeywordUpper = *keywordUpper

	var sql string
	var err error

	// Get SQL input
	if *sqlString != "" {
		// From command line argument
		sql = *sqlString
	} else if len(flag.Args()) > 0 {
		// From positional arguments
		sql = strings.Join(flag.Args(), " ")
	} else if *inputFile != "" {
		// From file
		sql, err = readFromFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// From stdin
		sql, err = readFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
	}

	if strings.TrimSpace(sql) == "" {
		fmt.Fprintf(os.Stderr, "Error: No SQL statement provided\n")
		showUsage()
		os.Exit(1)
	}

	// Format SQL
	formatted, err := formatter.Format(sql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Formatting failed: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if *outputFile != "" {
		err = writeToFile(*outputFile, formatted)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Formatting completed, result saved to: %s\n", *outputFile)
	} else {
		fmt.Println(formatted)
	}
}

func showUsage() {
	fmt.Printf(`sqlformatter - SQL Formatter Tool

Usage:
  sqlformatter [options] [SQL statement]

Options:
  -sql string      SQL statement to format
  -input string    Input SQL file
  -output string   Output file
  -indent int      Number of spaces for indentation (default: 2)
  -uppercase       Use uppercase for keywords (default: true)
  -help            Show help information

Examples:
  sqlformatter -sql "select * from users"
  echo "select * from users" | sqlformatter
  sqlformatter -input input.sql -output output.sql
  sqlformatter "select u.id, u.name from users u where u.age > 25"

`)
}

// readFromStdin reads from standard input
func readFromStdin() (string, error) {
	// Check if there's piped input
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// Interactive mode
		fmt.Print("Enter SQL statement (press Ctrl+D to finish):\n")
	}

	var result strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		result.WriteString(scanner.Text())
		result.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.String()), nil
}

// readFromFile reads from a file
func readFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// writeToFile writes content to a file
func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
