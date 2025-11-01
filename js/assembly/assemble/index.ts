import {
    Assembler,
    Message,
    CodeFragment,
    ToolResult,
    Result,
} from '../assembler-interface';
import {
    Program,
    Result as ProgramResult,
    wasmCreateProgram,
    UTF8ArrayToString
} from '../wasm-program';

function parseAsMessages(stderr: string): Message[] {
    const messages: Message[] = [];
    const lines = stderr.split('\n');

    for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim();
        if (!line) continue;

        const match = line.match(/^[^:]+:(\d+): (Error|Warning): (.+)$/);
        if (match) {
            messages.push({
                type: match[2].toLowerCase() as 'error' | 'warning',
                line: parseInt(match[1]),
                message: match[3]
            });
        }
    }

    return messages;
}

const ldRe = /^ld: (warning: )?([^ ]+:)?(.*)$/;
const ldLineRe = /^(.*):([0-9]+):(?:\(([^)]*)\):)?$/;
const ldInFunctionRe = /^ld: [^:]*: in function /;

function parseLdMessages(stderr: string): Message[] {
    const messages: Message[] = [];
    const lines = stderr.split('\n');
    let curPrefix = "";

    for (const line of lines) {
        const line_trim = line.trim();
        if (!line_trim) continue;

        const match = (curPrefix + line_trim).match(ldRe);
        if (match) {
            const lineMatch = match[2]?.match(ldLineRe);
            const line = lineMatch ? parseInt(lineMatch[2]) : undefined;
            messages.push({
                type: match[1] ? 'warning' : 'error',
                line,
                message: match[3].trim()
            });
        }

        if(ldInFunctionRe.test(line_trim)) {
            curPrefix = "ld: ";
        } else {
            curPrefix = "";
        }
    }

    return messages;
}

const hexDumpRe = /^ ([0-9a-fA-F]+)\s+(.*)/;
const debugLineRestRe = /^\s*([0-9]+|-)\s+(0x[0-9a-fA-F]+|0)(\s|$)/;
const sectionLineRe = /^Contents of section (.*):$/;

