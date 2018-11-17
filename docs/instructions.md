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

Register mode is when a register is used as the source of an instruction.
Semantics may differ between instructions where some may used the
value as an immediate value and others may use it as an address.

- `LOAD %D %A` - Load the value at the address stored in register A to register D.
- `STR %D %A` - Store the value in register D at the address stored in register A.

### Inherent Mode

Inherent mode is when an instruction either has no arguments or its arguments
are embedded in the instruction itself.

- `RTN`
- `HALT`

## RMB

RMB is not a real instruction. It simply reserves the specified number of bytes
in memory.

### Examples

- `RMB 10` - Reserve 10 bytes of memory.

## ORG

ORG is not a real instruction. It sets the origin address of the following code.

### Examples

- `ORG 0xC000` - The following code will start at address 0xC000.

## FCB

FCB is not a real instruction. Its stores literal data into memory. Each piece
of data must be a single byte or a string. A string will be written to memory
as individual bytes.

### Examples

- `FCB 0xCC` - Write the value `0xCC` as starting data to memory.

The following example stores the data starting at location 0x1000 (0x1000-0x1002)
and gives it the label data:

```
    ORG 0x1000
:data
    FCB 0x01, 0x25, 0x42
```

It can also be a string:

```
    ORG 0x1000
:errmsg
    FCB "An error as occured"
```

## FDB

FDB is not a real instruction. Its stores literal data into memory. Each piece
of data is written as a 16 bit number.

### Examples

- `FDB 0x1234` - Write the value `0x1234` as starting data to memory.

The following example stores the data starting at location 0x1000 (0x1000-0x1001)
and gives it the label data:

```
    ORG 0x1000
:data
    FDB 0x1234
```

It's very useful for storing addresses to subroutines:

```
    ORG 0x2000
:main
    ; Do work

    ORG 0xFFFE  ; Reset address
    FDB main    ; Start execution at the address labeled main
```

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

## ADD

Add a value to a register.

### Modes

- Immediate
- Address
- Register

### Examples

- `ADD %D #0x1234`
- `ADD %D 0xC000`
- `ADD %D %A`

## PUSH

Push a register value onto the stack.

### Modes

- Register

### Examples

- `PUSH %A`
- `PUSH %2`

## POP

Pop a value off the stack and store it into a register.

### Modes

- Register

### Examples

- `POP %A`
- `POP %2`

## NOOP

Do nothing.

### Modes

- Inherent

## RTN

Return from a subroutine call. The new program counter is popped off the stack
and execution is resumed.

### Modes

- Inherent

## HALT

Stop all execution.

### Modes

- Inherent

## JMP

Jump to an address if the source register is equal to the value of register 0.

### Modes

- Mixed

### Examples

- `JMP %1 end` - If the value in register 1 equals the value in register 0, jump to
the label "end".

## JMPA

Always jump to an address.

### Modes

- Address

### Examples

- `JMPA bg_loop`

## OR

Or a value to a register.

### Modes

- Immediate
- Address
- Register

### Examples

- `OR %D #0x1234`
- `OR %D 0xC000`
- `OR %D %A`

## AND

And a value to a register.

### Modes

- Immediate
- Address
- Register

### Examples

- `AND %D #0x1234`
- `AND %D 0xC000`
- `AND %D %A`

## XOR

Exclusive or a value to a register.

### Modes

- Immediate
- Address
- Register

### Examples

- `XOR %D #0x1234`
- `XOR %D 0xC000`
- `XOR %D %A`

## ROTR

Rotate the value of a register right.

### Modes

- Immediate

### Examples

- `ROTR %A #4`
- `ROTR %2 #2`

## ROTL

Rotate the value of a register left.

### Modes

- Immediate

### Examples

- `ROTL %A #4`
- `ROTL %2 #2`
