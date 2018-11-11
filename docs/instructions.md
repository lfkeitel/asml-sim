# Instruction Set

## Assembly Table

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XYXY' which is a 16-bit
number. In general, if a destination register is needed, it will be the first operand.

| Name    | Opcode | Arg1 | Arg2 | Arg3 |
|---------|--------|------|------|------|
| NOOP    |  0x00  |      |      |      |
| ADD     |  0x01  |  %D  |  %S1 |  %S2 |
| ADDI    |  0x02  |  %D  |  %S  |  B   |
| AND     |  0x03  |  %D  |  %S1 |  %S2 |
| OR      |  0x04  |  %D  |  %S1 |  %S2 |
| ROT     |  0x05  |  %D  |  B   |      |
| XOR     |  0x06  |  %D  |  %S1 |  %S2 |
| CALLI   |  0x07  |  H   |  L   |      |
| CALLR   |  0x08  |  %S  |      |      |
| RTN     |  0x09  |      |      |      |
| HALT    |  0x0A  |      |      |      |
| JMP     |  0x0B  |  %S  |  H   |  L   |
| JMPA    |  0x0C  |  H   |  L   |      |
| LDSP    |  0x0D  |  H   |  L   |      |
| LDSPI   |  0x0E  |  H   |  L   |      |
| LOADA   |  0x0F  |  %D  |  H   |  L   |
| LOADI   |  0x10  |  %D  |  H   |  L   |
| LOADR   |  0x11  |  %D  |  %S  |      |
| STRA    |  0x12  |  %S  |  H   |  L   |
| STRR    |  0x13  |  %S  |  %D  |      |
| MOVR    |  0x14  |  %D  |  %S  |      |
| POP     |  0x15  |  %S  |      |      |
| PUSH    |  0x16  |  %S  |      |      |

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
| CALLI    | Push the current program counter onto the stack and set the program counter to the address given. |
| CALLR    | Push the current program counter onto the stack and set the program counter to the value of register R. |
| HALT     | Halt execution. |
| JMP      | Jump to memory address XY if the value in register R equals the value in register 0. |
| JMPA     | Jump unconditionally to address. |
| LDSP     | Load the stack pointer with the contents of address. |
| LDSPI    | Load the stack pointer with an immediate value. |
| LOADA    | Load the value in memory address XY into register R. |
| LOADI    | Load the value XY into register R. |
| LOADR    | Load the value at the memory address stored in register S to register R. |
| MOVR     | Move the value of register S into register R. |
| NOOP     | Perform no operation. |
| OR       | OR the values of registers S and T and store the value in register R. |
| POP      | Read a value from the stack and store in register. The stack pointer is incremented the size of the destination register. |
| PUSH     | Store the value in register to the stack. The stack pointer is decremented the size of the source register. |
| ROT      | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| RTN      | Pop 2 bytes off the stack and set the program counter to that address. |
| STRA     | Store the value of register R into memory address XY. |
| STRR     | Store the value of register S into the memory address stored in register R. |
| XOR      | XOR the values of registers S and T and store the value in register R. |
