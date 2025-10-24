package vm

import (
	"testing"
)

func TestScripts(t *testing.T) {

	var tests = []struct {
		input string
		want  InterpretResult
	}{
		{`var x = 1; print x;`, INTERPRET_OK},
		{`var drink = "coffee";
	      var breakfast = "crossiant with " + drink;
	      print breakfast;`, INTERPRET_OK},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			InitVM()
			if result := Interpret(&test.input); result != test.want {
				t.Errorf("Script error: %q", test.input)
			}
			FreeVM()
		})
	}
}
