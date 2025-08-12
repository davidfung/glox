package table

type Entry string

type Table struct {
	count    int // may not need this
	capacity int // may not need this
	entries  []Entry
}

func initTable(table *Table) {
	table.count = 0     // no need, autoinit to zero value
	table.capacity = 0  // no need, autoinit to zero value
	table.entries = nil // no need, autoinit to zero value
}

func freeTable(table *Table) {
	initTable(table)
}
