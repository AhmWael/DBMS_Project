# awsql - Lightweight Go DBMS

[![Go Version](https://img.shields.io/badge/go-1.21-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## Overview

**awsql** is a lightweight, file-based DBMS implemented in Go. It supports basic SQL operations, Write-Ahead Logging (WAL), and a TCP server mode, allowing integration into existing systems similar to MySQL or PostgreSQL. The goal is to provide a minimal but functional SQL database for educational and small-scale projects.

---

## Features Supported

| Feature                   | Status                |
| ------------------------- | --------------------- |
| CREATE TABLE              | ✅ Supported           |
| INSERT INTO               | ✅ Supported           |
| SELECT                    | ✅ Supported           |
| WHERE clause              | ✅ Supported (basic)   |
| WAL (Write-Ahead Logging) | ✅ Supported           |
| UPDATE                    | ❌ Not yet             |
| DELETE                    | ❌ Not yet             |
| Indexing / Hashing        | ❌ Not yet             |
| Transactions (ACID)       | ❌ Partially (via WAL) |

---

## Architecture

* **main.go**: TCP server and connection handling
* **parser.go**: SQL parsing and tokenization
* **storage/io.go**: Table file operations (CREATE, INSERT, WAL)
* **storage/wal.go**: Write-Ahead Logging, crash recovery
* **models/query.go**: Query struct definitions
* **data/**: Folder storing table files and WAL logs

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/awsql.git
cd awsql
```

2. Run directly (requires Go installed):

```bash
go run *.go
```

3. Or build a standalone binary:

```bash
go build -o awsql *.go
./awsql
```

The server listens on TCP port `8888` by default.

---

## Usage

Connect using a TCP client (e.g., `nc`, `telnet`) or integrate with your application via TCP.

### Examples:

```sql
-- Create a table with columns
CREATE TABLE users (id, name, age);

-- Insert data into table
INSERT INTO users VALUES (1, 'Ali', 25);
INSERT INTO users VALUES (2, 'Sara', 30);

-- Query data
SELECT * FROM users;
```

Server responses will include parsed query structures, e.g., `PARSED: {Table: users, Columns: [...], Where: nil}`.

---

## WAL (Write-Ahead Logging)

* All operations are first logged to `data/wal.log` before writing to table files.
* On server startup, any uncommitted operations are replayed automatically.
* Ensures crash recovery and minimal data loss.

---

## Roadmap / TODO

* [ ] UPDATE queries
* [ ] DELETE queries
* [ ] Advanced WHERE clauses
* [ ] Indexing & query optimization
* [ ] Joins, Aggregations, GROUP BY
* [ ] Checkpointing WAL
* [ ] User authentication & multi-client handling

---

## Contribution

Contributions are welcome! You can:

* Open issues for bugs or feature requests
* Submit pull requests for improvements
* Suggest enhancements for SQL parsing, WAL, or server features

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
