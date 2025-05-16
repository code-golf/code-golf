// This is a reasonable version of the Emscripten JavaScript glue code, which is
// unfortunately truly awful code.

// It's based on the original Emscripten code, but unlike that code, it is well-written,
// it's modular rather than having a single absurdly huge function that does everything,
// doesn't use 10 lines of code to perform actions that can be done with one line of code,
// doesn't include any useless functionality for purely input/output programs,
// has proper abstractions, doesn't make a lot of assumptions about what the user wants,
// provides a way to run a program multiple times efficiently, doesn't use wrong
// terminology like referring to the memory as "heap" (the heap is only part of the
// memory, along with the stack, data, and bss sections...), includes basic features like
// stdin support, uses correct code for resizing the heap rather than resizing the memory
// multiple times and starting to grow linearly at some point, doesn't overwrite Math.random
// and Date.now in deterministic mode

import { AT_EMPTY_PATH, AT_FDCWD, AT_REMOVEDIR, AT_SYMLINK_NOFOLLOW, EACCES, EBADF, EEXIST, EINVAL, ENOENT, ENOTDIR, ENOTTY, ERANGE, ESPIPE, EXDEV, F_DUPFD, F_DUPFD_CLOEXEC, F_GETFD, F_GETFL, F_GETLK, F_SETFD, F_SETFL, F_SETLK, F_SETLKW, FD_CLOEXEC, O_APPEND, O_CLOEXEC, O_CREAT, O_DIRECTORY, O_EXCL, O_RDWR, O_TRUNC, O_WRONLY, WasmExports, WasmImports } from "./abi";
import { bzero, MemInterface, ptr, read32, readString, write16, write32, write64, writeString, writeSyscallStat, writeWasiStat } from "./mem";
import {ResizableBuffer} from "../resizable-buffer";
import { lengthBytesUTF8, stringToUTF8Array0, UTF8ArrayToString } from "./utf8";

export function makeEmptyImports(): WasmImports {
    return {env: {}, wasi_snapshot_preview1: {}};
}

export function makeThrowImports(wasmModule: WebAssembly.Module): WasmImports {
    let imports = makeEmptyImports();
    let importDescs = WebAssembly.Module.imports(wasmModule);
    for(const {module, name} of importDescs) {
        if(!(module in imports)) {
            (imports as any)[module] = {};
        }
        const handler = () => {throw new Error("Unimplemented WASM import: " + module + "." + name);};
        (imports as any)[module][name] = handler;
    }
    return imports;
}

export function traceImports(imports: WasmImports): WasmImports {
    const outImports = makeEmptyImports();
    let module: keyof WasmImports;
    for(module in imports) {
        const origModule = imports[module];
        const outModule: any = {};
        for(const name in origModule) {
            const orig = (origModule as any)[name];
            const fname = module + "." + name;
            outModule[name] = function(...args: unknown[]) {
                console.log(fname + "(" + args.map(x => String(x)).join(",") + ")");
                const ret = orig(...args);
                console.log(fname + "(" + args.map(x => String(x)).join(",") + ") = " + ret);
                return ret;
            }
        }
        outImports[module] = outModule;
    }
    return outImports;
}

export interface MemoryInterface extends MemInterface {
    memory?: WebAssembly.Memory;
    memU8: Uint8Array;
}

export interface ResizeMemoryInterface extends MemoryInterface {
    memorySize: number;
    dirtySize: number;
    maxMemorySize: number;
    memoryGrowFactor: number;
}

export function setMemory(intf: Partial<ResizeMemoryInterface>, exports: WasmExports): void {
    const memory = exports.memory;
    const b = memory.buffer;
    intf.memory = memory;
    intf.memU8 = new Uint8Array(b);
    intf.memorySize = intf.memU8.length;
    intf.dirtySize = intf.memU8.length;
}

export function defaultResizeMemoryInterface(): ResizeMemoryInterface {
    return {
        memory: undefined,
        memorySize: 0,
        dirtySize: 0,
        memU8: new Uint8Array(0),
        maxMemorySize: 2147483648,
        memoryGrowFactor: Math.pow(2.0, 1.0 / 4.0),
    };
}

