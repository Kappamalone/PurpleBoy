package main

func (cpu *gameboyCPU) handleInterrupts() {
	//Interrupt service routine
	if !cpu.IME {
		if cpu.IE != 0 && cpu.IF != 0 {
			cpu.HALT = false
		}
		return
	}

	//Handle an interupt if interrupt exists
	if cpu.IE != 0 && cpu.IF != 0 {
		cpu.PUSH(cpu.PC)
		cpu.HALT = false

		if bitSet(cpu.IE, 0) && bitSet(cpu.IF, 0) {
			cpu.VBlank()
		} else if bitSet(cpu.IE, 1) && bitSet(cpu.IF, 1) {
			cpu.LCDSTAT()
		} else if bitSet(cpu.IE, 2) && bitSet(cpu.IF, 2) {
			cpu.TIMER()
		} else if bitSet(cpu.IE, 3) && bitSet(cpu.IF, 3) {
			cpu.SERIAL()
		} else if bitSet(cpu.IE, 4) && bitSet(cpu.IF, 4) {
			cpu.JOYPAD()
		}
		
		cpu.cycles += 20 //Takes a total of 5 machine cycles
	}
}

/*
The general gist of handling these 5 interrupts are as follows.
1) Disable interrupt request held in the IF register
2) Jump to a INT VEC
3) Disable IME to prevent any more interrupts from being serviced
*/

func (cpu *gameboyCPU) VBlank() {
	//V-blank
	cpu.IF &= (0xFF-0x01)
	cpu.PC = 0x0040
	cpu.IME = false
	cpu.gb.debug.printConsole("VBLANK\n", "cyan")
}

func (cpu *gameboyCPU) LCDSTAT() {
	//LCD STAT
	cpu.IF &= (0xFF-0x02)
	cpu.PC = 0x0048
	cpu.IME = false
	cpu.gb.debug.printConsole("LCD STAT\n", "cyan")
}

func (cpu *gameboyCPU) TIMER() {
	//Timer
	cpu.IF &= (0xFF-0x04)
	cpu.PC = 0x0050
	cpu.IME = false
	cpu.gb.debug.printConsole("TIMER\n", "cyan")
}

func (cpu *gameboyCPU) SERIAL() {
	//Serial
	cpu.IF &= (0xFF-0x08)
	cpu.PC = 0x0058
	cpu.IME = false
	cpu.gb.debug.printConsole("SERIAL\n", "cyan")
}

func (cpu *gameboyCPU) JOYPAD() {
	//Joypad
	cpu.IF &= (0xFF-0x10)
	cpu.PC = 0x0060
	cpu.IME = false
	cpu.gb.debug.printConsole("JOYPAD\n", "cyan")
}
