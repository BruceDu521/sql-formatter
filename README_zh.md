# SQL Formatter

一个用Go语言编写的SQL格式化工具，既可以作为库使用，也提供CLI命令行工具。

[English Documentation](README.md)

## 功能特性

- 🔧 将压缩的SQL语句格式化为易读的多行格式
- 📚 可作为Go语言库集成到你的项目中
- 💻 提供命令行工具，支持多种输入方式
- ⚙️ 可配置的格式化选项（缩进大小、关键字大小写等）
- 🚀 支持复杂的SQL语句，包括JOIN、子查询等

## 支持的SQL语句

- SELECT（包括JOIN、WHERE、GROUP BY、HAVING、ORDER BY、LIMIT等）
- INSERT
- UPDATE  
- DELETE

## 安装

### 作为库使用

```bash
go get github.com/BruceDu521/sql-formatter
```

### 编译CLI工具

```bash
go build -o sqlformatter cmd/main.go
```

或从 [Releases](https://github.com/BruceDu521/sql-formatter/releases) 下载预编译的二进制文件。

## 使用方法

### 作为库使用

```go
package main

import (
    "fmt"
    "log"
    
    sqlformatter "github.com/BruceDu521/sql-formatter"
)

func main() {
    // 创建格式化器
    formatter := sqlformatter.NewFormatter()
    
    // 配置选项
    formatter.IndentSize = 2      // 缩进空格数
    formatter.KeywordUpper = true // 关键字大写
    
    // 格式化SQL
    sql := "select u.id, u.name from users u where u.age > 25"
    formatted, err := formatter.Format(sql)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(formatted)
}
```

### CLI命令行工具

#### 基本用法

```bash
# 直接输入SQL语句
./sqlformatter -sql "select * from users where id = 1"

# 从文件读取
./sqlformatter -input input.sql

# 输出到文件
./sqlformatter -input input.sql -output formatted.sql

# 使用管道
echo "select * from users" | ./sqlformatter
```

#### 命令行选项

```
选项:
  -sql string      要格式化的SQL语句
  -input string    输入SQL文件
  -output string   输出文件
  -indent int      缩进空格数 (默认: 2)
  -uppercase       关键字大写 (默认: true)
  -help            显示帮助信息
```

## 示例

### SELECT语句

**输入:**
```sql
select u.id, p.product_name, u.name from users u join products p on u.id = p.user_id where u.age > 25 and p.category = 'electronics' group by u.id order by p.price desc limit 10
```

**输出:**
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

### INSERT语句

**输入:**
```sql
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')
```

**输出:**
```sql
INSERT INTO users
  (name, email)
VALUES
  ('Alice', 'alice@example.com')
```

### UPDATE语句

**输入:**
```sql
UPDATE users SET name = 'John', email = 'john@example.com' WHERE id = 1
```

**输出:**
```sql
UPDATE users
SET
  name = 'John',
  email = 'john@example.com'
WHERE
  id = 1
```

### DELETE语句

**输入:**
```sql
DELETE FROM users WHERE age < 18
```

**输出:**
```sql
DELETE FROM users
WHERE
  age < 18
```

## 项目结构

```
sql-formatter/
├── formatter.go        # 核心格式化逻辑
├── cmd/
│   └── main.go        # CLI工具
├── example/
│   └── main.go        # 使用示例
├── formatter_test.go  # 单元测试
├── go.mod             # Go模块文件
├── README.md          # 项目说明（英文）
├── README_zh.md       # 项目说明（中文）
├── test.sql           # 测试SQL文件
└── formatted.sql      # 格式化输出示例
```

## 运行示例

```bash
# 运行库使用示例
go run example/main.go

# 运行CLI工具
go run cmd/main.go -help
go run cmd/main.go -sql "select * from users"
go run cmd/main.go -input test.sql

# 运行测试
go test -v

# 运行基准测试
go test -bench=.
```

## 开发计划

- [ ] 支持更多SQL语句类型（CREATE TABLE、ALTER TABLE等）
- [ ] 添加更多格式化选项
- [ ] 支持不同数据库方言
- [ ] 添加语法高亮
- [ ] 集成Vitess或TiDB解析器以提供更强大的解析能力

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License 