export class ResizableBuffer {
    #array?: Uint8Array;
    #length: number;
    #capacity: number;
    readonly #growthFactor: number;

    constructor(array?: Uint8Array, growthFactor: number = 2) {
        if (growthFactor <= 1) {
            throw new Error("Growth factor must be greater than 1.");
        }

        this.#array = array;
        this.#capacity = this.array.length;
        this.#length = this.array.length;
        this.#growthFactor = growthFactor;
    }

    static withCapacity(initialCapacity: number = 4, growthFactor: number = 2) {
        if (initialCapacity < 0) {
            throw new Error("Initial capacity cannot be negative.");
        }
        return new ResizableBuffer(new Uint8Array(initialCapacity), growthFactor);
    }

    get array(): Uint8Array {
        return this.#array ? this.#array.subarray(0, this.#length) : new Uint8Array(0);
    }

    get length(): number {
        return this.#length;
    }

    get capacity(): number {
        return this.#capacity;
    }

    #ensureCapacity(minCapacity: number): void {
        if (this.#capacity < minCapacity) {
            let newCapacity = this.#capacity === 0 ? 1 : this.#capacity;
            while (newCapacity < minCapacity) {
                newCapacity = Math.ceil(newCapacity * this.#growthFactor);
            }

            const newBuffer = new Uint8Array(newCapacity);

            if(this.#array)
                newBuffer.set(this.#array.subarray(0, this.#length), 0);
            this.#array = newBuffer;
            this.#capacity = newCapacity;
        }
    }

    resize(newLength: number): void {
        if (newLength < 0) {
            throw new RangeError("New length cannot be negative.");
        }
        this.#ensureCapacity(newLength);
        this.#length = newLength;
    }

    clear(): void {
        this.#length = 0;
    }

    write(bufferWriteOffset: number, sourceArray: Uint8Array, sourceArrayOffset: number, count: number) {
        if (bufferWriteOffset < 0) {
            throw new RangeError("bufferWriteOffset cannot be negative.");
        }
        if (sourceArrayOffset < 0) {
            throw new RangeError("sourceArrayOffset cannot be negative.");
        }
        if (count < 0) {
            throw new Error("count cannot be negative.");
        }

        const bytesToCopy = Math.min(count, sourceArray.length - sourceArrayOffset);
        if (bytesToCopy <= 0) {
            return 0;
        }

        const requiredLength = bufferWriteOffset + bytesToCopy;
        this.#ensureCapacity(requiredLength);

        this.#array!.set(sourceArray.subarray(sourceArrayOffset, sourceArrayOffset + bytesToCopy), bufferWriteOffset);

        if (requiredLength > this.#length) {
            this.#length = requiredLength;
        }
        return bytesToCopy;
    }

    append(sourceArray: Uint8Array, sourceArrayOffset: number, count: number): number {
        return this.write(this.#length, sourceArray, sourceArrayOffset, count);
    }

    writeByte(bufferWriteOffset: number, value: number): void {
        if (bufferWriteOffset < 0) {
            throw new RangeError("bufferWriteOffset cannot be negative.");
        }

        const requiredLength = bufferWriteOffset + 1;
        this.#ensureCapacity(requiredLength);

        this.#array![bufferWriteOffset] = value;

        if (requiredLength > this.#length) {
            this.#length = requiredLength;
        }
    }

    push(value: number): void {
        this.#ensureCapacity(this.#length + 1);
        this.#array![this.#length] = value;
        ++this.#length;
    }

    read(bufferReadOffset: number, targetArray: Uint8Array, targetArrayOffset: number, count: number): number {
        if (bufferReadOffset < 0) {
            throw new RangeError("bufferReadOffset cannot be negative.");
        }
        if (targetArrayOffset < 0) {
            throw new RangeError("targetArrayOffset cannot be negative.");
        }
        if (count < 0) {
            throw new RangeError("count cannot be negative.");
        }

        const bytesToCopy = Math.min(count, this.#length - bufferReadOffset, targetArray.length - targetArrayOffset);

        if (bytesToCopy <= 0) {
            return 0;
        }

        targetArray.set(this.#array!.subarray(bufferReadOffset, bufferReadOffset + bytesToCopy), targetArrayOffset);
        return bytesToCopy;
    }

    trim(): void {
        if (this.#length < this.#capacity) {
            const newBuffer = new Uint8Array(this.#length);
            if (this.#length > 0) {
                newBuffer.set(this.#array!.subarray(0, this.#length), 0);
            }
            this.#array = newBuffer;
            this.#capacity = this.#length;
        }
    }
}
