# A Simple Machine Language Simulator

This is a Go implementation of a teaching simulator originally written in Java. This simulator
is used to demonstrate machine programming in computer science courses. The instruction set is very
limited, but that's the challenge now isn't it.

## Building

This program requires Go 1.9+ to compile.

## Running the Simulator

By default, the simulator will load a file named `MachineIn.txt` as the code and write the state
and machine printer to standard out. Both of these can be changed using CLI flags as described below.

### Command Flags

- `-in`: Path to the input code file.
- `-out`: Path to the output file. If this is the text "stdout", output will be printed to standard output instead of a file.
- `-state`: Print the machine state before every instruction execution. (this will be very large for 16-bit programs)
- `-printmem`: Print the initial memory state after loading the code. Not instructions are executed.
- `-compile`: Compile a source file and output it to a binary file. The binary still needs the runtime to
execute. It can be used with the "-in" flag. The loader will automatically load the binary if
it has the ASML header.

## Architecture

This machine emulates a 16-bit CPU with 16-bit memory addresses. The total available memory is 64K.

This machine has 16 registers. Registers are numbered 0-F. Memory cell addresses are from 0x0000-0xFFFF.

A value in memory address 0xFFFF will result in the value's ASCII representation being printed to the machine printer.
The cell's value will then be reset to 0. Because of this, the memory location can't actually hold a value
between instruction loads.

NOTE: The store instructions work on two bytes of memory at a time. Meaning, there's no way to store a single
byte of data into a single memory location. A full 16-bit word must be written to memory. So, to write a value to
the printer address 0xFFFF, you really need to write a full 16-bit value to address 0xFFFE. The lower 8 bits will be
stored to the printer address and written to the printer. See the examples for a demonstration.

## Language

### Instructions

Each instruction is 1 byte followed by 0-3 1 byte arguments. Some instructions will combine arguments to form
a full 16-bit value, in particular for immediate values or memory locations.

Each instruction is written using mnemonics:

```
LOADI %1 0xFFFF
HALT
```

'%' denotes a register. (%0 - %F)

Literals have various forms depending on the base. The prefix '0x' denotes a hexadecimal number, the prefix
'0' is an octal number, anything else is assumed to be a decimal number. Literals are parsed as unsigned values.
To make negative numbers, use the two's compliment of the value in hex.

Single bytes can be used by enclosing them in single quotes. (`'H'`)

Strings can be used on a line by themselves by enclosing a string in double quotes. The containing bytes
are inserted as is an interpreted as raw data.

Raw bytes can be used when the line doesn't start with a comment, label, or instruction.
See the data section examples below.

### Comments

Comments are allowed by starting a line with a semicolon, blank lines are ok as well:

```
; Load FF into register 1
LOADI %1 0xFFFF

; Halt
HALT
```

### Labels

Labels can be used in place of a full byte argument. Labels are defined on their own
line starting with a colon:

```
:data
13 45 0x45 0644
```

Labels can be used anywhere a full byte would be used. The syntax is `~labelName`:

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

## Instruction Set

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XYXY' which is a 16-bit
number. In general, if a destination register is needed, it will be the first operand.

| Opcode | Operands | Description                                                                                                                                 |
|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| 0      |          | NOOP                                                                                                                                        |
| 1+     | RXY      | Load the value in memory address XY into register R.                                                                                        |
| 2      | RXY      | Load the value XY into register R.                                                                                                          |
| 3+     | RXY      | Store the value of register R into memory address XY.                                                                                       |
| 4      | RS       | Move the value of register S into register R.                                                                                               |
| 5      | RST      | Add the values in registers S and T using 2's compliment. The result will be stored in register R.                                          |
| 6      |          | Used internally by compiler.                                                                                                                |
| 7      | RST      | OR the values of registers S and T and store the value in register R.                                                                       |
| 8      | RST      | AND the values of registers S and T and store the value in register R.                                                                      |
| 9      | RST      | XOR the values of registers S and T and store the value in register R.                                                                      |
| A      | RX       | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| B+     | RXY      | Jump to memory address XY if the value in register R equals the value in register 0.                                                        |
| C      |          | Halt execution.                                                                                                                             |
| D      | RS       | Store the value of register S into the memory address stored in register R.                                                                 |
| E      | RS       | Load the value at the memory address stored in register S to register R.                                                                    |
| F      |          | Not implemented.                                                                                                                            |

\+ These instructions take another byte to make the two bytes of a 16-bit address.

## Mnemonics

| Opcode | Mnemonic | Syntax        |
|--------|----------|---------------|
| 0      | NOOP     | NOOP          |
| 1      | LOADA    | LOADA %R XYXY |
| 2      | LOADI    | LOADI %R XYXY |
| 3      | STRA     | STRA %R XYXY  |
| 4      | MOVR     | MOVR %R %S    |
| 5      | ADD      | ADD %R %S %T  |
| 6      |          |               |
| 7      | OR       | OR %R %S %T   |
| 8      | AND      | AND %R %S %T  |
| 9      | XOR      | XOR %R %S %T  |
| A      | ROT      | ROT %R X      |
| B      | JMP      | JMP %R XYXY   |
| C      | HALT     | HALT          |
| D      | STRR     | STRR %R %S    |
| E      | LOADR    | LOADR %R %S   |
| F      |          |               |

## Example

The following example prints the character 'X' to the printer 3 times using a loop:

```
; Register 1 - loop counter
LOADI %1 3

; Register 2 - constant -1
LOADI %2 0xFFFF

; Register 3 - char X
LOADI %3 'X'

:print_x
; Print X
STRA %3 0xFFFE

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
