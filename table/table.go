package table

import (
	"github.com/davidfung/glox/object"
	"github.com/davidfung/glox/value"
)

type Table struct {
	entries map[object.ObjString]value.Value
}

func InitTable(table *Table) {
	table.entries = make(map[object.ObjString]value.Value)
}

func FreeTable(table *Table) {
	InitTable(table)
}

// Return the value and ok=true if found,
// otherwise return the zero value of Value and ok=false
func TableGet(table *Table, key object.ObjString) (val value.Value, ok bool) {
	val, ok = table.entries[key]
	return val, ok
}

// This function adds the given key/value pair to the given hash table.
// If an entry for that key is already present, the new value overwrites
// the old value. The function returns true if a new entry was added.
func TableSet(table *Table, key object.ObjString, val value.Value) (newkey bool) {
	_, ok := table.entries[key]
	table.entries[key] = val
	newkey = !ok
	return
}

// Delete a map entry.  Return true if an entry is found and deleted.
// Return false if an entry is not found.
func TableDelete(table *Table, key object.ObjString) bool {
	_, found := table.entries[key]
	delete(table.entries, key)
	return found
}

// A helper function to copy all of the entries of one table into another.
func TableAddAll(from *Table, to *Table) {
	to.entries = from.entries
}
