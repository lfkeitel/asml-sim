# A Simple Machine Language Simulator

This is a Go implementation of a teaching simulator originally written in Java. This simulator
is used to demonstrate machine programming in computer science courses. The instruction set is very
limited, but that's the challenge now isn't it.

## Building

This program requires Go 1.9+ to compile.

## Running the Simulator

By default, the simulator will load a file named `MachineIn.txt` as the code and write the state
and machine printer to standard out. Both of these can be changed using CLI flags as described below.

`asml [OPTIONS] file...`

### Command Options

- `-out`: Path to the output file. If this is the text "stdout", output will be printed to standard output instead of a file.
- `-state`: Print the machine state before every instruction execution. (this will be very large for 16-bit programs)
- `-printmem`: Print the initial memory state after loading the code. Not instructions are executed.
- `-compile`: Compile a source file and output it to a binary file. The binary still needs the runtime to
execute. The compiled may be used in place of a source file.

## Architecture

This machine emulates a 8-bit CPU with 16-bit memory addresses. The total available memory is 64K.

The machine has 10 "physical" 8-bit registers named 0 through 9. Registers A - D are double width
registers overlaid on registers 2-9. Memory cell addresses are from 0x0000-0xFFFF.

```
Register Layout:

|---------------------------------------|
| 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
|---------------------------------------|
|       |   A   |   B   |   C   |   D   |
|---------------------------------------|
```

A value in memory address 0xFFFF will result in the value's ASCII representation being printed to the machine printer.
The cell's value will then be reset to 0. Because of this, the memory location can't actually hold a value
between instruction loads.

The number of bytes written to memory depends on the length of the source register. Single and double width
registers will write 1 or 2 bytes respectively starting at the address in the instruction.

## Language

### Instructions

Each instruction is 1 byte followed by 0-3 1 byte arguments. Some instructions will combine arguments to form
a 16-bit value, in particular for immediate values or memory locations.

Each instruction is written using mnemonics:

```
LOAD %1 #0xFF
HALT
```

'%' denotes a register. (%0 - %D)

Literals have various forms depending on the base. The prefix '0x' denotes a hexadecimal number, the prefix
'0' is an octal number, the prefix '!' is a binary number, anything else is assumed to be a decimal number.
Literals are parsed as unsigned values. To make negative numbers, use the two's compliment of the value in hex.

Strings can be used either on their own line or in place of an immediate value. When used as an immediate value,
the length of the string must fit in the destination register or operation. For example, loading register 1 with
a string must have a length of 1. Strings are created by enclosing text in double quotes. The containing bytes
are inserted as is and interpreted as raw data.

Raw bytes can be used when the line doesn't start with a comment, label, or instruction.
See the data section examples below.

Immediate values are prefixed with a pound sign: `#0xC000`.

### Comments

Comments are allowed by starting a line with a semicolon, blank lines are ok as well:

```
; Load FF into register 1
LOAD %1 #0xFF

; Halt
HALT
```

### Labels

Labels are defined on their own line starting with a colon:

```
:data
13 45 0x45 0644
```

Labels can be used anywhere a 1 or 2 byte argument would be used:

```
; Load the value in memory location 'data' (13) to register 1
LOAD %1 data

:data
13 42
```

Labels can also have simple math applied:

```
; Load the value in memory location 'data+1' (42) to register 1
LOAD %1 data+1

:data
13 42
```

The special label `$` references the address of the current instruction. Offsets may be used as normal.

NOTE: If a label address is loaded into a single-width register, only the lower byte of the address
is stored.

## Instructions

Instruction information can be found in the [documentation](docs/instructions.md).

## Example

The following example prints the character 'X' to the printer 3 times using a loop:

```
    ; Register 1 - loop counter
    LOAD %1 #3

    ; Register 2 - constant -1
    LOAD %2 #0xFF

    ; Register 3 - char X
    LOAD %3 #"X"

:print_x
    ; Print X
    STR %3 0xFF

    ; Add register 2 to register 1
    ADD %1 %2

    ; Check if loop counter is 0
    JMP %1 end

    ; Unconditional jump to print X
    JMPA print_x

:end
    ; Exit
    HALT
```
