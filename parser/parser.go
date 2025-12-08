package parser

import (
	"strings"
	"fmt"
	"awsql/models"
)

// ParseSQL parses a SQL query string and returns the corresponding query struct
func ParseSQL(query string) (interface{}, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	// Make everything uppercase for keyword matching
	upper := strings.ToUpper(query)

	if strings.HasPrefix(upper, "SELECT") {
		return parseSelect(query)
	} else if strings.HasPrefix(upper, "CREATE TABLE") {
		return parseCreateTable(query)
	} else if strings.HasPrefix(upper, "INSERT INTO") {
		return parseInsert(query)
	}

	return nil, nil
}

// parseSelect parses a SELECT query string into a SelectQuery struct
func parseSelect(query string) (models.SelectQuery, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	// Split by spaces first
	parts := strings.Fields(query)

	// Must start with SELECT
	if len(parts) < 4 || strings.ToUpper(parts[0]) != "SELECT" {
		return models.SelectQuery{}, fmt.Errorf("invalid SELECT syntax")
	}

	// Find FROM index
	fromIdx := -1
	for i, tok := range parts {
		if strings.ToUpper(tok) == "FROM" {
			fromIdx = i
			break
		}
	}

	if fromIdx == -1 || fromIdx == 1 {
		return models.SelectQuery{}, fmt.Errorf("missing FROM clause")
	}

	// Columns are everything between SELECT and FROM
	columnsStr := strings.Join(parts[1:fromIdx], " ")
	columns := strings.Split(columnsStr, ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	// Table name is next token after FROM
	table := parts[fromIdx+1]

	// Check if WHERE clause exists
	var cond *models.Condition
	whereIdx := -1
	for i, tok := range parts {
		if strings.ToUpper(tok) == "WHERE" {
			whereIdx = i
			break
		}
	}

	if whereIdx != -1 {
		// Simple condition: left operator right
		if len(parts) < whereIdx+4 {
			return models.SelectQuery{}, fmt.Errorf("invalid WHERE clause")
		}
		cond = &models.Condition{
			Left:     parts[whereIdx+1],
			Operator: parts[whereIdx+2],
			Right:    parts[whereIdx+3],
		}
	}

	return models.SelectQuery{
		Columns: columns,
		Table:   table,
		Where:   cond,
	}, nil
}

// parseCreateTable parses a CREATE TABLE query string into a CreateTableQuery struct
func parseCreateTable(query string) (models.CreateTableQuery, error) {
    // Remove CREATE TABLE
    q := strings.TrimSpace(query[len("CREATE TABLE"):])

    // Split table name and columns
    openIdx := strings.Index(q, "(")
    closeIdx := strings.Index(q, ")")
    if openIdx == -1 || closeIdx == -1 || closeIdx < openIdx {
        return models.CreateTableQuery{}, fmt.Errorf("invalid CREATE TABLE syntax")
    }

    table := strings.TrimSpace(q[:openIdx])
    colsStr := q[openIdx+1 : closeIdx]

    cols := strings.Split(colsStr, ",")
    for i := range cols {
        cols[i] = strings.TrimSpace(cols[i])
    }

    return models.CreateTableQuery{
        Table:   table,
        Columns: cols,
    }, nil
}

// parseInsert parses an INSERT INTO query string into an InsertQuery struct
func parseInsert(query string) (models.InsertQuery, error) {
    q := strings.TrimSpace(query[len("INSERT INTO"):])
    // Split table and values
    parts := strings.SplitN(q, "VALUES", 2)
    if len(parts) != 2 {
        return models.InsertQuery{}, fmt.Errorf("invalid INSERT syntax")
    }

    table := strings.TrimSpace(parts[0])

    valsStr := strings.TrimSpace(parts[1])
    valsStr = strings.TrimPrefix(valsStr, "(")
    valsStr = strings.TrimSuffix(valsStr, ")")

    vals := strings.Split(valsStr, ",")
    for i := range vals {
        vals[i] = strings.TrimSpace(vals[i])
    }

    return models.InsertQuery{
        Table:  table,
        Values: vals,
    }, nil
}