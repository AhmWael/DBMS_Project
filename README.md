# awsql - Lightweight Go DBMS

[![Go Version](https://img.shields.io/badge/go-1.21-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## Overview

**awsql** is a lightweight, educational database management system (DBMS) written in Go.
It implements a minimal SQL engine, a TCP-based server, Write-Ahead Logging (WAL), and a heap-based storage engine with fixed-size paging, inspired by PostgreSQL’s internal architecture.

The primary goal of awsql is to learn and demonstrate database internals, not to compete with production systems.

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

```
DBMS_Project/
├── config/
│   └── version.go           # Server name and version
│
├── models/
│   └── query.go             # Query and condition structs
│
├── parser/
│   └── parser.go            # SQL parsing logic
│
├── storage/
│   ├── io.go                # High-level table I/O
│   │
│   ├── page/
│   │   ├── constants.go     # Page size (8KB) and constants
│   │   ├── manager.go      # Page read/write manager
│   │   └── page.go         # Page structure and helpers
│   │
│   └── table/
│       ├── header.go        # Table metadata and header layout
│       └── heap.go          # Heap table implementation
│
├── wal/
│   └── wal.go               # Write-Ahead Logging
│
├── data/
│   ├── *.tbl                # Heap table files (binary)
│   └── wal.log              # WAL file
│
├── main.go                  # TCP server and query execution
├── go.mod
└── README.md
```

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/AhmWael/DBMS_Project.git
cd DBMS_Project
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
