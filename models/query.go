package models

type Condition struct {
	Left     string
	Operator string
	Right    string
}

type SelectQuery struct {
	Columns []string
	Table   string
	Where   *Condition
}

type CreateTableQuery struct {
    Table   string
    Columns []string
}

type InsertQuery struct {
    Table   string
    Values  []string
}