package vm

import (
	"testing"
)

func TestScripts(t *testing.T) {
	InitVM()
	//	s := `var drink = "coffee";
	//
	// var breakfast = "crossiant with " + drink;
	// print breakfast;`
	// s := `var drink = "coffee"; print coffee;`
	 s := `var x = 1; print x;`
	result := Interpret(&s)
	if result != INTERPRET_OK {
		t.Error("Script error")
	}
	FreeVM()
}
