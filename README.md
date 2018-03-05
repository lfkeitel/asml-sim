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
- `-state`: Print the machine state before every instruction execution.
- `-printmem`: Print the initial memory state after loading the code. Not instructions are executed.
- `-compile`: Compile a source file and output it to a binary file. The binary still needs the runtime to
execute. It can be used with the "-in" flag. The loader will automatically load the binary if
it has the ASML header.

## Architecture

This machine emulates an 8-bit CPU with 8-bit addressing. This limits the available memory to a whopping 256 bytes.
Effectively 255 bytes as the highest address is used for printing to the screen.

This machine has 16 registers. Registers are numbered 0-F in hex so their addresses are 4 bits.
Memory cell addresses are from 00-FF in hex so their addresses are 8 bits.

Storing a value in cell FF will result in the value's ASCII representation being printed to the machine printer.
The cell's value will then be reset to 0. Because of this, memory address FF can't actually hold a value
between instruction loads.

## Language

### Instructions

Each instruction is 16-bits (2-bytes) long. The first 4 bits are the opcode as described below.
The next 3 4-bit sets are the operands for the instruction. Some instructions will combine operands
into a single value, but conceptually they are still separate operands.

```
 0 1 2 3 4 5 6 7  8 9 A B C D E F
|_______|_______||_______|_______|
 opcode   op1      op2     op3
```

Each instruction is written using mnemonics:

```
LOADI %1 FF
HALT
```

'%' denotes a register. (%0 - %F)

Hex numbers use no special syntax. For readability, the syntax `0x##` is also valid.

Single bytes can be used by enclosing them in single quotes. (`'H'`)

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

Labels can be used in place of a full byte argument. Labels are defined on their own
line starting with a colon:

```
:data
13 45
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

## Instruction Set

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XY' which is an 8-bit,
hexadecimal number. '0' is a literal numeral 0. In general, if a destination register is needed,
it will be the first operand.

| Opcode | Operands | Description                                                                                                                                 |
|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| 0      | 000      | NOOP                                                                                                                                        |
| 1      | RXY      | Load the value in memory address XY into register R.                                                                                        |
| 2      | RXY      | Load the value XY into register R.                                                                                                          |
| 3      | RXY      | Store the value of register R into memory address XY.                                                                                       |
| 4      | 0RS      | Move the value of register S into register R.                                                                                               |
| 5      | RST      | Add the values in registers S and T using 2's compliment. The result will be stored in register R.                                          |
| 6      | RST      | Not implemented.                                                                                                                            |
| 7      | RST      | OR the values of registers S and T and store the value in register R.                                                                       |
| 8      | RST      | AND the values of registers S and T and store the value in register R.                                                                      |
| 9      | RST      | XOR the values of registers S and T and store the value in register R.                                                                      |
| A      | R0X      | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| B      | RXY      | Jump to memory address XY if the value in register R equals the value in register 0.                                                        |
| C      | 000      | Halt execution.                                                                                                                             |
| D*     | 0RS      | Store the value of register S into the memory address stored in register R.                                                                 |
| E*     | 0RS      | Load the value at the memory address stored in register S to register R.                                                                    |
| F      |          | Not implemented.                                                                                                                            |

\* These instructions are not part of the original language. They were added to implement a
stack using the value of a register as a memory address.

## Mnemonics

| Opcode | Mnemonic | Syntax       |
|--------|----------|--------------|
| 0      | NOOP     | NOOP         |
| 1      | LOADA    | LOADA %R XY  |
| 2      | LOADI    | LOADI %R XY  |
| 3      | STRA     | STRA %R XY   |
| 4      | MOVR     | MOVR %R %S   |
| 5      | ADD      | ADD %R %S %T |
| 6      |          |              |
| 7      | OR       | OR %R %S %T  |
| 8      | AND      | AND %R %S %T |
| 9      | XOR      | XOR %R %S %T |
| A      | ROT      | ROT %R X     |
| B      | JMP      | JMP %R XY    |
| C      | HALT     | HALT         |
| D      | STRR     | STRR %R %S   |
| E      | LOADR    | LOADR %R %S  |
| F      |          |              |

## Example

The following example prints the character 'X' to the printer 3 times using a loop:

```
; Register 1 - loop counter
LOADI %1 3

; Register 2 - constant -1
LOADI %2 FF

; Register 3 - char X
LOADI %3 'X'

:print_x
; Print X
STRA %3 FF

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
