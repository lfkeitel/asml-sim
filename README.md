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
registers overlaid on registers 2-9. Registers E-F are quad width. Memory cell addresses are from 0x0000-0xFFFF.

```
Register Layout:

|---------------------------------------|
| 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
|---------------------------------------|
|       |   A   |   B   |   C   |   D   |
|---------------------------------------|
|       |       E       |       F       |
|---------------------------------------|
```

A value in memory address 0xFFFF will result in the value's ASCII representation being printed to the machine printer.
The cell's value will then be reset to 0. Because of this, the memory location can't actually hold a value
between instruction loads.

The number of bytes written to memory depends on the length of the source register. Single, double, and quad width
registers will write 1, 2, and 4 bytes respectivly starting at the address in the instruction.

## Language

### Instructions

Each instruction is 1 byte followed by 0-3 1 byte arguments. Some instructions will combine arguments to form
a 16-bit value, in particular for immediate values or memory locations.

Each instruction is written using mnemonics:

```
LOADI %1 0xFF
HALT
```

'%' denotes a register. (%0 - %F)

Literals have various forms depending on the base. The prefix '0x' denotes a hexadecimal number, the prefix
'0' is an octal number, anything else is assumed to be a decimal number. Literals are parsed as unsigned values.
To make negative numbers, use the two's compliment of the value in hex.

Single bytes can be used by enclosing them in single quotes. (`'H'`)

Strings can be used on a line by themselves by enclosing a string in double quotes. The containing bytes
are inserted as is and interpreted as raw data.

Raw bytes can be used when the line doesn't start with a comment, label, or instruction.
See the data section examples below.

### Comments

Comments are allowed by starting a line with a semicolon, blank lines are ok as well:

```
; Load FF into register 1
LOADI %1 0xFF

; Halt
HALT
```

### Labels

Labels are defined on their own line starting with a colon:

```
:data
13 45 0x45 0644
```

Labels can be used anywhere a 1 or 2 byte argument would be used. The syntax is `~labelName`:

```
; Load the value in memory location 'data' (13) to register 1
LOADA %1 ~data

:data
13 42
```

Labels can also have simple math applied:

```
; Load the value in memory location 'data+1' (42) to register 1
LOADA %1 ~data+1

:data
13 42
```

The first and second byte of a two byte address can be retrieved using a special syntax.

To get the first byte: `~^label`

To get the second byte: <code>~\`label</code>

To get the full addresss: `~label`

The special label `~$` references the address of the current instruction. Offsets may be used as normal.

NOTE: If a label address is loaded into a single-width register, only the lower byte of the address
is stored. The higher byte is assumed to be 0x00.

## Instruction Set

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XYXY' which is a 16-bit
number. In general, if a destination register is needed, it will be the first operand.

| Name    | Opcode | Arg1 | Arg2 | Arg3 |
|---------+--------+------+------+------+
| NOOP    |   0    |      |      |      |
| LOADA   |   1    |  %D  |  H   |  L   |
| LOADI   |   2    |  %D  |  H   |  L   |
| STRA    |   3    |  %S  |  H   |  L   |
| MOVR    |   4    |  %D  |  %S  |      |
| ADD     |   5    |  %D  |  %S1 |  %S2 |
| FLAGS   |   6-   |      |      |      |
| OR      |   7    |  %D  |  %S1 |  %S2 |
| AND     |   8    |  %D  |  %S1 |  %S2 |
| XOR     |   9    |  %D  |  %S1 |  %S2 |
| ROT     |   A    |  %D  |  /B  |      |
| JMP     |   B    |  %S  |  H   |  L   |
| HALT    |   C    |      |      |      |
| STRR    |   D    |  %S  |  %D  |      |
| LOADR   |   E    |  %D  |  %s  |      |

Each opcode is one byte. Each arg is one byte.

%D destination register
%S1 source register 1
%S2 source register 2
H higher byte of 2-byte value
L lower byte of 2-byte value
B one byte value
B/ half byte high
/B half byte low
B/B two half byte values

## Instruction Descriptions

| Mnemonic | Description                                                                                                                                 |
|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| NOOP     | NOOP                                                                                                                                        |
| LOADA    | Load the value in memory address XY into register R.                                                                                        |
| LOADI    | Load the value XY into register R.                                                                                                          |
| STRA     | Store the value of register R into memory address XY.                                                                                       |
| MOVR     | Move the value of register S into register R.                                                                                               |
| ADD      | Add the values in registers S and T using 2's compliment. The result will be stored in register R.                                          |
|          | Used internally by compiler.                                                                                                                |
| OR       | OR the values of registers S and T and store the value in register R.                                                                       |
| AND      | AND the values of registers S and T and store the value in register R.                                                                      |
| XOR      | XOR the values of registers S and T and store the value in register R.                                                                      |
| ROT      | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| JMP      | Jump to memory address XY if the value in register R equals the value in register 0.                                                        |
| HALT     | Halt execution.                                                                                                                             |
| STRR     | Store the value of register S into the memory address stored in register R.                                                                 |
| LOADR    | Load the value at the memory address stored in register S to register R.                                                                    |
|          | Not implemented.                                                                                                                            |

## Example

The following example prints the character 'X' to the printer 3 times using a loop:

```
; Register 1 - loop counter
LOADI %1 3

; Register 2 - constant -1
LOADI %2 0xFF

; Register 3 - char X
LOADI %3 'X'

:print_x
; Print X
STRA %3 0xFF

; Add r1 and r2 (r1 - 1), store in r1
ADD %1 %1 %2

; Check if loop counter is 0
JMP %1 ~end

; Unconditional jump to print X
JMP %0 ~print_x

:end
; Exit
HALT
```

## Runtime

A small "runtime" is available by using the directive `@runtime`. The directive
must be on its own line.

The runtime code is available [here](examples/runtime.asml). It just sets up a
few registers with common values, sets up a stack, and return link register.
It also provides two labels `~exit` and `~return`. `~exit` simply points to a
HALT instruction. `~return` makes jumping to a memory location simple. It will
take the value in register `%D`, modify a JMP instruction, and then jump to that
location. ASML doesn't provide an instruction to jump to an address in a register
so the code is edited at runtime with the appropriate address. This is sometimes
referred to as a [Wheeler Jump](https://en.wikipedia.org/wiki/Goto#Wheeler_Jump)
after David Wheeler who first developed the technique.

If a program uses the runtime, the "reserved" registers MUST be used for their
intended purposes. Currently, the only "reserved" register is %D, the link register.
This register is used for returning to a calling function. The other registers
are purely for the programmer's convenience and can be changed or reused as needed.
