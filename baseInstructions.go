package main
func (cpu *gameboyCPU) decodeAndExecute(opcode){
var cycles int
var instruction string
switch opcode{


//Row 1: 0x00
case 0x00:
		//0x00: NOP

//NOP()
cycles = 4
instruction = "NOP"

case 0x01:
		//0x01: LD BC,d16

//LD()
cycles = 12
instruction = "LD BC,d16"

case 0x02:
		//0x02: LD (BC),A

//LD()
cycles = 8
instruction = "LD (BC),A"

case 0x03:
		//0x03: INC BC

//INC()
cycles = 8
instruction = "INC BC"

case 0x04:
		//0x04: INC B

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC B"

case 0x05:
		//0x05: DEC B

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC B"

case 0x06:
		//0x06: LD B,d8

//LD()
cycles = 8
instruction = "LD B,d8"

case 0x07:
		//0x07: RLCA

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"RLCA")
//RLCA()
cycles = 4
instruction = "RLCA"

case 0x08:
		//0x08: LD (a16),SP

//LD()
cycles = 20
instruction = "LD (a16),SP"

case 0x09:
		//0x09: ADD HL,BC

cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD HL,BC"

case 0x0A:
		//0x0A: LD A,(BC)

//LD()
cycles = 8
instruction = "LD A,(BC)"

case 0x0B:
		//0x0B: DEC BC

//DEC()
cycles = 8
instruction = "DEC BC"

case 0x0C:
		//0x0C: INC C

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC C"

case 0x0D:
		//0x0D: DEC C

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC C"

case 0x0E:
		//0x0E: LD C,d8

//LD()
cycles = 8
instruction = "LD C,d8"

case 0x0F:
		//0x0F: RRCA

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"RRCA")
//RRCA()
cycles = 4
instruction = "RRCA"



//Row 2: 0x01
case 0x10:
		//0x10: STOP 0

//STOP()
cycles = 4
instruction = "STOP 0"

case 0x11:
		//0x11: LD DE,d16

//LD()
cycles = 12
instruction = "LD DE,d16"

case 0x12:
		//0x12: LD (DE),A

//LD()
cycles = 8
instruction = "LD (DE),A"

case 0x13:
		//0x13: INC DE

//INC()
cycles = 8
instruction = "INC DE"

case 0x14:
		//0x14: INC D

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC D"

case 0x15:
		//0x15: DEC D

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC D"

case 0x16:
		//0x16: LD D,d8

//LD()
cycles = 8
instruction = "LD D,d8"

case 0x17:
		//0x17: RLA

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"RLA")
//RLA()
cycles = 4
instruction = "RLA"

case 0x18:
		//0x18: JR r8

//JR()
cycles = 12
instruction = "JR r8"

case 0x19:
		//0x19: ADD HL,DE

cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD HL,DE"

case 0x1A:
		//0x1A: LD A,(DE)

//LD()
cycles = 8
instruction = "LD A,(DE)"

case 0x1B:
		//0x1B: DEC DE

//DEC()
cycles = 8
instruction = "DEC DE"

case 0x1C:
		//0x1C: INC E

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC E"

case 0x1D:
		//0x1D: DEC E

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC E"

case 0x1E:
		//0x1E: LD E,d8

//LD()
cycles = 8
instruction = "LD E,d8"

case 0x1F:
		//0x1F: RRA

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"RRA")
//RRA()
cycles = 4
instruction = "RRA"



//Row 3: 0x02
case 0x20:
		//0x20: JR NZ,r8

//JR()
cycles = 12 // 8
instruction = "JR NZ,r8"

case 0x21:
		//0x21: LD HL,d16

//LD()
cycles = 12
instruction = "LD HL,d16"

case 0x22:
		//0x22: LD (HL+),A

//LD()
cycles = 8
instruction = "LD (HL+),A"

case 0x23:
		//0x23: INC HL

//INC()
cycles = 8
instruction = "INC HL"

case 0x24:
		//0x24: INC H

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC H"

case 0x25:
		//0x25: DEC H

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC H"

case 0x26:
		//0x26: LD H,d8

//LD()
cycles = 8
instruction = "LD H,d8"

case 0x27:
		//0x27: DAA

cpu.zflag()
		//cpu.hflag( , ,"DAA")
		//cpu.cflag( , ,"DAA")
//DAA()
cycles = 4
instruction = "DAA"

case 0x28:
		//0x28: JR Z,r8

