#!/usr/bin/node
import { assemble, LocatedError, Machine } from "fluffy-6502"
import { readFileSync, writeFileSync, writeSync } from "fs";
const code = readFileSync(0, "utf-8");
try {
    const { memory } = assemble(code);
    let byte_count = new Machine(memory).nz_bytes();
    writeFileSync(3, byte_count.toString());
    let in_ptr = 0x4000;
    for(const arg of process.argv.slice(2)) {
        for(const byte of Buffer.from(arg, "utf-8")) {
            memory[in_ptr] = byte;
            in_ptr++;
        }
        in_ptr++;
    }
    const machine = new Machine(memory);
    machine.track_accesses = false;
    machine.run_until_jam();
    let output = [];
    for(let i = 0x8000; i < 0x10000; i++) {
        const byte = machine.memory[i];
        if(byte == 0) {
            break;
        } else {
            output.push(byte)
        }
    }
    writeSync(1, Buffer.from(output));
    console.error(`PC: ${machine.pc.toString(16)}`);
    console.error(`A: ${machine.a.toString(16)}`);
    console.error(`X: ${machine.x.toString(16)}`);
    console.error(`Y: ${machine.y.toString(16)}`);
    console.error(`S: ${machine.s.toString(16)}`);
    console.error(`Flags: ${machine.get_p(1).toString(2)}`);
} catch(e) {
    if(e instanceof LocatedError) {
        let lines = code.split("\n");
        let line = 0;
        let column = e.start;
        while(column >= lines[line].length) {
            column -= lines[line].length + 1;
            line++;
        }
        console.error(`Error at line ${line+1} column ${column+1}: ${e.message}`);
    } else {
        console.error(e);
    }
}