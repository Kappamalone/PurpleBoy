"""
Generates the decode and execute go code to reduce time and errors writing
like 500 instructions by hand. It reads the needed parameters from the 
very nicely formatted excel cpu opcode table.

The general write structure:

func (cpu *gameboyCPU) decodeAndExecute(opcode){
var cycles int
var instruction int
switch opcode{
case 0xHEX:
IF Z -> zflag else toggle Z
IF N -> toggle N
IF H -> hflag else toggle H
IF C -> cflag else toggle C

OPCODE()
IF bytes == 2 then pc++ elif bytes == 2 then pc += 2
cycle = cycles
instruction = instruction
}
}


"""