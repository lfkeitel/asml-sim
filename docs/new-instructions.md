# Instructions

## Modes

Most instructions have different modes depending on their arguments.

### Immediate Mode

Immediate mode is used when a value directly follows the opcode in memory. This
can be used to load and use constant values.

- `LOAD %D #0x1234` - Load the immediate value 0x1234 to register D.

### Address Mode

Address mode is used when interacting with values stored in memory. For example
to load or store a value from a memory location.

- `LOAD %D 0xC000` - Load the value at address 0xC000 to register D.
- `STR %D 0xC000` - Store the value in register D to memory location 0xC000.

### Register Mode

Register mode is used to replace an immediate value or address with a value stored
in a register. Semantics may differ between instructions where some may used the
value as an immediate value and others may use it as an address.

- `LOAD %D %A` - Load the value at the address stored in register A to register D.
- `STR %D %A` - Store the value in register D at the address stored in register A.

## LOAD

Load a value into a register.

### Modes

- Immediate
- Address
- Register

### Examples

- `LOAD %D #0x1234`
- `LOAD %D 0xC000`
- `LOAD %D %A`

## STR

Store a value into memory.

### Modes

- Address
- Register

### Examples

- `STR %D 0xC000`
- `STR %D %A`

## XFER

Transfer data between registers.

### Modes

- Register

### Examples

- `XFER %A %D` - Move the data from register D to register A

## CALL

Call a subroutine by setting the program counter to an address. The current
stack pointer will be pushed onto the stack.

### Modes

- Address
- Register

### Examples

- `CALL 0xC000`
- `CALL %A`
- `CALL sub_label`

## LDSP

Load a value into a the stack pointer.

### Modes

- Immediate
- Address
- Register

### Examples

- `LDSP %D #0x1234`
- `LDSP %D 0xC000`
- `LDSP %D %A`
