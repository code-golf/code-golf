#!/usr/bin/env node

const fs = require("fs");
const { Runtime } = require("/usr/lib/stax/stax");

function showDebug(runtime) {
    const state = runtime.getDebugState();

    const lines = [
        `X: ${state.x}`,
        `Y: ${state.y}`,
        `i: ${state.index}`,
        `_: ${state._}`,
        "Main stack",
        ...state.main.map((value, index) => `${index + 1}. ${value}`),
        "Input stack",
        ...state.input.map((value, index) => `${index + 1}. ${value}`),
        "",
    ];

    process.stdout.write(lines.join("\n") + "\n");
}


if (process.argv[2] !== "--version") {
    const program = fs.readFileSync(0, "utf8");
    const stdin = process.argv.slice(3);
    const runtime = new Runtime(output => process.stdout.write(output));

    for (const state of runtime.runProgram(program, stdin)) {
        // Handle debug output from "|`"
        if (state.break) {
            showDebug(runtime)
        }
    }
}