//JR()
cycles = 12 // 8
instruction = "JR Z,r8"

case 0x29:
		//0x29: ADD HL,HL

cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD HL,HL"

case 0x2A:
		//0x2A: LD A,(HL+)

//LD()
cycles = 8
instruction = "LD A,(HL+)"

case 0x2B:
		//0x2B: DEC HL

//DEC()
cycles = 8
instruction = "DEC HL"

case 0x2C:
		//0x2C: INC L

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC L"

case 0x2D:
		//0x2D: DEC L

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC L"

case 0x2E:
		//0x2E: LD L,d8

//LD()
cycles = 8
instruction = "LD L,d8"

case 0x2F:
		//0x2F: CPL

cpu.setFlag("N",1)
cpu.setFlag("H",1)
//CPL()
cycles = 4
instruction = "CPL"



//Row 4: 0x03
case 0x30:
		//0x30: JR NC,r8

//JR()
cycles = 12 // 8
instruction = "JR NC,r8"

case 0x31:
		//0x31: LD SP,d16

//LD()
cycles = 12
instruction = "LD SP,d16"

case 0x32:
		//0x32: LD (HL-),A

//LD()
cycles = 8
instruction = "LD (HL-),A"

case 0x33:
		//0x33: INC SP

//INC()
cycles = 8
instruction = "INC SP"

case 0x34:
		//0x34: INC (HL)

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 12
instruction = "INC (HL)"

case 0x35:
		//0x35: DEC (HL)

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 12
instruction = "DEC (HL)"

case 0x36:
		//0x36: LD (HL),d8

//LD()
cycles = 12
instruction = "LD (HL),d8"

case 0x37:
		//0x37: SCF

cpu.setFlag("N",0)
		//cpu.hflag( , ,"SCF")
cpu.setFlag("C",1)
//SCF()
cycles = 4
instruction = "SCF"

case 0x38:
		//0x38: JR C,r8

//JR()
cycles = 12 // 8
instruction = "JR C,r8"

case 0x39:
		//0x39: ADD HL,SP

cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD HL,SP"

case 0x3A:
		//0x3A: LD A,(HL-)

//LD()
cycles = 8
instruction = "LD A,(HL-)"

case 0x3B:
		//0x3B: DEC SP

//DEC()
cycles = 8
instruction = "DEC SP"

case 0x3C:
		//0x3C: INC A

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"INC")
//INC()
cycles = 4
instruction = "INC A"

case 0x3D:
		//0x3D: DEC A

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"DEC")
//DEC()
cycles = 4
instruction = "DEC A"

case 0x3E:
		//0x3E: LD A,d8

//LD()
cycles = 8
instruction = "LD A,d8"

case 0x3F:
		//0x3F: CCF

cpu.setFlag("N",0)
		//cpu.hflag( , ,"CCF")
		//cpu.cflag( , ,"CCF")
//CCF()
cycles = 4
instruction = "CCF"



//Row 5: 0x04
case 0x40:
		//0x40: LD B,B

//LD()
cycles = 4
instruction = "LD B,B"

case 0x41:
		//0x41: LD B,C

//LD()
cycles = 4
instruction = "LD B,C"

case 0x42:
		//0x42: LD B,D

//LD()
cycles = 4
instruction = "LD B,D"

case 0x43:
		//0x43: LD B,E

//LD()
cycles = 4
instruction = "LD B,E"

case 0x44:
		//0x44: LD B,H

//LD()
cycles = 4
instruction = "LD B,H"

case 0x45:
		//0x45: LD B,L

//LD()
cycles = 4
instruction = "LD B,L"

case 0x46:
		//0x46: LD B,(HL)

//LD()
cycles = 8
instruction = "LD B,(HL)"

case 0x47:
		//0x47: LD B,A

//LD()
cycles = 4
instruction = "LD B,A"

case 0x48:
		//0x48: LD C,B

//LD()
cycles = 4
instruction = "LD C,B"

case 0x49:
		//0x49: LD C,C

//LD()
cycles = 4
instruction = "LD C,C"

case 0x4A:
		//0x4A: LD C,D

//LD()
cycles = 4
instruction = "LD C,D"

case 0x4B:
		//0x4B: LD C,E

//LD()
cycles = 4
instruction = "LD C,E"

case 0x4C:
		//0x4C: LD C,H

//LD()
cycles = 4
instruction = "LD C,H"

case 0x4D:
		//0x4D: LD C,L

//LD()
cycles = 4
instruction = "LD C,L"

