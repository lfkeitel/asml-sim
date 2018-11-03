# Instruction Set

## Assembly Table

In the table below, the first column is the opcode in hexadecimal. The second column is a syntax
for the operands of an opcode. An 'R', 'S', or 'T' represents a register address. An 'X' or 'Y'
represents a single hexadecimal digit. Many instructions will accept an 'XYXY' which is a 16-bit
number. In general, if a destination register is needed, it will be the first operand.

| Name    | Opcode | Arg1 | Arg2 | Arg3 |
|---------|--------|------|------|------|
| NOOP    |  0x00  |      |      |      |
| LOADA   |  0x01  |  %D  |  H   |  L   |
| LOADI   |  0x02  |  %D  |  H   |  L   |
| STRA    |  0x03  |  %S  |  H   |  L   |
| MOVR    |  0x04  |  %D  |  %S  |      |
| ADD     |  0x05  |  %D  |  %S1 |  %S2 |
| ADDI    |  0x06  |  %D  |  %S  |  B   |
| OR      |  0x07  |  %D  |  %S1 |  %S2 |
| AND     |  0x08  |  %D  |  %S1 |  %S2 |
| XOR     |  0x09  |  %D  |  %S1 |  %S2 |
| ROT     |  0x0A  |  %D  |  B   |      |
| JMP     |  0x0B  |  %S  |  H   |  L   |
| HALT    |  0x0C  |      |      |      |
| STRR    |  0x0D  |  %S  |  %D  |      |
| LOADR   |  0x0E  |  %D  |  %S  |      |
| JMPA    |  0x0F  |  H   |  L   |      |
| LDSP    |  0x10  |  H   |  L   |      |
| LDSPI   |  0x11  |  H   |  L   |      |
| PUSH    |  0x12  |  %S  |      |      |
| POP     |  0x13  |  %S  |      |      |
| CALL    |  0x14  |  H   |  L   |      |
| CALLR   |  0x15  |  %S  |      |      |
| RTN     |  0x16  |      |      |      |

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
| NOOP     | NOOP |
| LOADA    | Load the value in memory address XY into register R. |
| LOADI    | Load the value XY into register R. |
| STRA     | Store the value of register R into memory address XY. |
| MOVR     | Move the value of register S into register R. |
| ADD      | Add the values in registers S and T using 2's compliment. The result will be stored in register R. |
| ADDI     | Add the immediate value B to register S and store the result in register D. |
| OR       | OR the values of registers S and T and store the value in register R. |
| AND      | AND the values of registers S and T and store the value in register R. |
| XOR      | XOR the values of registers S and T and store the value in register R. |
| ROT      | Rotate the bits of register R to the right X places. Bits are shifted to the right and lower order bits are moved to the higher order bits. |
| JMP      | Jump to memory address XY if the value in register R equals the value in register 0. |
| HALT     | Halt execution. |
| STRR     | Store the value of register S into the memory address stored in register R. |
| LOADR    | Load the value at the memory address stored in register S to register R. |
| JMPA     | Jump unconditionally to address. |
| LDSP     | Load the stack pointer with the contents of address. |
| LDSPI    | Load the stack pointer with an immediate value. |
| PUSH     | Store the value in register to the stack. The stack pointer is decremented the size of the source register. |
| POP      | Read a value from the stack and store in register. The stack pointer is incremented the size of the destination register. |
| CALL     | Push the current program counter onto the stack and set the program counter to the address given. |
| CALLR    | Push the current program counter onto the stack and set the program counter to the value of register R. |
| RTN      | Pop 2 bytes off the stack and set the program counter to that address. |
