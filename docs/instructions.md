# Instruction Set

**NOTE:** Please refer to new instruction guide.

## Assembly Table

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XYXY' which is a 16-bit
number. In general, if a destination register is needed, it will be the first operand.

| Name    | Arg1 | Arg2 | Arg3 |
|---------|------|------|------|
| NOOP    |      |      |      |
| ADD     |  %D  |  %S1 |  %S2 |
| ADDI    |  %D  |  %S  |  B   |
| AND     |  %D  |  %S1 |  %S2 |
| OR      |  %D  |  %S1 |  %S2 |
| ROT     |  %D  |  B   |      |
| XOR     |  %D  |  %S1 |  %S2 |
| RTN     |      |      |      |
| HALT    |      |      |      |
| JMP     |  %S  |  H   |  L   |
| JMPA    |  H   |  L   |      |
| POP     |  %S  |      |      |
| PUSH    |  %S  |      |      |

Each opcode is one byte. Each arg is one byte.

%D destination register
%S1 source register 1
%S2 source register 2
H higher byte of 2-byte value
L lower byte of 2-byte value
B one byte value

## Descriptions

| Mnemonic | Description |
|----------|-------------|
| ADD      | Add the values in registers S and T using 2's compliment. The result will be stored in register R. |
| ADDI     | Add the immediate value B to register S and store the result in register D. |
| AND      | AND the values of registers S and T and store the value in register R. |
| HALT     | Halt execution. |
| JMP      | Jump to memory address XY if the value in register R equals the value in register 0. |
| JMPA     | Jump unconditionally to address. |
| NOOP     | Perform no operation. |
| OR       | OR the values of registers S and T and store the value in register R. |
| POP      | Read a value from the stack and store in register. The stack pointer is incremented the size of the destination register. |
| PUSH     | Store the value in register to the stack. The stack pointer is decremented the size of the source register. |
| ROT      | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| RTN      | Pop 2 bytes off the stack and set the program counter to that address. |
| XOR      | XOR the values of registers S and T and store the value in register R. |