case 0x4E:
		//0x4E: LD C,(HL)

//LD()
cycles = 8
instruction = "LD C,(HL)"

case 0x4F:
		//0x4F: LD C,A

//LD()
cycles = 4
instruction = "LD C,A"



//Row 6: 0x05
case 0x50:
		//0x50: LD D,B

//LD()
cycles = 4
instruction = "LD D,B"

case 0x51:
		//0x51: LD D,C

//LD()
cycles = 4
instruction = "LD D,C"

case 0x52:
		//0x52: LD D,D

//LD()
cycles = 4
instruction = "LD D,D"

case 0x53:
		//0x53: LD D,E

//LD()
cycles = 4
instruction = "LD D,E"

case 0x54:
		//0x54: LD D,H

//LD()
cycles = 4
instruction = "LD D,H"

case 0x55:
		//0x55: LD D,L

//LD()
cycles = 4
instruction = "LD D,L"

case 0x56:
		//0x56: LD D,(HL)

//LD()
cycles = 8
instruction = "LD D,(HL)"

case 0x57:
		//0x57: LD D,A

//LD()
cycles = 4
instruction = "LD D,A"

case 0x58:
		//0x58: LD E,B

//LD()
cycles = 4
instruction = "LD E,B"

case 0x59:
		//0x59: LD E,C

//LD()
cycles = 4
instruction = "LD E,C"

case 0x5A:
		//0x5A: LD E,D

//LD()
cycles = 4
instruction = "LD E,D"

case 0x5B:
		//0x5B: LD E,E

//LD()
cycles = 4
instruction = "LD E,E"

case 0x5C:
		//0x5C: LD E,H

//LD()
cycles = 4
instruction = "LD E,H"

case 0x5D:
		//0x5D: LD E,L

//LD()
cycles = 4
instruction = "LD E,L"

case 0x5E:
		//0x5E: LD E,(HL)

//LD()
cycles = 8
instruction = "LD E,(HL)"

case 0x5F:
		//0x5F: LD E,A

//LD()
cycles = 4
instruction = "LD E,A"



//Row 7: 0x06
case 0x60:
		//0x60: LD H,B

//LD()
cycles = 4
instruction = "LD H,B"

case 0x61:
		//0x61: LD H,C

//LD()
cycles = 4
instruction = "LD H,C"

case 0x62:
		//0x62: LD H,D

//LD()
cycles = 4
instruction = "LD H,D"

case 0x63:
		//0x63: LD H,E

//LD()
cycles = 4
instruction = "LD H,E"

case 0x64:
		//0x64: LD H,H

//LD()
cycles = 4
instruction = "LD H,H"

case 0x65:
		//0x65: LD H,L

//LD()
cycles = 4
instruction = "LD H,L"

case 0x66:
		//0x66: LD H,(HL)

//LD()
cycles = 8
instruction = "LD H,(HL)"

case 0x67:
		//0x67: LD H,A

//LD()
cycles = 4
instruction = "LD H,A"

case 0x68:
		//0x68: LD L,B

//LD()
cycles = 4
instruction = "LD L,B"

case 0x69:
		//0x69: LD L,C

//LD()
cycles = 4
instruction = "LD L,C"

case 0x6A:
		//0x6A: LD L,D

//LD()
cycles = 4
instruction = "LD L,D"

case 0x6B:
		//0x6B: LD L,E

//LD()
cycles = 4
instruction = "LD L,E"

case 0x6C:
		//0x6C: LD L,H

//LD()
cycles = 4
instruction = "LD L,H"

case 0x6D:
		//0x6D: LD L,L

//LD()
cycles = 4
instruction = "LD L,L"

case 0x6E:
		//0x6E: LD L,(HL)

//LD()
cycles = 8
instruction = "LD L,(HL)"

case 0x6F:
		//0x6F: LD L,A

//LD()
cycles = 4
instruction = "LD L,A"



//Row 8: 0x07
case 0x70:
		//0x70: LD (HL),B

//LD()
cycles = 8
instruction = "LD (HL),B"

case 0x71:
		//0x71: LD (HL),C

//LD()
cycles = 8
instruction = "LD (HL),C"

case 0x72:
		//0x72: LD (HL),D

//LD()
cycles = 8
instruction = "LD (HL),D"

case 0x73:
		//0x73: LD (HL),E

//LD()
cycles = 8
instruction = "LD (HL),E"

case 0x74:
		//0x74: LD (HL),H

