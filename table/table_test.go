package table

import (
	"testing"

	"github.com/davidfung/glox/objval"
	"github.com/davidfung/glox/value"
)

func TestTable(t *testing.T) {
	var table Table
	val := value.Value{Type_: value.VAL_NUMBER, Val: float64(1)}

	InitTable(&table)
	created := TableSet(&table, "hello", val)
	if !created {
		t.Error("map entry not being created")
	}

	created = TableSet(&table, "hello", val)
	if created {
		t.Error("map entry should not be created")
	}

	created = TableSet(&table, "world", val)
	if !created {
		t.Error("map entry not being created")
	}

	var table2 Table
	TableAddAll(&table, &table2)
	if &table == &table2 {
		t.Error("table copy error")
	}
	if len(table.entries) != len(table2.entries) {
		t.Error("table copy error")
	}

	val, found := TableGet(&table, "hello")
	if objval.AS_NUMBER(val) != float64(1) || !found {
		t.Error("Table entry retrival error")
	}

	val, found = TableGet(&table2, "world")
	if objval.AS_NUMBER(val) != float64(1) || !found {
		t.Error("Table entry retrival error")
	}

	_, found = TableGet(&table2, "not exist")
	if found {
		t.Error("Table entry retrival error")
	}

	found = TableDelete(&table2, "hello")
	if !found || len(table2.entries) != 1 {
		t.Error("table entry deletion error")
	}
}