export function installAbort(imports: WasmImports): void {
    function _abort_js(): never {
        throw new WebAssembly.RuntimeError("Aborted");
    }
    imports.env._abort_js = _abort_js;
}

export function installAssert(imports: WasmImports, intf: MemInterface): void {
    function __assert_fail(condition: number, filename: number, line: number, func: number): never {
        const file = filename ? readString(intf, filename) : "unknown filename";
        const funcName = func ? readString(intf, func) : "unknown function";
        throw new WebAssembly.RuntimeError(`Assertion failed: ${readString(intf, condition)}, at: ${file}:${line}:${funcName}`);
    }
    imports.env.__assert_fail = __assert_fail;
}

export class ExitStatus extends Error {
    name: string = "ExitStatus";
    status: number;

    constructor(status: number) {
        super(`Program terminated with exit(${status})`);
        this.status = status;
        Object.setPrototypeOf(this, ExitStatus.prototype);
    }
}

export function installExit(imports: WasmImports): void {
    function exit(status: number): never {
        throw new ExitStatus(status);
    }
    imports.env.exit = exit;
    imports.wasi_snapshot_preview1.proc_exit = exit;
}

export function installDummyRuntimeKeepAlive(imports: WasmImports): void {
    imports.env._emscripten_runtime_keepalive_clear = function _emscripten_runtime_keepalive_clear() {};
}

export function installEmptyEnviron(imports: WasmImports, intf: MemInterface): void {
    imports.wasi_snapshot_preview1.environ_sizes_get = function environ_sizes_get(penviron_count: number, penviron_buf_size: number): number {
        write32(intf, penviron_count, 0);
        write32(intf, penviron_buf_size, 0);
        return 0;
    }
    imports.wasi_snapshot_preview1.environ_get = function environ_get(__environ: number, environ_buf: number): number {
        return 0;
    }
}

export function installResizeHeap(imports: WasmImports, intf: ResizeMemoryInterface): void {
    function emscripten_resize_heap(requestedSize: number): boolean {
        const oldPhysicalSize = intf.memU8.length;
        const oldLogicalSize = intf.memorySize ?? oldPhysicalSize;
        requestedSize >>>= 0;
        if (requestedSize <= oldLogicalSize)
            return true;

        const maxMemorySize = intf.maxMemorySize;
        const memory = intf.memory;
        if (!memory || requestedSize > maxMemorySize) {
            return false;
        }

        let newSizePages = Math.max(1, (oldPhysicalSize + 65535) >> 16);
        let requestedSizePages = (requestedSize + 65535) >> 16;

        if(requestedSizePages > newSizePages) {
            while (requestedSizePages > newSizePages) {
                newSizePages *= intf.memoryGrowFactor;
            }

            newSizePages = newSizePages | 0;

            const growPages = newSizePages - (oldPhysicalSize >> 16);
            try {
                memory.grow(growPages);
            } catch (e) {
                return false;
            }
            const b = memory.buffer;
            intf.memU8 = new Uint8Array(b);
        }

        const requestedSizeRoundedUpToPages = requestedSizePages << 16;

        const toClearEnd = Math.min(intf.dirtySize ?? oldPhysicalSize, requestedSizeRoundedUpToPages);
        if(oldLogicalSize < toClearEnd)
            intf.memU8.fill(0, oldLogicalSize, toClearEnd);

        intf.memorySize = requestedSizeRoundedUpToPages;
        intf.dirtySize = Math.max(intf.dirtySize, intf.memorySize);

        return true;
    }

    imports.env.emscripten_resize_heap = emscripten_resize_heap;
}

export class File<T> {
    constructor(public obj: T & {close?(): number}, public flags: number, private rc: number = 1) {
    }

    ref() {
        ++this.rc;
        return this;
    }

    close() {
        if(--this.rc == 0) {
            if(this.obj.close) {
                return this.obj.close();
            }
        }
        return 0;
    }
}