//LD()
cycles = 8
instruction = "LD (HL),H"

case 0x75:
		//0x75: LD (HL),L

//LD()
cycles = 8
instruction = "LD (HL),L"

case 0x76:
		//0x76: HALT

//HALT()
cycles = 4
instruction = "HALT"

case 0x77:
		//0x77: LD (HL),A

//LD()
cycles = 8
instruction = "LD (HL),A"

case 0x78:
		//0x78: LD A,B

//LD()
cycles = 4
instruction = "LD A,B"

case 0x79:
		//0x79: LD A,C

//LD()
cycles = 4
instruction = "LD A,C"

case 0x7A:
		//0x7A: LD A,D

//LD()
cycles = 4
instruction = "LD A,D"

case 0x7B:
		//0x7B: LD A,E

//LD()
cycles = 4
instruction = "LD A,E"

case 0x7C:
		//0x7C: LD A,H

//LD()
cycles = 4
instruction = "LD A,H"

case 0x7D:
		//0x7D: LD A,L

//LD()
cycles = 4
instruction = "LD A,L"

case 0x7E:
		//0x7E: LD A,(HL)

//LD()
cycles = 8
instruction = "LD A,(HL)"

case 0x7F:
		//0x7F: LD A,A

//LD()
cycles = 4
instruction = "LD A,A"



//Row 9: 0x08
case 0x80:
		//0x80: ADD A,B

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,B"

case 0x81:
		//0x81: ADD A,C

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,C"

case 0x82:
		//0x82: ADD A,D

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,D"

case 0x83:
		//0x83: ADD A,E

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,E"

case 0x84:
		//0x84: ADD A,H

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,H"

case 0x85:
		//0x85: ADD A,L

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,L"

case 0x86:
		//0x86: ADD A,(HL)

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD A,(HL)"

case 0x87:
		//0x87: ADD A,A

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 4
instruction = "ADD A,A"

case 0x88:
		//0x88: ADC A,B

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,B"

case 0x89:
		//0x89: ADC A,C

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,C"

case 0x8A:
		//0x8A: ADC A,D

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,D"

case 0x8B:
		//0x8B: ADC A,E

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,E"

case 0x8C:
		//0x8C: ADC A,H

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,H"

case 0x8D:
		//0x8D: ADC A,L

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,L"

case 0x8E:
		//0x8E: ADC A,(HL)

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 8
instruction = "ADC A,(HL)"

case 0x8F:
		//0x8F: ADC A,A

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 4
instruction = "ADC A,A"



//Row 10: 0x09
case 0x90:
		//0x90: SUB B

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB B"

case 0x91:
		//0x91: SUB C

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB C"

case 0x92:
		//0x92: SUB D

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB D"

case 0x93:
		//0x93: SUB E

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB E"

case 0x94:
		//0x94: SUB H

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB H"

case 0x95:
		//0x95: SUB L

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB L"

case 0x96:
		//0x96: SUB (HL)

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 8
instruction = "SUB (HL)"

case 0x97:
		//0x97: SUB A

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 4
instruction = "SUB A"

case 0x98:
		//0x98: SBC A,B

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,B"

case 0x99:
		//0x99: SBC A,C

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,C"

case 0x9A:
		//0x9A: SBC A,D

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,D"

case 0x9B:
		//0x9B: SBC A,E

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,E"

case 0x9C:
		//0x9C: SBC A,H

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,H"

case 0x9D:
		//0x9D: SBC A,L

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,L"

case 0x9E:
		//0x9E: SBC A,(HL)

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 8
instruction = "SBC A,(HL)"

case 0x9F:
		//0x9F: SBC A,A

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 4
instruction = "SBC A,A"



//Row 11: 0x0A
case 0xA0:
		//0xA0: AND B

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND B"

case 0xA1:
		//0xA1: AND C

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND C"

case 0xA2:
		//0xA2: AND D

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND D"

case 0xA3:
		//0xA3: AND E

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND E"

case 0xA4:
		//0xA4: AND H

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND H"

case 0xA5:
		//0xA5: AND L

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND L"

case 0xA6:
		//0xA6: AND (HL)

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 8
instruction = "AND (HL)"

case 0xA7:
		//0xA7: AND A

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 4
instruction = "AND A"

case 0xA8:
		//0xA8: XOR B

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR B"

case 0xA9:
		//0xA9: XOR C

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR C"

case 0xAA:
		//0xAA: XOR D

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR D"

