package table

import (
	"testing"

	"github.com/davidfung/glox/objval"
	"github.com/davidfung/glox/value"
)

func TestTable(t *testing.T) {
	var table Table
	var ok bool
	var newkey bool
	val := value.Value{Type_: value.VAL_NUMBER, Val: float64(1)}

	// table: hello
	InitTable(&table)
	newkey = TableSet(&table, "hello", val)
	if !newkey {
		t.Error("map entry not being created")
	}

	// table: hello
	newkey = TableSet(&table, "hello", val)
	if newkey {
		t.Error("map entry should not be created")
	}

	// table: hello, world
	newkey = TableSet(&table, "world", val)
	if !newkey {
		t.Error("map entry not being created")
	}

	// table: hello, world
	// table2: hello, world
	var table2 Table
	TableAddAll(&table, &table2)
	if &table == &table2 {
		t.Error("table copy error")
	}
	if len(table.entries) != len(table2.entries) {
		t.Error("table copy error")
	}

	// table: hello, world
	// table2: hello, world
	val, ok = TableGet(&table, "hello")
	if objval.AS_NUMBER(val) != float64(1) || !ok {
		t.Error("Table entry retrival error")
	}

	// table: hello, world
	// table2: hello, world
	val, ok = TableGet(&table2, "world")
	if objval.AS_NUMBER(val) != float64(1) || !ok {
		t.Error("Table entry retrival error")
	}

	// table: hello, world
	// table2: hello, world
	_, ok = TableGet(&table2, "not exist")
	if ok {
		t.Error("Table entry retrival error")
	}

	ok = TableDelete(&table2, "hello")
	if !ok || len(table2.entries) != 1 {
		t.Error("table entry deletion error")
	}
}