export class Fd<T> {
    constructor(public file: File<T>, public flags: number) {}
}

export interface FileInterface<F> extends MemInterface {
    fds: Fd<F>[]
}

export function getFile<T>(intf: FileInterface<T>, fd: number): File<T> | undefined {
    const fdo = intf.fds[fd];
    if(!fdo) {
        return undefined;
    }
    return fdo.file;
}

function fd_maybep_readwrite<T>(intf: FileInterface<T>, getOp: (file: File<T>) => ((array: Uint8Array, off: number, len: number) => number) | undefined, fd: number, iov: number, iovcnt: number, pnum: number, start: (file: File<T>) => number, finish: (file: File<T>, startData: number) => void): number {
    const file = getFile(intf, fd);
    let num = 0;
    let ret = 0;
    if (!file) {
        ret = -EBADF;
    } else {
        let op = getOp(file);
        if(!op) {
            console.log("EINVAL: " + JSON.stringify(file));
            ret = -EINVAL;
        } else {
            const startData = start(file);
            if(startData < 0) {
                ret = startData;
            } else {
                op = op.bind(file.obj);
                for (let i = 0; i < iovcnt; i++) {
                    const ptr = read32(intf, iov);
                    const len = read32(intf, iov + 4);
                    /*
                    if(trace) {
                        console.log("RW " + kind + " " + fd + " " + ptr + " " + len);
                        if(kind == "write")
                            console.log("TEXT: " + UTF8ArrayToString(intf.memU8, ptr, len));
                    }
                    */

                    iov += 8;
                    const res = op(intf.memU8, ptr, len);
                    if(res < 0) {
                        ret = res;
                        break;
                    }
                    /*
                    if(trace && kind == "read") {
                        console.log("TEXT: " + UTF8ArrayToString(intf.memU8, ptr, len));
                    }
                    */

                    num += res;
                    if(res < len) {
                        break;
                    }
                }

                finish(file, startData);
            }
        }
    }
    /*
    if(trace) {
        console.log("RW " + kind + " = " + ret + " bytes " + num);
    }
    */
    write32(intf, pnum, num);
    return ret;
}

function fd_readwrite<T>(intf: FileInterface<T>, getOp: (file: File<T>) => ((array: Uint8Array, off: number, len: number) => number) | undefined, fd: number, iov: number, iovcnt: number, pnum: number): number {
    return -fd_maybep_readwrite(intf, getOp, fd, iov, iovcnt, pnum, () => 0, () => {});
}

function fd_preadwrite<T>(intf: FileInterface<T & SeekFile>, getOp: (file: File<T>) => ((array: Uint8Array, off: number, len: number) => number) | undefined, fd: number, iov: number, iovcnt: number, pnum: number, offset: number): number {
    return -fd_maybep_readwrite(intf, getOp, fd, iov, iovcnt, pnum, (file) => {
        const position = file.obj.position;
        if (position === undefined || file.obj.seek === undefined) {
            return -ESPIPE;
        }
        const seekRes = file.obj.seek(offset);
        if(seekRes < 0) {
            return seekRes;
        }
        return position;
    }, (file, position: number) => {
        file.obj.seek!(position);
    });
}

export type ReadInterface = FileInterface<{
    read?: (ta: Uint8Array, pos: number, len: number) => number;
}>;

export function installRead(imports: WasmImports, intf: ReadInterface): void {
    imports.wasi_snapshot_preview1.fd_read = function fd_read(fd: number, iov: number, iovcnt: number, pnum: number) {
        return -fd_readwrite(intf, file => file.obj.read, fd, iov, iovcnt, pnum);
    }
}

export type WriteInterface = FileInterface<{
    write?: (ta: Uint8Array, pos: number, len: number) => number;
}>

export function installWrite(imports: WasmImports, intf: WriteInterface): void {
    imports.wasi_snapshot_preview1.fd_write = function fd_read(fd: number, iov: number, iovcnt: number, pnum: number) {
        return -fd_readwrite(intf, file => file.obj.write, fd, iov, iovcnt, pnum);
    }
}

