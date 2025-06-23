package sqlformatter

import (
	"fmt"
	"regexp"
	"strings"
)

// Formatter SQL formatter configuration
type Formatter struct {
	IndentSize   int
	KeywordUpper bool
}

// NewFormatter creates a new formatter instance
func NewFormatter() *Formatter {
	return &Formatter{
		IndentSize:   2,
		KeywordUpper: true,
	}
}

// Format formats the given SQL statement
func (f *Formatter) Format(sql string) (string, error) {
	if strings.TrimSpace(sql) == "" {
		return "", fmt.Errorf("SQL statement cannot be empty")
	}

	// 清理并标准化SQL
	cleaned := f.cleanSQL(sql)

	// 格式化SQL
	formatted := f.formatSQL(cleaned)

	return formatted, nil
}

// cleanSQL 清理SQL语句
func (f *Formatter) cleanSQL(sql string) string {
	// 移除多余的空白字符
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	// 移除首尾空白
	sql = strings.TrimSpace(sql)
	return sql
}

// formatSQL 格式化SQL语句
func (f *Formatter) formatSQL(sql string) string {
	// SQL关键字模式
	keywords := []string{
		"SELECT", "FROM", "WHERE", "GROUP BY", "HAVING", "ORDER BY", "LIMIT",
		"INSERT", "INTO", "VALUES", "UPDATE", "SET", "DELETE",
		"JOIN", "INNER JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN",
		"UNION", "UNION ALL", "CASE", "WHEN", "THEN", "ELSE", "END",
	}

	// 对关键字进行处理
	for _, keyword := range keywords {
		pattern := `(?i)\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)
		replacement := f.keyword(keyword)
		sql = re.ReplaceAllString(sql, replacement)
	}

	// 检测SQL类型并格式化
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if strings.HasPrefix(sqlUpper, "SELECT") {
		return f.formatSelectStatement(sql)
	} else if strings.HasPrefix(sqlUpper, "INSERT") {
		return f.formatInsertStatement(sql)
	} else if strings.HasPrefix(sqlUpper, "UPDATE") {
		return f.formatUpdateStatement(sql)
	} else if strings.HasPrefix(sqlUpper, "DELETE") {
		return f.formatDeleteStatement(sql)
	}

	return sql
}

// formatSelectStatement 格式化SELECT语句
func (f *Formatter) formatSelectStatement(sql string) string {
	// 使用正则表达式分割SQL的各个部分
	parts := f.splitSelectSQL(sql)

	var result strings.Builder
	indent := f.getIndent(1)

	// SELECT部分
	if selectPart := parts["SELECT"]; selectPart != "" {
		result.WriteString(f.keyword("SELECT"))
		result.WriteString("\n")
		selectColumns := f.formatSelectColumns(selectPart)
		result.WriteString(indent + selectColumns)
	}

	// FROM部分
	if fromPart := parts["FROM"]; fromPart != "" {
		result.WriteString("\n" + f.keyword("FROM"))
		result.WriteString("\n")
		fromClause := f.formatFromClause(fromPart)
		result.WriteString(indent + fromClause)
	}

	// WHERE部分
	if wherePart := parts["WHERE"]; wherePart != "" {
		result.WriteString("\n" + f.keyword("WHERE"))
		result.WriteString("\n")
		result.WriteString(indent + wherePart)
	}

	// GROUP BY部分
	if groupByPart := parts["GROUP BY"]; groupByPart != "" {
		result.WriteString("\n" + f.keyword("GROUP BY"))
		result.WriteString("\n")
		result.WriteString(indent + groupByPart)
	}

	// HAVING部分
	if havingPart := parts["HAVING"]; havingPart != "" {
		result.WriteString("\n" + f.keyword("HAVING"))
		result.WriteString("\n")
		result.WriteString(indent + havingPart)
	}

	// ORDER BY部分
	if orderByPart := parts["ORDER BY"]; orderByPart != "" {
		result.WriteString("\n" + f.keyword("ORDER BY"))
		result.WriteString("\n")
		result.WriteString(indent + orderByPart)
	}

	// LIMIT部分
	if limitPart := parts["LIMIT"]; limitPart != "" {
		result.WriteString("\n" + f.keyword("LIMIT"))
		result.WriteString("\n")
		result.WriteString(indent + limitPart)
	}

	return result.String()
}

// splitSelectSQL 分割SELECT SQL的各个部分
func (f *Formatter) splitSelectSQL(sql string) map[string]string {
	parts := make(map[string]string)

	// 定义关键字的正则表达式模式
	patterns := map[string]string{
		"SELECT":   `(?i)\bSELECT\s+(.*?)(?:\s+FROM|\s*$)`,
		"FROM":     `(?i)\bFROM\s+(.*?)(?:\s+WHERE|\s+GROUP\s+BY|\s+HAVING|\s+ORDER\s+BY|\s+LIMIT|\s*$)`,
		"WHERE":    `(?i)\bWHERE\s+(.*?)(?:\s+GROUP\s+BY|\s+HAVING|\s+ORDER\s+BY|\s+LIMIT|\s*$)`,
		"GROUP BY": `(?i)\bGROUP\s+BY\s+(.*?)(?:\s+HAVING|\s+ORDER\s+BY|\s+LIMIT|\s*$)`,
		"HAVING":   `(?i)\bHAVING\s+(.*?)(?:\s+ORDER\s+BY|\s+LIMIT|\s*$)`,
		"ORDER BY": `(?i)\bORDER\s+BY\s+(.*?)(?:\s+LIMIT|\s*$)`,
		"LIMIT":    `(?i)\bLIMIT\s+(.*?)(?:\s*$)`,
	}

	for keyword, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			parts[keyword] = strings.TrimSpace(matches[1])
		}
	}

	return parts
}

// formatSelectColumns 格式化SELECT列
func (f *Formatter) formatSelectColumns(selectPart string) string {
	// 分割列名
	columns := f.splitColumns(selectPart)

	if len(columns) == 0 {
		return selectPart
	}

	var result strings.Builder
	for i, col := range columns {
		if i > 0 {
			result.WriteString(",\n" + f.getIndent(1))
		}
		result.WriteString(strings.TrimSpace(col))
	}

	return result.String()
}

// formatFromClause 格式化FROM子句
func (f *Formatter) formatFromClause(fromPart string) string {
	// 处理JOIN
	joinPattern := `(?i)\b(INNER\s+JOIN|LEFT\s+JOIN|RIGHT\s+JOIN|FULL\s+JOIN|JOIN)\b`
	re := regexp.MustCompile(joinPattern)

	// 分割JOIN部分
	parts := re.Split(fromPart, -1)
	joins := re.FindAllString(fromPart, -1)

	var result strings.Builder
	result.WriteString(strings.TrimSpace(parts[0])) // 主表

	for i, join := range joins {
		if i+1 < len(parts) {
			result.WriteString("\n" + f.getIndent(1))
			result.WriteString(f.keyword(join) + " " + strings.TrimSpace(parts[i+1]))
		}
	}

	return result.String()
}

// splitColumns 分割列名（考虑函数调用中的逗号）
func (f *Formatter) splitColumns(columnsStr string) []string {
	var columns []string
	var current strings.Builder
	var parenCount int

	for _, char := range columnsStr {
		switch char {
		case '(':
			parenCount++
			current.WriteRune(char)
		case ')':
			parenCount--
			current.WriteRune(char)
		case ',':
			if parenCount == 0 {
				columns = append(columns, current.String())
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		columns = append(columns, current.String())
	}

	return columns
}

// formatInsertStatement 格式化INSERT语句
func (f *Formatter) formatInsertStatement(sql string) string {
	// INSERT INTO table (col1, col2) VALUES (val1, val2)
	insertPattern := `(?i)\bINSERT\s+INTO\s+(\S+)\s*\(([^)]+)\)\s+VALUES\s*\(([^)]+)\)`
	re := regexp.MustCompile(insertPattern)
	matches := re.FindStringSubmatch(sql)

	if len(matches) >= 4 {
		tableName := strings.TrimSpace(matches[1])
		columns := strings.TrimSpace(matches[2])
		values := strings.TrimSpace(matches[3])

		var result strings.Builder
		indent := f.getIndent(1)

		result.WriteString(f.keyword("INSERT INTO") + " " + tableName)
		result.WriteString("\n" + indent + "(" + f.formatColumnList(columns) + ")")
		result.WriteString("\n" + f.keyword("VALUES"))
		result.WriteString("\n" + indent + "(" + f.formatValueList(values) + ")")

		return result.String()
	}

	// 如果不匹配标准格式，返回格式化的关键字版本
	return f.formatKeywords(sql)
}

// formatUpdateStatement 格式化UPDATE语句
func (f *Formatter) formatUpdateStatement(sql string) string {
	// 分割UPDATE语句的各个部分
	parts := f.splitUpdateSQL(sql)

	var result strings.Builder
	indent := f.getIndent(1)

	// UPDATE部分
	if updatePart := parts["UPDATE"]; updatePart != "" {
		result.WriteString(f.keyword("UPDATE") + " " + updatePart)
	}

	// SET部分
	if setPart := parts["SET"]; setPart != "" {
		result.WriteString("\n" + f.keyword("SET"))
		result.WriteString("\n" + indent + f.formatSetClause(setPart))
	}

	// WHERE部分
	if wherePart := parts["WHERE"]; wherePart != "" {
		result.WriteString("\n" + f.keyword("WHERE"))
		result.WriteString("\n" + indent + wherePart)
	}

	return result.String()
}

// formatDeleteStatement 格式化DELETE语句
func (f *Formatter) formatDeleteStatement(sql string) string {
	// 分割DELETE语句的各个部分
	parts := f.splitDeleteSQL(sql)

	var result strings.Builder
	indent := f.getIndent(1)

	// DELETE FROM部分
	if fromPart := parts["FROM"]; fromPart != "" {
		result.WriteString(f.keyword("DELETE FROM") + " " + fromPart)
	}

	// WHERE部分
	if wherePart := parts["WHERE"]; wherePart != "" {
		result.WriteString("\n" + f.keyword("WHERE"))
		result.WriteString("\n" + indent + wherePart)
	}

	return result.String()
}

// formatColumnList 格式化列列表
func (f *Formatter) formatColumnList(columns string) string {
	cols := f.splitColumns(columns)
	if len(cols) <= 1 {
		return columns
	}

	var result strings.Builder
	for i, col := range cols {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(strings.TrimSpace(col))
	}
	return result.String()
}

// formatValueList 格式化值列表
func (f *Formatter) formatValueList(values string) string {
	vals := f.splitColumns(values) // 复用splitColumns逻辑
	if len(vals) <= 1 {
		return values
	}

	var result strings.Builder
	for i, val := range vals {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(strings.TrimSpace(val))
	}
	return result.String()
}

// formatSetClause 格式化SET子句
func (f *Formatter) formatSetClause(setPart string) string {
	// 分割SET子句中的赋值语句
	assignments := f.splitColumns(setPart)
	if len(assignments) <= 1 {
		return setPart
	}

	var result strings.Builder
	indent := f.getIndent(1)
	for i, assignment := range assignments {
		if i > 0 {
			result.WriteString(",\n" + indent)
		}
		result.WriteString(strings.TrimSpace(assignment))
	}
	return result.String()
}

// splitUpdateSQL 分割UPDATE SQL的各个部分
func (f *Formatter) splitUpdateSQL(sql string) map[string]string {
	parts := make(map[string]string)

	patterns := map[string]string{
		"UPDATE": `(?i)\bUPDATE\s+(.*?)(?:\s+SET|\s*$)`,
		"SET":    `(?i)\bSET\s+(.*?)(?:\s+WHERE|\s*$)`,
		"WHERE":  `(?i)\bWHERE\s+(.*?)(?:\s*$)`,
	}

	for keyword, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			parts[keyword] = strings.TrimSpace(matches[1])
		}
	}

	return parts
}

// splitDeleteSQL 分割DELETE SQL的各个部分
func (f *Formatter) splitDeleteSQL(sql string) map[string]string {
	parts := make(map[string]string)

	patterns := map[string]string{
		"FROM":  `(?i)\bDELETE\s+FROM\s+(.*?)(?:\s+WHERE|\s*$)`,
		"WHERE": `(?i)\bWHERE\s+(.*?)(?:\s*$)`,
	}

	for keyword, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(sql)
		if len(matches) > 1 {
			parts[keyword] = strings.TrimSpace(matches[1])
		}
	}

	return parts
}

// formatKeywords 格式化关键字（备用方法）
func (f *Formatter) formatKeywords(sql string) string {
	keywords := []string{
		"INSERT INTO", "VALUES", "UPDATE", "SET", "DELETE FROM", "WHERE",
	}

	result := sql
	for _, keyword := range keywords {
		pattern := `(?i)\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, f.keyword(keyword))
	}

	return result
}

// keyword 处理关键字大小写
func (f *Formatter) keyword(word string) string {
	if f.KeywordUpper {
		return strings.ToUpper(word)
	}
	return strings.ToLower(word)
}

// getIndent 获取缩进字符串
func (f *Formatter) getIndent(level int) string {
	return strings.Repeat(" ", level*f.IndentSize)
}
