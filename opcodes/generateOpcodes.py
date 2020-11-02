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
(cpu.setFlag("X",0))
(cpu.zflag())
(cpu.hflag())
(cpu.cflag(,,opcodeMnemonic))

OPCODE()
IF bytes == 2 then pc++ elif bytes == 2 then pc += 2
cycle = cycles
instruction = instruction
}
}


"""
import csv

#Initialise file
with open("baseInstructions.go","w+") as file:
    pass

starterLines = [
    "package main\n"
    "func (cpu *gameboyCPU) decodeAndExecute(opcode uint8)(int,string){",
    "var cycles int",
    "var instruction string",
    "switch opcode{"
]

endLines = ["}","return cycles,instruction","}"]

def writeLine(line):
    with open("baseInstructions.go","a") as file:
        file.write(line+"\n")

for line in starterLines:
    writeLine(line)

seperationCounter = 0 #Every 16 entries have an extra newline
seperationTag = 0 #Comments the hex row

#Read each line of the csv file and construct the go file
with open('GB_Opcodes.csv','r') as file:
    reader = csv.reader(file)
    next(reader,None) #skip header
    for row in reader:
        instructionSet = row[0]
        instruction = row[1]
        if instructionSet == "Base" and instruction != "NOT USED":
            mnemonic = row[2]
            byteSize,cycles = row[5].split('  ')
            z,n,h,c = row[6].split(' ')
            hexvalue = row[7]

            if seperationCounter % 16 == 0:
                 writeLine('\n')
                 writeLine(f'//-----------Row {seperationTag+1}: 0x{format(seperationTag,"01X")}-----------')
                 seperationTag += 1
            
            writeLine(f'case 0x{hexvalue}:')
            writeLine(f'\t\t//0x{hexvalue}: {instruction}\n')
            if z != '-':
                if z == '1':
                    writeLine(f'cpu.setFlag("Z",1)')
                elif z == '0':
                    writeLine(f'cpu.setFlag("Z",0)')
                else:
                    writeLine(f'\t\t//cpu.zflag()')
            
            if n != '-':
                if n == '1':
                    writeLine(f'cpu.setFlag("N",1)')
                elif n == '0':
                    writeLine(f'cpu.setFlag("N",0)')

            if h != '-':
                if h == '1':
                    writeLine(f'cpu.setFlag("H",1)')
                elif z == '0':
                    writeLine(f'cpu.setFlag("H",0)')
                else:
                    writeLine(f'\t\t//cpu.hflag( , ,"{mnemonic}")')

            if c != '-':
                if c == '1':
                    writeLine(f'cpu.setFlag("C",1)')
                elif c == '0':
                    writeLine(f'cpu.setFlag("C",0)')
                else:
                    writeLine(f'\t\t//cpu.cflag( , ,"{mnemonic}")')

            writeLine(f'//{mnemonic}()')
            if int(byteSize) == 2:
                writeLine(f'cpu.PC++')
            elif int(byteSize) == 3:
                writeLine(f'cpu.PC += 2')
            
            if "/" not in cycles:
                writeLine(f'cycles = {cycles}')
            else:
                writeLine(f'cycles = {cycles.split("/")[0]} // {cycles.split("/")[1]}')
            writeLine(f'instruction = "{instruction}"\n')

        seperationCounter += 1


for line in endLines:
    writeLine(line)