case 0xAB:
		//0xAB: XOR E

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR E"

case 0xAC:
		//0xAC: XOR H

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR H"

case 0xAD:
		//0xAD: XOR L

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR L"

case 0xAE:
		//0xAE: XOR (HL)

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 8
instruction = "XOR (HL)"

case 0xAF:
		//0xAF: XOR A

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 4
instruction = "XOR A"



//Row 12: 0x0B
case 0xB0:
		//0xB0: OR B

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR B"

case 0xB1:
		//0xB1: OR C

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR C"

case 0xB2:
		//0xB2: OR D

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR D"

case 0xB3:
		//0xB3: OR E

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR E"

case 0xB4:
		//0xB4: OR H

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR H"

case 0xB5:
		//0xB5: OR L

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR L"

case 0xB6:
		//0xB6: OR (HL)

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 8
instruction = "OR (HL)"

case 0xB7:
		//0xB7: OR A

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 4
instruction = "OR A"

case 0xB8:
		//0xB8: CP B

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP B"

case 0xB9:
		//0xB9: CP C

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP C"

case 0xBA:
		//0xBA: CP D

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP D"

case 0xBB:
		//0xBB: CP E

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP E"

case 0xBC:
		//0xBC: CP H

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP H"

case 0xBD:
		//0xBD: CP L

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP L"

case 0xBE:
		//0xBE: CP (HL)

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 8
instruction = "CP (HL)"

case 0xBF:
		//0xBF: CP A

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 4
instruction = "CP A"



//Row 13: 0x0C
case 0xC0:
		//0xC0: RET NZ

//RET()
cycles = 20 // 8
instruction = "RET NZ"

case 0xC1:
		//0xC1: POP BC

//POP()
cycles = 12
instruction = "POP BC"

case 0xC2:
		//0xC2: JP NZ,a16

//JP()
cycles = 16 // 12
instruction = "JP NZ,a16"

case 0xC3:
		//0xC3: JP a16

//JP()
cycles = 16
instruction = "JP a16"

case 0xC4:
		//0xC4: CALL NZ,a16

//CALL()
cycles = 24 // 12
instruction = "CALL NZ,a16"

case 0xC5:
		//0xC5: PUSH BC

//PUSH()
cycles = 16
instruction = "PUSH BC"

case 0xC6:
		//0xC6: ADD A,d8

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADD")
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 8
instruction = "ADD A,d8"

case 0xC7:
		//0xC7: RST 00H

//RST()
cycles = 16
instruction = "RST 00H"

case 0xC8:
		//0xC8: RET Z

//RET()
cycles = 20 // 8
instruction = "RET Z"

case 0xC9:
		//0xC9: RET

//RET()
cycles = 16
instruction = "RET"

case 0xCA:
		//0xCA: JP Z,a16

//JP()
cycles = 16 // 12
instruction = "JP Z,a16"

case 0xCB:
		//0xCB: PREFIX CB

//()
cycles = 4
instruction = "PREFIX CB"

case 0xCC:
		//0xCC: CALL Z,a16

//CALL()
cycles = 24 // 12
instruction = "CALL Z,a16"

case 0xCD:
		//0xCD: CALL a16

//CALL()
cycles = 24
instruction = "CALL a16"

case 0xCE:
		//0xCE: ADC A,d8

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"ADC")
		//cpu.cflag( , ,"ADC")
//ADC()
cycles = 8
instruction = "ADC A,d8"

case 0xCF:
		//0xCF: RST 08H

//RST()
cycles = 16
instruction = "RST 08H"



//Row 14: 0x0D
case 0xD0:
		//0xD0: RET NC

//RET()
cycles = 20 // 8
instruction = "RET NC"

case 0xD1:
		//0xD1: POP DE

//POP()
cycles = 12
instruction = "POP DE"

case 0xD2:
		//0xD2: JP NC,a16

//JP()
cycles = 16 // 12
instruction = "JP NC,a16"

case 0xD4:
		//0xD4: CALL NC,a16

//CALL()
cycles = 24 // 12
instruction = "CALL NC,a16"

case 0xD5:
		//0xD5: PUSH DE

//PUSH()
cycles = 16
instruction = "PUSH DE"

case 0xD6:
		//0xD6: SUB d8

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SUB")
		//cpu.cflag( , ,"SUB")
//SUB()
cycles = 8
instruction = "SUB d8"

case 0xD7:
		//0xD7: RST 10H

//RST()
cycles = 16
instruction = "RST 10H"

