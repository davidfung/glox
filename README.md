# Glox

Glox is an implementation of the lox compiler in Go, and is based on the clox compiler written in C (clox).  Lox is the creation of Robert Nystrom (craftinginterpreters.com).

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
  - Skip implementing hash map and string interning in glox.

## Pratt Parser

In the C code, use forward declaration to handle a declaration cycle in the Pratt Parser.  In Go, use an init() to fix the invalid initialization cycle:

Package level variable "rules" depends on binary() in initialication, binary() depends on getRules(), which depends on "rules".

## Tagged Union

There are two ways to implement a tagged union in Go.  Either use a struct (higher performance?) or an interface (more space efficient).  I picked interface because the implementation is more interesting.

## Chunk OpCode

In order to add a type to the chunk opcode structure, will need to create a type interface and use type constraint in other functions such as writeByte() and writeBytes(), because the parameters that they take can be an opCode or a data byte (uint8).

## BINARY_OP Macro

It is not straightforward to convert the BINARY_OP macro in vm.c to Go becuase the macro takes a macro and an operator as parameters.  

Typically a C macro is converted to a Go function.  A C macro does not have type because it is just text substitution.  But a Go function is statically typed.  It is impossible to convert a C macro which takes an arbitrary function as a parameter.

Also there is no way to pass an operator as a function argument in Go.

## Import Cycle

There is an import cycle between value.go and object.go.  The cycle is broken by moving code from object.go to value.go.

## Struct Inheritance

Like tagged union, we use interface to implement struct inheritance in clox.

## Keywords

Some words are keywords in Go but not in C.  So have to be named differently:

  - type -> type_
  - string() -> str()

## ObjString

  clox ObjString translated to Go regular string, because there is no need to keep any housekeeping info for the string.

## End
