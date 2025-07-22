# Glox

Glox is an implementation of the lox compiler in Go, and is based on the clox compiler written in C (clox).  Lox is the creation of Robert Nystrom (craftinginterpreters.com).

## Major differences from clox
  - Go has a gc, hence all memory related stuffs are gone.
  - Go has no pointer arithmetic, hence need to use index.
  - Go does not have enum, hence use const and type.
  - Go does not have inline, so define anonymous function in function scope.
  - Go does not have explicit conditional compilation, so at compiler discretion.
  - Go does not have C-like macro, hence use function instead

## Pratt Parser

In the C code, use forward declaration to handle a declaration cycle in the Pratt Parser.  In Go, use an init() to fix the invalid initialization cycle:

Package level variable "rules" depends on binary() in initialication, binary() depends on getRules(), which depends on "rules".

## Tagged Union

There are two ways to implement a tagged union in Go.  Either use a struct (higher performance?) or an interface (more space efficient).  I picked interface because the implementation is more interesting.

## Chunk OpCode

In order to add a type to the chunk opcode structure, will need to create a type interface and use
type constraint in other functions such as writeByte() and writeBytes(), because the parameters that
they take can be an opCode or a data byte (uint8)

## End
