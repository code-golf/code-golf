import { ptr } from "./mem";

export const EACCES = 2;
export const EBADF = 8;
export const EEXIST = 20;
export const EFAULT = 21;
export const EINVAL = 28;
export const EIO = 29;
export const EISDIR = 44;
export const ENOENT = 44;
export const ENOSYS = 52;
export const ENOTDIR = 54;
export const ENOTEMPTY = 55;
export const ENOTTY = 59;
export const ENXIO = 60;
export const EPERM = 63;
export const ERANGE = 68;
export const ESPIPE = 70;
export const EXDEV = 75;

export const O_RDONLY = 0;
export const O_WRONLY = 1;
export const O_RDWR = 2;

export const O_CREAT = 0o100;
export const O_EXCL = 0o200;
export const O_TRUNC = 0o1000;
export const O_APPEND = 0o2000;
export const O_DIRECTORY = 0o200000;
export const O_CLOEXEC = 0o02000000;

export const F_DUPFD = 0;
export const F_GETFD = 1;
export const F_SETFD = 2;
export const F_GETFL = 3;
export const F_SETFL = 4;

export const F_SETOWN = 8;
export const F_GETOWN = 9;
export const F_SETSIG = 10;
export const F_GETSIG = 11;

export const F_GETLK = 12;
export const F_SETLK = 13;
export const F_SETLKW = 14;
export const F_DUPFD_CLOEXEC = 1030;

export const FD_CLOEXEC = 1;

export const AT_FDCWD = -100;

export const AT_SYMLINK_NOFOLLOW = 0x100;
export const AT_REMOVEDIR = 0x200;
export const AT_SYMLINK_FOLLOW = 0x400;
export const AT_EMPTY_PATH = 0x1000;

export type WasmExports = {
    memory: WebAssembly.Memory;
    _emscripten_stack_alloc: (size: number) => number;
    emscripten_stack_get_current?: () => number;
    _emscripten_stack_restore?: (sp: number) => void;
    __wasm_call_ctors: () => void;
    __main_argc_argv: (argc: number, argv: number) => number;
    __indirect_function_table: WebAssembly.Table;
    _emscripten_timeout: (ms: number) => number;
}

export type WasmSyscalls = Partial<{
    __syscall_chmod: (path: ptr, mode: number) => void;
    __syscall_fcntl64: (fd: number, cmd: number, varargs: ptr) => number;
    __syscall_fstat64: (fd: number, buf: ptr) => void;
    __syscall_getcwd: (buf: ptr, size: number) => void;
    __syscall_ioctl: (fd: number, op: number, varargs: ptr) => number;
    __syscall_lstat64: (path: ptr, buf: ptr) => void;
    __syscall_newfstatat: (dirfd: number, path: ptr, buf: ptr, flags: number) => void;
    __syscall_openat: (dirfd: number, path: ptr, flags: number, varargs: ptr) => void;
    __syscall_readlinkat: (dirfd: number, path: ptr, buf: ptr, bufsize: number) => void;
    __syscall_stat64: (path: ptr, buf: ptr) => void;
    __syscall_unlinkat: (dirfd: number, path: ptr, flags: number) => void;
    __syscall_renameat: (olddirfd: number, oldpath: ptr, newdirfd: number, newpath: ptr) => void;
    __syscall_faccessat: (dirfd: number, path: ptr, amode: number, flags: number) => void;
    __syscall_dup: (fd: number) => void;
}>;

type EmscriptenApis = Partial<{
    __assert_fail: (condition: number, filename: ptr, line: number, func: ptr) => never;
    __call_sighandler: (fp: ptr, sig: number) => void;
    _abort_js: () => never;
    _emscripten_runtime_keepalive_clear: () => void;
    _localtime_js: (time: number, tmPtr: ptr) => void;
    _setitimer_js: (which: number, timeout_ms: number) => number;
    _tzset_js: (timezone: ptr, daylight: ptr, std_name: ptr, dst_name: ptr) => void;
    emscripten_date_now: () => number;
    emscripten_resize_heap: (requestedSize: number) => boolean;
    exit: (status: number) => never;
}>

export type WasiImports = Partial<{
    environ_get: (__environ: ptr, environ_buf: ptr) => number;
    environ_sizes_get: (penviron_count: ptr, penviron_buf_size: ptr) => number;
    fd_close: (fd: number) => number;
    fd_read: (fd: number, iov: ptr, iovcnt: number, pnum: ptr) => number;
    fd_pread: (fd: number, iov: ptr, iovcnt: number, offset: bigint, pnum: ptr) => number;
    fd_seek: (fd: number, offset: bigint, whence: number, pnewOffset: ptr) => number;
    fd_write: (fd: number, iov: ptr, iovcnt: number, pnum: ptr) => number;
    fd_pwrite: (fd: number, iov: ptr, iovcnt: number, offset: bigint, pnum: ptr) => number;
    fd_fdstat_get: (fd: number, buf: ptr) => number;
    proc_exit: (code: number) => never;
    clock_time_get: (clk_id: number, ignored_precision: bigint, ptime: ptr) => number;
}>

export type WasmImports = {
    env: EmscriptenApis & WasmSyscalls,
    wasi_snapshot_preview1: WasiImports,
}
