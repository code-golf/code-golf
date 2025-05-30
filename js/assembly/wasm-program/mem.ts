import { stringToUTF8Array0, UTF8ArrayToString } from "./utf8";

export interface MemInterface {
    memU8: Uint8Array
}

export type ptr = number;

export function read8(intf: MemInterface, ptr: ptr): number {
    return intf.memU8[ptr];
}

export function write8(intf: MemInterface, ptr: ptr, value: number) {
    intf.memU8[ptr] = value & 0xff;
}

export function bzero(intf: MemInterface, ptr: ptr, size: number) {
    for(let i = 0; i < size; ++i) {
        write8(intf, ptr + i, 0);
    }
}

// using UInt32Array is incorrect because wasm is little-endian while UInt32Array is native endian, and also it doesn't work for unaligned pointers
export function read16(intf: MemInterface, ptr: ptr): number {
    return intf.memU8[ptr] | (intf.memU8[ptr + 1] << 8);
}

export function write16(intf: MemInterface, ptr: ptr, value: number) {
    intf.memU8[ptr] = value & 0xff;
    intf.memU8[ptr + 1] = (value >> 8) & 0xff;
}

// using UInt32Array is incorrect because wasm is little-endian while UInt32Array is native endian, and also it doesn't work for unaligned pointers
export function read32(intf: MemInterface, ptr: ptr): number {
    return intf.memU8[ptr] | (intf.memU8[ptr + 1] << 8) | (intf.memU8[ptr + 2] << 16) | (intf.memU8[ptr + 3] << 24);
}

export function write32(intf: MemInterface, ptr: ptr, value: number) {
    intf.memU8[ptr] = value & 0xff;
    intf.memU8[ptr + 1] = (value >> 8) & 0xff;
    intf.memU8[ptr + 2] = (value >> 16) & 0xff;
    intf.memU8[ptr + 3] = (value >> 24) & 0xff;
}

export function write64(intf: MemInterface, ptr: ptr, value: number) {
    intf.memU8[ptr] = value & 0xff;
    intf.memU8[ptr + 1] = (value >> 8) & 0xff;
    intf.memU8[ptr + 2] = (value >> 16) & 0xff;
    intf.memU8[ptr + 3] = (value >> 24) & 0xff;
    intf.memU8[ptr + 4] = (value >> 32) & 0xff;
    intf.memU8[ptr + 5] = (value >> 40) & 0xff;
    intf.memU8[ptr + 6] = (value >> 48) & 0xff;
    intf.memU8[ptr + 7] = (value >> 56) & 0xff;
}

export function readString(intf: MemInterface, ptr: ptr, maxBytesToRead?: number): string {
    return ptr ? UTF8ArrayToString(intf.memU8, ptr, maxBytesToRead) as string : "";
}

export function writeString(intf: MemInterface, ptr: ptr, size: number, str: string) {
    stringToUTF8Array0(str, intf.memU8, ptr, size);
}

export function writeSyscallStat(intf: MemInterface, buf: ptr, size: number | undefined) {
    bzero(intf, buf, 96);

    const dev = size !== undefined ? 0x101 : 13;
    const blksize = 4096;
    const nlink = 1;
    const mode = (size !== undefined ? 0o0100000 : 0o010000) | 0o700;
    write32(intf, buf, dev);
    write32(intf, buf + 4, mode);
    write32(intf, buf + 8, nlink);
    write64(intf, buf + 24, size !== undefined ? size : 0);
    write32(intf, buf + 32, blksize);
    write32(intf, buf + 36, size !== undefined ? Math.ceil(size / blksize) : 0);
    return 0
}

export function writeWasiStat(intf: MemInterface, buf: ptr, size: number | undefined) {
    write8(intf, buf, size !== undefined ? 4 : 2);
    bzero(intf, buf + 1, 23);
}