function parseObjdumpOutput(output: string, filename: string): CodeFragment[] {
    console.log("OBJDUMP: ", output);
    const lines = output.split('\n');
    let part = 0;
    let bytes = new Map<number, number>();
    let sections = new Map<number, string>();

    let preFragments = [];
    const filenameSp = filename + " ";
    let lastLineNum = undefined;

    let lastAddress = undefined;
    let section = undefined;
    for (const line of lines) {
        let thisLineNum = undefined;
        let thisAddress = undefined;
        const sectionMatch = line.match(sectionLineRe);
        if (sectionMatch) {
            section = sectionMatch[1];
            if(!section.startsWith(".debug_")) {
                part = 1;
            } else {
                part = 0;
            }
        } else if (line.startsWith('Contents of the .debug_line section:')) {
            part = 2;
            continue;
        } else if(part == 1) {
            const match = line.match(hexDumpRe);
            if(match) {
                let address = parseInt(match[1], 16);
                const hexDigits = match[2].substring(0, 35).replaceAll(' ', '');
                for(let i = 0; i < hexDigits.length; i += 2) {
                    bytes.set(address, parseInt(hexDigits.substring(i, i + 2), 16));
                    if(section)
                        sections.set(address, section);
                    ++address;
                }
            }
        } else if(part == 2) {
            if(line.startsWith(filenameSp)) {
                const rest = line.substring(filenameSp.length);
                const match = rest.match(debugLineRestRe);
                if(match) {
                    thisLineNum = match[1] === '-' ? undefined : parseInt(match[1]);
                    thisAddress = match[2] === "0" ? 0 : parseInt(match[2].substring(2), 16);
                    if(lastLineNum !== undefined && lastAddress !== undefined && lastAddress < thisAddress) {
                        preFragments.push({line: lastLineNum, address: lastAddress, length: thisAddress - lastAddress});
                    }
                }
            }
        }
        lastLineNum = thisLineNum;
        lastAddress = thisAddress;
    }

    const fragments = [];
    for(const preFragment of preFragments) {
        let fragSection = undefined;
        let fragBytes = undefined;
        let fragAddress = 0;
        for(let curAddress = preFragment.address; curAddress < preFragment.address + preFragment.length; ++curAddress) {
            const byte = bytes.get(curAddress);
            const curSection = sections.get(curAddress);
            bytes.delete(curAddress);
            sections.delete(curAddress);
            if(curSection !== undefined && curSection !== fragSection) {
                if(fragBytes) {
                    fragments.push({line: preFragment.line, address: fragAddress, section: fragSection, bytes: new Uint8Array(fragBytes)});
                }
                fragSection = curSection;
                fragBytes = undefined;
                fragAddress = curAddress;
            }
            if(byte !== undefined) {
                if(!fragBytes) {
                    fragBytes = [];
                }
                fragBytes.push(byte);
            }
        }
        if(fragBytes) {
            fragments.push({line: preFragment.line, address: fragAddress, section: fragSection, bytes: new Uint8Array(fragBytes)});
        }
    }

    const addresses = Array.from(bytes.keys());
    addresses.sort((a, b) => a - b);
    const fragBytes = [];
    let fragAddress = undefined;
    let fragSection = undefined;
    for(const address of addresses) {
        const section = sections.get(address);
        if(fragAddress !== undefined && (address != fragAddress + fragBytes.length || section !== fragSection)) {
            fragments.push({address: fragAddress, section: fragSection, bytes: new Uint8Array(fragBytes)});
            fragBytes.length = 0;
            fragAddress = address;
            fragSection = section;
        }
        fragBytes.push(bytes.get(address)!);
    }
    if(fragBytes.length > 0 && fragAddress !== undefined) {
        fragments.push({address: fragAddress, section: fragSection, bytes: new Uint8Array(fragBytes)});
    }

    return fragments;
}

const sourceName = "a.S";
const objectName = "a.o";
const executableName = "a.out";
const linkerScriptName = "a.lds";

export class ToolchainAssembler implements Assembler {
    private gas: Program;
    private ld: Program;
    private objdump: Program;
    private linkerScript: string | undefined = undefined;
    private asArgs: (src: string) => string[] = (_src: string) => [];
    private ldArgs: (src: string) => string[] = (_src: string) => [];

    private constructor(gas: Program, ld: Program, objdump: Program) {
        this.gas = gas;
        this.ld = ld;
        this.objdump = objdump;
    }

    assemble(source: string | Uint8Array): ToolResult<Uint8Array> {
        const messages: Message[] = [];

        try {
            const sourceStr = typeof source === "string" ? source : UTF8ArrayToString(source);

            // TODO: -mno-arch for RISC-V only; add proper options
            const asArgs = ["as", "-g", "-o", objectName, sourceName];
            asArgs.push(...this.asArgs(sourceStr));
            const gasRes: ProgramResult = this.gas(
                asArgs,
                "",
                {[sourceName]: source}
            );

            messages.push(...parseAsMessages(UTF8ArrayToString(gasRes.stderr)));

            console.log("gas", UTF8ArrayToString(gasRes.stderr));

            if (gasRes.status !== 0) {
                return { messages };
            }

            const objectFile = gasRes.files[objectName];
            if (!objectFile) {
                 messages.push({ type: 'error', message: 'Assembler succeeded but did not produce an object file.'});
                 return { messages };
            }

            const ldArgs = ["ld", objectName, "-o", executableName];
            const ldFiles: Record<string, Uint8Array | string> = {[objectName]: objectFile};
            if(this.linkerScript) {
                ldArgs.push("-T", linkerScriptName);
                ldFiles[linkerScriptName] = this.linkerScript;
            }
            ldArgs.push(...this.ldArgs(sourceStr));
            console.log(this.linkerScript);
            const ldRes: ProgramResult = this.ld(
                ldArgs,
                "",
                ldFiles
            );

            messages.push(...parseLdMessages(UTF8ArrayToString(ldRes.stderr)));

            if (ldRes.status !== 0) {
                return { messages };
            }

            const executable = ldRes.files["a.out"];
            if (!executable) {
                messages.push({ type: 'error', message: 'Linker succeeded but did not produce an output file.'});
                return { messages };
            }

            return { value: new Uint8Array(executable), messages };
        } catch (e: any) {
            console.error(e);
            messages.push({type: 'error', message: e.message || "An unexpected error occurred during assembly."});
            return { messages };
        }
    }