export type CloseInterface = FileInterface<{
    close?: () => number;
}>;

export function installClose(imports: WasmImports, intf: CloseInterface) {
    function fd_close_negerr(fd: number) {
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        delete intf.fds[fd];
        return file.close();
    }
    imports.wasi_snapshot_preview1.fd_close = function fd_close(fd) {
        return -fd_close_negerr(fd);
    }
}

export function installEnottyIoctl(imports: WasmImports, intf: FileInterface<unknown>) {
    imports.env.__syscall_ioctl = function __syscall_ioctl(fd) {
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        return -ENOTTY;
    }
}

export type FstatInterface = FileInterface<{
    size?: number;
}>

export function installFstat(imports: WasmImports, intf: FstatInterface) {
    function fd_fdstat_get_negerr(fd: number, buf: ptr) {
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        writeWasiStat(intf, buf, file.obj.size);
        return 0;
    }
    imports.wasi_snapshot_preview1.fd_fdstat_get = function fd_fdstat_get(fd, buf) {
        return -fd_fdstat_get_negerr(fd, buf);
    }
    imports.env.__syscall_fstat64 = function __syscall_fstat64(fd, buf) {
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        writeSyscallStat(intf, buf, file.obj.size);
        return 0;
    }
}

export function installFcntlDup(imports: WasmImports, intf: FileInterface<unknown>) {
    imports.env.__syscall_dup = function __syscall_dup(fd: number) {
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        const newfd = intf.fds.length;
        intf.fds[newfd] = new Fd(file.ref(), 0);
        return newfd;
    }

    imports.env.__syscall_fcntl64 = function __syscall_fcntl64(fd, cmd, args) {
        const fdo = intf.fds[fd];
        if(!fdo) {
            return -EBADF;
        }
        const file = fdo.file;
        switch (cmd) {
            case F_DUPFD:
            case F_DUPFD_CLOEXEC: {
                let newfd = read32(intf, args);
                if (newfd < 0) {
                    return -EINVAL;
                }
                while (intf.fds[newfd]) {
                    ++newfd;
                }
                intf.fds[newfd] = new Fd(file.ref(), (cmd == F_DUPFD_CLOEXEC) ? FD_CLOEXEC : 0);
                return newfd;
            }
            case F_GETFD:
                return fdo.flags;
            case F_SETFD:
                fdo.flags = read32(intf, args) & FD_CLOEXEC;
                return 0;
            case F_GETFL:
                return file.flags;
            case F_SETFL:
                return 0;
            case F_GETLK: {
                const arg = read32(intf, args);
                write16(intf, arg, 2);
                return 0;
            }
            case F_SETLK:
            case F_SETLKW:
                return 0;
            default:
                return -EINVAL;
        }
    }
}

interface SeekFile {
    seek?(offset: number): number;
    position?: number;
    size?: number;
}

export type SeekInterface = FileInterface<SeekFile>;

export function installSeek(imports: WasmImports, intf: SeekInterface) {
    function fd_seek_negerr(fd: number, noffset: bigint, whence: number, pnewOffset: ptr) {
        const offset = Number(noffset);
        const file = getFile(intf, fd);
        if (!file) {
            return -EBADF;
        }
        if(file.obj.seek === undefined) {
            return -ESPIPE;
        }
        let newPosition;
        if(whence === 0) {
            newPosition = offset;
        } else if(whence === 1) {
            const position = file.obj.position;
            if(position === undefined) {
                return -ESPIPE;
            }
            newPosition = position + offset;
        } else if(whence === 2) {
            const size = file.obj.size;
            if(size === undefined) {
                return -ESPIPE;
            }
            newPosition = size + offset;
        } else {
            return -EINVAL;
        }
        const res = file.obj.seek(newPosition);
        if(res < 0)
            return res;
        write64(intf, pnewOffset, newPosition);
        return 0;
    }
    imports.wasi_snapshot_preview1.fd_seek = function fd_seek(fd: number, noffset: bigint, whence: number, pnewOffset: ptr) {
        return -fd_seek_negerr(fd, noffset, whence, pnewOffset);
    }
}

