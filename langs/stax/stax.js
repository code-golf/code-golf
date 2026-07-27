#!/usr/bin/env node

const fs = require("fs");
const { Runtime } = require("/usr/lib/stax/stax");

if (process.argv[2] !== "--version") {
    const program = fs.readFileSync(0, "utf8");
    const stdin = process.argv.slice(3);
    const runtime = new Runtime(output => process.stdout.write(output));

    for (const _ of runtime.runProgram(program, stdin)) { }
}
