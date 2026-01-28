package vm

import (
	"testing"
)

type tests struct {
	input string
	want  InterpretResult
}

func TestScripts(t *testing.T) {
	tests := initTestTable()
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

func initTestTable() []tests {
	var tests = []tests{
		{`
		fun one() {
			return 1;
		}
		var a = one();
		print a;
		`, INTERPRET_OK},
		{`
		fun areWeHavingItYet() {
            print "Yes we are!";
        }
		print areWeHavingItYet;
		`, INTERPRET_OK},
		{`
		var x = 1; print x;
		`, INTERPRET_OK},
		{`
		var drink = "coffee";
	    var breakfast = "crossiant with " + drink;
	    print breakfast;
		`, INTERPRET_OK},
		{`
		var a = 1; var b = 2; var c = 3; var d = 4;
		print a * b = c + d;
		`, INTERPRET_COMPILE_ERROR},
		{`
		var a = a
		`, INTERPRET_COMPILE_ERROR},
		{`
		print (1==2);
		`, INTERPRET_OK},
		{`
		if (1 == 2) {
            print "IF BLOCK";
		} else {
            print "ELSE BLOCK";
		}
		print "CONTINUE BLOCK";
		`, INTERPRET_OK},
		{`
		if (true and false) {} else {}
		`, INTERPRET_OK},
		{`
		if (true or false) {} else {}
		`, INTERPRET_OK},
		{`
        for (var i=1;i<=3;i=i+1) {
        	print i;
		}
		`, INTERPRET_OK},
	}
	return tests
}