case 0xD8:
		//0xD8: RET C

//RET()
cycles = 20 // 8
instruction = "RET C"

case 0xD9:
		//0xD9: RETI

//RETI()
cycles = 16
instruction = "RETI"

case 0xDA:
		//0xDA: JP C,a16

//JP()
cycles = 16 // 12
instruction = "JP C,a16"

case 0xDC:
		//0xDC: CALL C,a16

//CALL()
cycles = 24 // 12
instruction = "CALL C,a16"

case 0xDE:
		//0xDE: SBC A,d8

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"SBC")
		//cpu.cflag( , ,"SBC")
//SBC()
cycles = 8
instruction = "SBC A,d8"

case 0xDF:
		//0xDF: RST 18H

//RST()
cycles = 16
instruction = "RST 18H"



//Row 15: 0x0E
case 0xE0:
		//0xE0: LDH (a8),A

//LDH()
cycles = 12
instruction = "LDH (a8),A"

case 0xE1:
		//0xE1: POP HL

//POP()
cycles = 12
instruction = "POP HL"

case 0xE2:
		//0xE2: LD (C),A

//LD()
cycles = 8
instruction = "LD (C),A"

case 0xE5:
		//0xE5: PUSH HL

//PUSH()
cycles = 16
instruction = "PUSH HL"

case 0xE6:
		//0xE6: AND d8

cpu.zflag()
cpu.setFlag("N",0)
cpu.setFlag("H",1)
cpu.setFlag("C",0)
//AND()
cycles = 8
instruction = "AND d8"

case 0xE7:
		//0xE7: RST 20H

//RST()
cycles = 16
instruction = "RST 20H"

case 0xE8:
		//0xE8: ADD SP,r8

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"ADD")
//ADD()
cycles = 16
instruction = "ADD SP,r8"

case 0xE9:
		//0xE9: JP (HL)

//JP()
cycles = 4
instruction = "JP (HL)"

case 0xEA:
		//0xEA: LD (a16),A

//LD()
cycles = 16
instruction = "LD (a16),A"

case 0xEE:
		//0xEE: XOR d8

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"XOR")
cpu.setFlag("C",0)
//XOR()
cycles = 8
instruction = "XOR d8"

case 0xEF:
		//0xEF: RST 28H

//RST()
cycles = 16
instruction = "RST 28H"



//Row 16: 0x0F
case 0xF0:
		//0xF0: LDH A,(a8)

//LDH()
cycles = 12
instruction = "LDH A,(a8)"

case 0xF1:
		//0xF1: POP AF

cpu.zflag()
		//cpu.hflag( , ,"POP")
		//cpu.cflag( , ,"POP")
//POP()
cycles = 12
instruction = "POP AF"

case 0xF2:
		//0xF2: LD A,(C)

//LD()
cycles = 8
instruction = "LD A,(C)"

case 0xF3:
		//0xF3: DI

//DI()
cycles = 4
instruction = "DI"

case 0xF5:
		//0xF5: PUSH AF

//PUSH()
cycles = 16
instruction = "PUSH AF"

case 0xF6:
		//0xF6: OR d8

cpu.zflag()
cpu.setFlag("N",0)
		//cpu.hflag( , ,"OR")
cpu.setFlag("C",0)
//OR()
cycles = 8
instruction = "OR d8"

case 0xF7:
		//0xF7: RST 30H

//RST()
cycles = 16
instruction = "RST 30H"

case 0xF8:
		//0xF8: LD HL,SP+r8

cpu.setFlag("Z",0)
cpu.setFlag("N",0)
cpu.setFlag("H",0)
		//cpu.cflag( , ,"LD")
//LD()
cycles = 12
instruction = "LD HL,SP+r8"

case 0xF9:
		//0xF9: LD SP,HL

//LD()
cycles = 8
instruction = "LD SP,HL"

case 0xFA:
		//0xFA: LD A,(a16)

//LD()
cycles = 16
instruction = "LD A,(a16)"

case 0xFB:
		//0xFB: EI

//EI()
cycles = 4
instruction = "EI"

case 0xFE:
		//0xFE: CP d8

cpu.zflag()
cpu.setFlag("N",1)
		//cpu.hflag( , ,"CP")
		//cpu.cflag( , ,"CP")
//CP()
cycles = 8
instruction = "CP d8"

case 0xFF:
		//0xFF: RST 38H

//RST()
cycles = 16
instruction = "RST 38H"

}
}
