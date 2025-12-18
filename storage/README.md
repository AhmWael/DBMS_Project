# Storage Layer - awsql

This folder contains the core file and page management logic for **awsql**. It explains how tables are stored on disk, the page structure, and how rows are inserted and read.

---

## 1. File Structure

* Each table is stored as a single file in `data/` with the `.tbl` extension:

```
data/
 ├─ users.tbl       # Table "users"
 └─ wal.log         # Write-Ahead Log (WAL)
```

* Table files consist of **fixed-size pages** (default 8 KB each), similar to PostgreSQL.

---

## 2. Page Layout

Each page has a small header and a data section:

| Offset | Size    | Purpose                    |
| ------ | ------- | -------------------------- |
| 0      | 8 bytes | `NextPageID` (linked list) |
| 8      | 2 bytes | `UsedBytes` in page        |
| 10     | rest    | Row data (variable length) |

* `NextPageID` points to the next page if the table has multiple pages.
* `UsedBytes` tracks how much of the page is filled.
* Rows are stored as: `[length (2 bytes)] + [row bytes]`.

---

## 3. Table Header (Page 0)

* The **first page** (page 0) stores the table header:

```text
| NumColumns (2 bytes) | Column Names | NumPages (2 bytes) | Data Page IDs |
```

* Column names are stored with a 2-byte length prefix.
* Data pages IDs list all allocated pages for this table.
* This allows multipage tables with linked pages.

---

## 4. Row Storage & Insertion

* Rows are stored as comma-separated values: `"1,'Ali',25"`.

* When inserting:

  1. Check the last page for available space.
  2. If the page is full, **allocate a new page**, link it via `NextPageID`.
  3. Update the table header with the new page ID.

* Maximum row size: slightly less than page size minus header (about 8 KB - 10 bytes).

---

## 5. Reading Rows

* To read a table:

  1. Load the header (page 0) to get column info and page IDs.
  2. Read each data page in order.
  3. Follow `NextPageID` for multipage tables.
  4. Parse each row using the 2-byte length prefix.

---

## 6. Multipage Support

* Each table can grow dynamically across multiple pages.
* Pages are linked as a singly-linked list using `NextPageID`.
* The header tracks all pages for easier traversal.

---

## 7. Summary

* **Page-based storage** allows fixed-size reads/writes.
* **Header page** stores metadata and column information.
* **Linked data pages** support multi-page tables.
* Row insertion automatically handles page allocation.

