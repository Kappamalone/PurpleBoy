package main

var (
	interrupts = [5]string{"VBLANK ", "LCDSTAT", "TIMER  ", "SERIAL ", "JOYPAD "}
)

func (cpu *gameboyCPU) handleInterrupts() {
	//Interrupt service routine
	if ((cpu.IE & cpu.IF) & 0x1F) != 0 {
		cpu.HALT = false
	}

	if !cpu.IME {
		return
	}

	//Handle an interupt if interrupt exists
	if (cpu.IE & cpu.IF) != 0 {
		for bit := 0; bit < 5; bit++ {
			if bitSet(cpu.IE, uint8(bit)) && bitSet(cpu.IF, uint8(bit)) {
				cpu.PUSH(cpu.PC)
				cpu.IF &= ^(1 << bit)                 //Disable requested interrupt
				cpu.PC = uint16(0x0040 + (0x8 * bit)) //Jump to INT vec
				cpu.IME = false
				cpu.cycles += 20

				if isDebugging {
					cpu.gb.debug.printInterrupt(interrupts[bit])
				}

				break //Only service one interrupt at a time
			}
		}
	}
}

func (cpu *gameboyCPU) requestVblank() {
	cpu.IF |= 0x1
}

func (cpu *gameboyCPU) requestSTAT() {
	cpu.IF |= 0x2
}

func (cpu *gameboyCPU) requestTimer() {
	cpu.IF |= 0x4
}

func (cpu *gameboyCPU) requestSerial() {
	cpu.IF |= 0x8
}

func (cpu *gameboyCPU) requestJoypad() {
	cpu.IF |= 0x10
}
