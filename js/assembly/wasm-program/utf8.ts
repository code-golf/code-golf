
export function lengthBytesUTF8(str: string): number {
    let len = 0;
    for (let i = 0; i < str.length; ++i) {
        const c = str.charCodeAt(i);
        if (c <= 127) {
            len++;
        } else if (c <= 2047) {
            len += 2;
        } else if (c >= 55296 && c <= 57343) {
            len += 4;
            ++i;
        } else {
            len += 3;
        }
    }
    return len;
}

export function stringToUTF8Array(str: string, heap: Uint8Array, startIdx: number, maxBytesToWrite: number): number {
    if (maxBytesToWrite <= 0) return 0;
    let outIdx = startIdx;
    const limitIdx = startIdx + maxBytesToWrite;
    for (let i = 0; i < str.length; ++i) {
        let u = str.charCodeAt(i);
        if (u >= 55296 && u <= 57343) {
            const u1 = str.charCodeAt(++i);
            u = 65536 + ((u & 1023) << 10) | u1 & 1023;
        }
        if (u <= 127) {
            if (outIdx >= limitIdx) break;
            heap[outIdx++] = u;
        } else if (u <= 2047) {
            if (outIdx + 1 >= limitIdx) break;
            heap[outIdx++] = 192 | u >> 6;
            heap[outIdx++] = 128 | u & 63;
        } else if (u <= 65535) {
            if (outIdx + 2 >= limitIdx) break;
            heap[outIdx++] = 224 | u >> 12;
            heap[outIdx++] = 128 | u >> 6 & 63;
            heap[outIdx++] = 128 | u & 63;
        } else {
            if (outIdx + 3 >= limitIdx) break;
            heap[outIdx++] = 240 | u >> 18;
            heap[outIdx++] = 128 | u >> 12 & 63;
            heap[outIdx++] = 128 | u >> 6 & 63;
            heap[outIdx++] = 128 | u & 63;
        }
    }
    return outIdx - startIdx;
}

export function stringToUTF8Array0(str: string, heap: Uint8Array, startIdx: number, maxBytesToWrite: number): number {
    if (maxBytesToWrite <= 0) return 0;
    const endIdx = stringToUTF8Array(str, heap, startIdx, maxBytesToWrite - 1);
    heap[endIdx] = 0;
    return endIdx - startIdx;
}

export function stringToNewUTF8Array(str: string | Uint8Array | ArrayBufferView): Uint8Array {
    if (ArrayBuffer.isView(str)) {
        return str as Uint8Array;
    } else {
        const len = lengthBytesUTF8(str as string);
        const arr = new Uint8Array(len);
        stringToUTF8Array(str as string, arr, 0, len);
        return arr;
    }
}

const UTF8Decoder = typeof TextDecoder != "undefined" ? new TextDecoder() : undefined;

export function UTF8ArrayToString(heapOrArray: Uint8Array, idx: number = 0, maxBytesToRead: number = NaN): string {
    let endIdx = idx + maxBytesToRead;
    let endPtr = idx;
    while (heapOrArray[endPtr] && !(endPtr >= endIdx)) ++endPtr;
    if (endPtr - idx > 16 && heapOrArray.buffer && UTF8Decoder) {
        return UTF8Decoder.decode(heapOrArray.subarray(idx, endPtr));
    }
    let str = "";
    while (idx < endPtr) {
        let u0 = heapOrArray[idx++];
        if (!(u0 & 128)) {
            str += String.fromCharCode(u0);
            continue;
        }
        const u1 = heapOrArray[idx++] & 63;
        if ((u0 & 224) == 192) {
            str += String.fromCharCode((u0 & 31) << 6 | u1);
            continue;
        }
        const u2 = heapOrArray[idx++] & 63;
        if ((u0 & 240) == 224) {
            u0 = (u0 & 15) << 12 | u1 << 6 | u2;
        } else {
            u0 = (u0 & 7) << 18 | u1 << 12 | u2 << 6 | heapOrArray[idx++] & 63;
        }
        if (u0 < 65536) {
            str += String.fromCharCode(u0);
        } else {
            const ch = u0 - 65536;
            str += String.fromCharCode(55296 | ch >> 10, 56320 | ch & 1023);
        }
    }
    return str;
}