export function installPRead(imports: WasmImports, intf: SeekInterface & ReadInterface): void {
    imports.wasi_snapshot_preview1.fd_pread = function fd_pread(fd: number, iov: number, iovcnt: number, offset: bigint, pnum: number) {
        return -fd_preadwrite(intf, file => file.obj.read, fd, iov, iovcnt, pnum, Number(offset));
    }
}

export function installPWrite(imports: WasmImports, intf: SeekInterface & WriteInterface): void {
    imports.wasi_snapshot_preview1.fd_pwrite = function fd_pwrite(fd: number, iov: number, iovcnt: number, offset: bigint, pnum: number) {
        return -fd_preadwrite(intf, file => file.obj.write, fd, iov, iovcnt, pnum, Number(offset));
    }
}

export function installZeroClock(imports: WasmImports, intf: MemInterface): void {
    imports.env.emscripten_date_now = function emscripten_date_now() {
        return 0;
    }
    imports.env._tzset_js = function _tzset_js(timezone, daylight, std_name, dst_name) {
        write32(intf, timezone, 0);
        write32(intf, daylight, 0);
        writeString(intf, std_name, 17, "");
        writeString(intf, dst_name, 17, "");
        return 0;
    }
    imports.env._localtime_js = function _localtime_js(_time, tmPtr) {
        bzero(intf, tmPtr, 40);
        return 0;
    }
    imports.wasi_snapshot_preview1.clock_time_get = function _clock_time_get(_clk_id, _precision, ptime) {
        write64(intf, ptime, 0);
        return 0
    }
    // the timer never triggers, since our simulated CPU is infinitely fast and thus clock time never passes
    // this is suitable for computation-only programs that shouldn't be setting timers in the first place
    imports.env._setitimer_js = function __setitimer_js(_which, _timeout_ms) {
        return 0;
    };
}

export interface IndirectFunctionTableInterface {
    indirectFunctionTable?: WebAssembly.Table;
}

export function installIndirectFunctionTable(imports: WasmImports, intf: IndirectFunctionTableInterface): void {
    imports.env.__call_sighandler = function(fp: ptr, sig: number) {
        if(intf.indirectFunctionTable) {
            return intf.indirectFunctionTable.get(fp)(sig);
        }
    }
}

export function setIndirectFunctionTable(intf: IndirectFunctionTableInterface, exports: WasmExports): void {
    intf.indirectFunctionTable = exports.__indirect_function_table;
}

export class BufferFile {
    constructor(public buffer: ResizableBuffer, public position = 0) {
    }

    get size() {
        return this.buffer.length;
    }

    read?(ta: Uint8Array, pos: number, len: number): number {
        const n = this.buffer.read(this.position, ta, pos, len);
        this.position += n;
        return n;
    }

    write?(ta: Uint8Array, pos: number, len: number): number {
        const n = this.buffer.write(this.position, ta, pos, len);
        this.position += n;
        return n;
    }

    seek(position: number) {
        this.position = position;
        return 0;
    }
}

export class RoBufferFile extends BufferFile {
    override write = undefined;
}

export class WoBufferFile extends BufferFile {
    override read = undefined;
}

export class LineBufferedOutFile {
    private buffer = new ResizableBuffer();
    constructor(private handler: (array: Uint8Array) => void) {}

    write(ta: Uint8Array, pos: number, len: number): number{
        for (let i = 0; i < len; i++) {
            const curr = ta[pos + i];
            this.buffer.push(curr);
            if (curr === 10) {
                this.handler(this.buffer.array);
                this.buffer.clear()
            }
        }
        return len;
    }

    close() {
        if (this.buffer.length > 0) {
            this.handler(this.buffer.array);
            this.buffer.clear();
        }
        return 0;
    }
}

export interface DirlessFsInterface extends FileInterface<{size?: number}> {
    fs: Map<string, ResizableBuffer>;
}

