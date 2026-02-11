# Glox

Glox is an implementation of the lox compiler in Go, and is based on the clox compiler written in C (clox).  Lox is the creation of Robert Nystrom (craftinginterpreters.com).

Many comments in the code are taken directly from the book Crafting Interpreters.

## Lox Grammar
  - https://craftinginterpreters.com/appendix-i.html

## Major differences from clox
  - Go has a gc, hence all memory related stuffs are gone.
  - Go has no pointer arithmetic, hence need to use index.
  - Go does not have enum, hence use const and type.
  - Go does not have inline, so define anonymous function in function scope.
  - Go does not have explicit conditional compilation, so at compiler discretion.
  - Go does not have C-like macro, hence use function instead.
  - Go does not use header files, hence stuffs put in .h will be put in the corresponding .go file instead.
  - Go nil replaces C NULL.
  - Go use Uppercase for export symbols.
  - Although no need to implement the hash map mechanics in the Go implementation, still have to implement the hashmap access api.
  - String interning: we do not implement string interning in the Go implementation.
  - Some words are keywords in Go but not in C.  So have to be named differently:
    - type -> type_
    - string() -> str() // compiler.go
  - glox does not have common.h, so need to scatter the definitions defined in common.h elsewhere.
  - Avoid using pointer in glox because that goes against the philosophy of Go which make it difficult to convert from C to Go.
  - In Go, cannot assign nil to a struct.

## Implementation

Although some mechanics required by C is not necessary in Go, we still choose to implement some of them as a programming exercise as well as to mirror the struct of clox as much as reasonable.

## Pratt Parser

In the C code, use forward declaration to handle a declaration cycle in the Pratt Parser.  In Go, use an init() to fix the invalid initialization cycle:

Package level variable "rules" depends on binary() in initialication, binary() depends on getRules(), which depends on "rules".

## Tagged Union

There are two ways to implement a tagged union in Go.  Either use a struct (higher performance?) or an interface (more space efficient).  I picked interface because the implementation is more interesting.

## Chunk OpCode

In order to add a type to the chunk opcode structure, will need to create a type interface and use type constraint in other functions such as writeByte() and writeBytes(), because the parameters that they take can be an opCode or a data byte (uint8).

## Macros

The READ_BYTE, READ_SHORT, READ_CONSTANT, READ_STRING  defined in vm.run() are implemented as inner functions defined in vm.run().

### BINARY_OP

It is not straightforward to convert the BINARY_OP macro in vm.c to Go becuase the macro takes a macro and an operator as parameters.  

Typically a C macro is converted to a Go function.  A C macro does not have type because it is just text substitution.  But a Go function is statically typed.  It is impossible to convert a C macro which takes an arbitrary function as a parameter.

Also there is no way to pass an operator as a function argument in Go.

BINARY_OP macro is implemented as an ordinary GO function binary_op()

## Import Cycle

There is an import cycle between value.go and object.go.  The cycle is broken by extracting code which depending on both object.go & value.go to a new file objval.go.

## Struct Inheritance

Like tagged union, we use interface to implement struct inheritance in clox.

## Object & Value

An object is a struct with 2 fields representing its type and its object value.  ObjString is an alias to string in glox.

structure Object {
  type // ObjFunction or ObjString
  value
}

A value is a struct with 2 fields representing its type and value.  The value can be a primitive data type, or an object.

structure Value {
  type // bool, nil, number, object
  value
}

Difference from clox: VAL_UNDEFINED ValueType is added in glox, so that the zero value of Value will not show up with type VAL_BOOL.

## String Interning

Since glox use Go built-in map datatype to implement the table api, there is no need to implement string interning because the map datatype already take care of it.

## Constants Pool

If we add the same values twice in the constant pool, it will occupies two slots in the constant pool.  That is duplication exists.  Do we want to de-duplicate it?

## Local Variables

Lox, like many programming languages, stores local variables in the stack.

The locals array in the compiler struct has the exact same layout  as the VM's stack at runtime.  The variable's index in the locals array is the same as its stack slot.  How convenient.

## Testing

Since Lox does not have error handling construct, most of the testings are just running valid Lox scripts (vm_test.go) to see if they crash without checking their output.

## Troubleshooting

The following is the most common bugs:
  - missing statement
  - incorrect equality test
  - confusion in object and value conversion

If I have to do it again, I will not translate the C code literally to Go, but just the concept.  For example, avoid using pointers.

When reading code, always keep in mind whether it is executing at compile time or run time.  Hint: is the code in compiler.go or vm.go :P

## Future Improvement

  - Save the generated bytecodes to a file and later load and run the bytecodes directly without the need of parsing the Lox source code.  However, will need to persist the constants too because the bytecodes will not work without them.

## Native Functions

A programming language implementation reaches out and touches the material world through native functions.

At the language level, Lox is fairly complete—it’s got closures, classes, inheritance, and other fun stuff. One reason it feels like a toy language is because it has almost no native capabilities. We could turn it into a real language by adding a long list of them.

Native functions are different from Lox functions. When they are called, they don’t push a CallFrame, because there’s no bytecode code for that frame to point to. They have no bytecode chunk. Instead, they somehow reference a piece of native C code.

Without something like a foreign function interface, users can’t define their own native functions. That’s our job as VM implementers. Glox defineNative() is the foreign function interface.

## End
