import { __exportStar } from "tslib";
import { O_RDONLY, O_RDWR, O_WRONLY, WasmExports } from "./abi";
import { BufferFile, CloseInterface, defaultResizeMemoryInterface, DirlessFsInterface, Fd, File, initAndCallMain, installAbort, installAssert, installClose, installDirlessFileSystem, installZeroClock, installEmptyEnviron, installEnottyIoctl, installExit, installFstat, installRead, installResizeHeap, installSeek, installWrite, LineBufferedOutFile, makeEmptyImports, makeThrowImports, ReadInterface, ResizeMemoryInterface, RoBufferFile, SeekInterface, setMemory, traceImports, WriteInterface, installFcntlDup, installPRead, installPWrite, installIndirectFunctionTable, setIndirectFunctionTable, IndirectFunctionTableInterface, installDummyRuntimeKeepAlive } from "./impl";
import { ResizableBuffer } from "../resizable-buffer";
import { stringToNewUTF8Array, UTF8ArrayToString } from "./utf8";
import { ptr } from "./mem";

export { UTF8ArrayToString } from "./utf8";

export async function wasmCompileUrl(url: string): Promise<WebAssembly.Module> {
    const responsePromise = fetch(url, {
        credentials: "same-origin"
    });
    try {
        return await WebAssembly.compileStreaming(responsePromise);
    } catch (e) {
        const response = await responsePromise;
        if (!response.ok)
            throw new Error(response.status + " : " + response.url);
        return await WebAssembly.compile(await response.arrayBuffer());
    }
}

export interface Result {
    status: number;
    stdout: Uint8Array;
    stderr: Uint8Array;
    stdinPos: number;
    files: Record<string, Uint8Array>;
}

export type Program = (args: string[], stdin: string | Uint8Array, files: Record<string, string | Uint8Array>, options?: {stderrHandler?: (line: string) => void}) => Result;

export async function wasmCreateProgram(module: WebAssembly.Module, reusable: boolean, options?: {reinstance?: boolean, traceSyscalls?: boolean, allowUndefinedImports?: boolean}): Promise<Program> {
    let intf: ResizeMemoryInterface & ReadInterface & WriteInterface & CloseInterface & SeekInterface & DirlessFsInterface & IndirectFunctionTableInterface = {
        fds: [],
        fs: new Map(),
        ...defaultResizeMemoryInterface()
    };

    let imports = options?.allowUndefinedImports ? makeThrowImports(module) : makeEmptyImports();
    installAbort(imports);
    installAssert(imports, intf);
    installExit(imports);
    installDummyRuntimeKeepAlive(imports);
    installResizeHeap(imports, intf);
    installIndirectFunctionTable(imports, intf);
    installEmptyEnviron(imports, intf);
    installZeroClock(imports, intf);
    installRead(imports, intf);
    installWrite(imports, intf);
    installClose(imports, intf);
    installSeek(imports, intf);
    installPRead(imports, intf);
    installPWrite(imports, intf);
    installFstat(imports, intf);
    installFcntlDup(imports, intf);
    installEnottyIoctl(imports, intf);
    installDirlessFileSystem(imports, intf);

    let instance: WebAssembly.Instance | undefined = undefined;
    const reuseInstance = !options?.reinstance;

    if(options?.traceSyscalls)
        imports = traceImports(imports);

    {
        instance = await WebAssembly.instantiate(module, imports);
        const exports = instance.exports as WasmExports;
        setMemory(intf, exports);
        setIndirectFunctionTable(intf, exports);
    }

    let saveMem: Uint8Array | undefined = undefined;
    let saveSP: number | undefined = undefined;

    return (args: string[], stdin: string | Uint8Array, files: Record<string, string | Uint8Array>, options?: {stderrHandler?: (line: string) => void}): Result => {
        if (!intf) {
            throw new Error("to call a program multiple times, you need to create it with the reusable flag set to true");
        }

        if(!instance) {
            instance = new WebAssembly.Instance(module, imports);
            const exports = instance.exports as WasmExports;
            setMemory(intf, exports);
            setIndirectFunctionTable(intf, exports);
        }

        if(reusable && reuseInstance) {
            if (!saveMem) {
                saveMem = new Uint8Array(intf.memU8);
            } else {
                intf.memU8.set(saveMem);
                intf.memorySize = saveMem.length;
            }
        }
        //console.log("run# " + (++runs) + " saveMem size: " + saveMem?.length + " memU8 size:", intf.memU8.length);

        const stdinFile = new RoBufferFile(new ResizableBuffer(stringToNewUTF8Array(stdin)));
        const stdoutFile = new BufferFile(new ResizableBuffer());
        const stderrFile = new BufferFile(new ResizableBuffer());
        let stderrFileC: File<BufferFile | LineBufferedOutFile> = new File(stderrFile, O_RDWR);
        const stderrHandler = options?.stderrHandler;
        if(stderrHandler) {
            const stderrLog = new LineBufferedOutFile((buf) => {
                let str = UTF8ArrayToString(buf);
                if(str.endsWith("\n"))
                    str = str.slice(0, -1);
                stderrHandler(str);
                stderrFile.buffer.append(buf, 0, buf.length);
            });
            stderrFileC = new File(stderrLog, O_WRONLY);
        }

        intf.fds = [new Fd(new File(stdinFile, O_RDONLY), 0), new Fd(new File(stdoutFile, O_RDWR), 0), new Fd(stderrFileC, 0)];
        intf.fs.clear();
        for(const fileName in files) {
            const file = files[fileName];
            if(file) {
                intf.fs.set(fileName, new ResizableBuffer(stringToNewUTF8Array(file)));
            }
        }

        let status;
        {
            const exports = instance.exports as WasmExports;
            if(reusable && reuseInstance) {
                if(saveSP === undefined) {
                    // don't use emscripten_stack_get_current() since we also want to make sure that the stack pointer is aligned
                    saveSP = exports._emscripten_stack_alloc(0);
                } else {
                    if(exports._emscripten_stack_restore) {
                        exports._emscripten_stack_restore(saveSP);
                    } else {
                        const curSP = exports.emscripten_stack_get_current ? exports.emscripten_stack_get_current() : exports._emscripten_stack_alloc(0);
                        // stack_alloc just subtracts and aligns, so since saveSP is aligned, this will set the stack pointer to saveSP
                        exports._emscripten_stack_alloc(curSP - saveSP)
                    }
                }
            } else {
                // align the stack pointer
                exports._emscripten_stack_alloc(0);
            }
            status = initAndCallMain(exports, args);
        }

        const outFiles: Record<string, Uint8Array> = {};
        for(const [name, file] of intf.fs.entries()) {
            file.trim();
            outFiles[name] = file.array;
        }
        intf.fds = [];
        intf.fs.clear();

        const stdout = stdoutFile.buffer.array;
        const stderr = stderrFile.buffer.array;
        const stdinPos = stdinFile.position;

        if(!reuseInstance) {
            instance = undefined;
        }

        if (!reusable) {
            // free memory
            instance = undefined!;
            intf = undefined!;
        }

        return { status, stdout, stderr, stdinPos, files: outFiles };
    };
}