// doesn't support directories, symlinks, file modes, uid/gid and inode numbers
export function installDirlessFileSystem(imports: WasmImports, intf: DirlessFsInterface): void {
    const leadingSlashesRegex = /^\/+/;
    const leadingDotSlashesRegex = /^(\.\/)+/;
    const multiSlashes = /\/\/+/g;

    function resolveDirlessPath(intf: FileInterface<{size?: number}>, dirfd: number, ppath: ptr | string): string | number | File<{size?: number}> {
        let path = typeof ppath === "string" ? ppath : readString(intf, ppath);
        path = path.replaceAll(multiSlashes, "/");
        if(!path.startsWith("/")) {
            if(dirfd != AT_FDCWD) {
                const file = getFile(intf, dirfd);
                if(!file)
                    return -EBADF;
                return -ENOTDIR;
            }
        } else {
            path = path.replace(leadingSlashesRegex, '');
        }
        path = path.replace(leadingDotSlashesRegex, '');
        if(path.includes("/")) {
            if(path.startsWith("/proc/self/fd/")) {
                const fd = Number.parseInt(path.substring(14));
                const file = getFile(intf, fd);
                if(!file) {
                    return -ENOENT;
                }
                return file;
            }
            if(path === "/dev/stdin")
                return 0;
            if(path === "/dev/stdout")
                return 1;
            if(path === "/dev/stderr")
                return 2;
            return -ENOENT;
        }
        return path;
    }

    function fstatat(dirfd: number, ppath: ptr, buf: ptr, flags: number) {
        const origPath = readString(intf, ppath);
        const path = resolveDirlessPath(intf, dirfd, origPath);
        if(typeof path === "number") {
            return path;
        }
        let size;
        if(typeof path === "object") {
            size = path.obj.size;
        } else if(origPath === "" && (flags & AT_EMPTY_PATH)) {
            const file = getFile(intf, dirfd);
            if(!file) {
                return -ENOENT;
            }
            size = file.obj.size;
        } else {
            const inode = intf.fs.get(path);
            if(!inode) {
                return -ENOENT;
            }
            size = inode.length;
        }
        writeSyscallStat(intf, buf, size);
        return 0;
    }

    imports.env.__syscall_stat64 = function __syscall_stat64(ppath, buf) {
        return fstatat(AT_FDCWD, ppath, buf, 0);
    }

    imports.env.__syscall_lstat64 = function __syscall_lstat64(ppath, buf) {
        return fstatat(AT_FDCWD, ppath, buf, AT_SYMLINK_NOFOLLOW);
    }

    imports.env.__syscall_newfstatat = function __syscall_newfstatat(dirfd: number, ppath: number, buf: number, flags: number) {
        return fstatat(dirfd, ppath, buf, flags);
    }

    imports.env.__syscall_faccessat = function __syscall_faccessat(dirfd: number, ppath: number, mode: number, flags: number) {
        const path = resolveDirlessPath(intf, dirfd, ppath);
        if(typeof path === "number") {
            return path;
        }
        if(typeof path === "object") {
            return 0;
        }
        if(!intf.fs.has(path)) {
            return -ENOENT;
        }
        return 0;
    }

    imports.env.__syscall_openat = function __syscall_openat(dirfd: number, ppath: number, flags: number, _args: number) {
        const path = resolveDirlessPath(intf, dirfd, ppath);
        if(typeof path === "number") {
            return path;
        }
        let file;
        if(typeof path === "object") {
            file = path.ref();
        } else {
            let inode = intf.fs.get(path);
            if(inode) {
                if(flags & (O_CREAT | O_EXCL)) {
                    return -EEXIST;
                }
                if(flags & O_DIRECTORY) {
                    return -ENOTDIR;
                }

                if(flags & O_TRUNC) {
                    inode.clear();
                }
            } else {
                if(!(flags & O_CREAT)) {
                    return -ENOENT;
                }
                if(path === "")
                    return -ENOENT;
                inode = new ResizableBuffer();
                intf.fs.set(path, inode);
            }
            const position = (flags & O_APPEND) ? inode.length : 0
            let obj;
            if(flags & O_WRONLY) {
                obj = new WoBufferFile(inode, position);
            } else if(flags & O_RDWR) {
                obj = new BufferFile(inode, position);
            } else {
                obj = new RoBufferFile(inode);
            }
            file = new File(obj, flags);
        }

        const fd = intf.fds.length;
        intf.fds[fd] = new Fd(file, (flags & O_CLOEXEC) ? FD_CLOEXEC : 0);
        return fd;
    }

    imports.env.__syscall_readlinkat = function __syscall_readlinkat(dirfd, ppath, _buf, _bufsize) {
        const path = resolveDirlessPath(intf, dirfd, ppath);
        if(typeof path === "number") {
            return path;
        }
        return -EINVAL;
    }
    imports.env.__syscall_getcwd = function __syscall_getcwd(buf: number, size: number) {
        if(size < 2)
            return -ERANGE;
        writeString(intf, buf, size, "/");
        return 0;
    }
    imports.env.__syscall_chmod = function __syscall_chmod(ppath: number, mode: number) {
        const path = resolveDirlessPath(intf, AT_FDCWD, ppath);
        if(typeof path === "number") {
            return path;
        }
        return 0;
    }
    imports.env.__syscall_unlinkat = function __syscall_unlinkat(dirfd: number, ppath: number, flags: number) {
        const path = resolveDirlessPath(intf, dirfd, ppath);
        if(typeof path === "number") {
            return path;
        }
        if(typeof path === "object") {
            return -EACCES;
        }
        if(!intf.fs.has(path)) {
            return -ENOENT;
        }
        if(flags & AT_REMOVEDIR) {
            return -ENOTDIR;
        }
        intf.fs.delete(path);
        return 0;
    }
    imports.env.__syscall_renameat = function __syscall_renameat(olddirfd: number, poldpath: ptr, newdirfd: number, pnewpath: ptr) {
        const oldpath = resolveDirlessPath(intf, olddirfd, poldpath);
        if(typeof oldpath === "number")
            return oldpath;
        const newpath = resolveDirlessPath(intf, newdirfd, pnewpath);
        if(typeof newpath === "number")
            return newpath;
        if(typeof oldpath === "object" || typeof newpath === "object") {
            return -EXDEV;
        }
        if(!intf.fs.has(oldpath)) {
            return -ENOENT;
        }
        intf.fs.set(newpath, intf.fs.get(oldpath)!);
        intf.fs.delete(oldpath);
        return 0;
    }
}