    parseExecutable(executable: Uint8Array): Result<CodeFragment[]> {
        try {
            const objdumpRes = this.objdump(
                ["objdump", "-x", "-s", "--dwarf=decodedline", executableName],
                "",
                {[executableName]: executable}
            );

            if (objdumpRes.status !== 0) {
                return { error: UTF8ArrayToString(objdumpRes.stderr) };
            }

            const output = UTF8ArrayToString(objdumpRes.stdout);
            const value = parseObjdumpOutput(output, sourceName);

            return { value };
        } catch (e: any) {
            console.error(e);
            return { error: e.message || "An unexpected error occurred during parsing." };
        }
    }

    clone(): ToolchainAssembler {
        return Object.assign(new ToolchainAssembler(this.gas, this.ld, this.objdump), this);
    }

    withLinkerScript(linkerScript: string): ToolchainAssembler {
        const obj = this.clone();
        obj.linkerScript = linkerScript;
        return obj;
    }

    withAsArgs(asArgs: (src: string) => string[]): ToolchainAssembler {
        const obj = this.clone();
        obj.asArgs = asArgs;
        return obj;
    }

    withLdArgs(ldArgs: (src: string) => string[]): ToolchainAssembler {
        const obj = this.clone();
        obj.ldArgs = ldArgs;
        return obj;
    }

    static async create(as: Promise<WebAssembly.Module>, ld: Promise<WebAssembly.Module>, objdump: Promise<WebAssembly.Module>): Promise<ToolchainAssembler> {
        const [asProg, ldProg, objdumpProg] = await Promise.all([
            as.then(x => wasmCreateProgram(x, true)),
            ld.then(x => wasmCreateProgram(x, true)),
            objdump.then(x => wasmCreateProgram(x, true)),
        ]);

        return new ToolchainAssembler(asProg, ldProg, objdumpProg);
    }
}

export class PreprocessedAssembler implements Assembler {
    private constructor(private preprocessor: Program, private assembler: Assembler) {
    }

    assemble(source: string): ToolResult<Uint8Array> {
        const asppRes: ProgramResult = this.preprocessor(
            ["aspp"],
            source,
            {}
        );

        if(asppRes.status !== 0) {
            const messages = [{ type: 'error', message: 'Aspp failed: ' + UTF8ArrayToString(asppRes.stderr)}];
            return { messages };
        }

        const result = this.assembler.assemble(asppRes.stdout);
        for(const message of result.messages) {
            // remove column numbers since preprocessing makes them incorrect
            message.startColumn = undefined;
            message.endColumn = undefined;
        }
        return result;
    }

    parseExecutable(executable: Uint8Array): Result<CodeFragment[]> {
        return this.assembler.parseExecutable(executable);
    }

    static async create(preprocessor: Promise<WebAssembly.Module>, assembler: Promise<Assembler>): Promise<PreprocessedAssembler> {
        const [preprocessorProg, assemblerValue] = await Promise.all([
            preprocessor.then(x => wasmCreateProgram(x, true)),
            assembler
        ]);

        return new PreprocessedAssembler(preprocessorProg, assemblerValue);
    }
}