export function initAndCallMain(exports: WasmExports, args: string[]): number {
    const memory = exports.memory;
    const b = memory.buffer;
    const memU8 = new Uint8Array(b);
    const memU32 = new Uint32Array(b);
    const stackAlloc = exports._emscripten_stack_alloc;

    function stringToUTF8OnStack(str: string | Uint8Array | number[]): number {
        let ret: number;
        if (Array.isArray(str) || ArrayBuffer.isView(str)) {
            ret = stackAlloc(str.length);
            memU8.set(str as Uint8Array, ret);
        } else {
            const size = lengthBytesUTF8(str as string) + 1;
            ret = stackAlloc(size);
            stringToUTF8Array0(str as string, memU8, ret, size);
        }
        return ret;
    }

    const callCtors = exports.__wasm_call_ctors;
    callCtors();

    const entryFunction = exports.__main_argc_argv;
    const argc = args.length;
    const argv = stackAlloc((argc + 1) * 4);
    let argv_ptr = argv;
    args.forEach(arg => {
        memU32[argv_ptr >> 2] = stringToUTF8OnStack(arg);
        argv_ptr += 4;
    });
    // seems that an extra \0 is needed so that the last argument is not ignored
    stackAlloc(1);
    memU32[argv_ptr >> 2] = 0;
    try {
        const ret = entryFunction(argc, argv);
        return ret;
    } catch (e) {
        if (e instanceof ExitStatus || e == "unwind") {
            return (e as ExitStatus).status;
        }
        throw e;
    }
}
