#include <stdbool.h>
#include <stddef.h>

// Output buffer structure
typedef struct {
	char *buffer;
	size_t capacity;
	size_t len;
	void (*write)(void*, const char*, size_t);
	void* arg;
} out_t;

// Initialize output buffer
static inline void out_init(out_t *out, char *buffer, size_t capacity, void (*write)(void*, const char*, size_t), void* arg) {
	out->buffer = buffer;
	out->capacity = capacity;
	out->len = 0;
	out->write = write;
	out->arg = arg;
}

// Input buffer structure
typedef struct {
	char *buffer;
	size_t pos;
	size_t len;
	size_t read_len;
	size_t capacity;
	size_t (*read)(void*, char*, size_t);
	void* arg;
} in_t;

// Initialize input buffer
static inline void in_init(in_t *in, char *buffer, size_t len, size_t read_len, size_t capacity, size_t (*read)(void*, char*, size_t), void* arg) {
	in->buffer = buffer;
	in->capacity = capacity;
	in->read_len = read_len;
	in->len = len;
	in->pos = 0;
	in->read = read;
	in->arg = arg;
}

void transform_assembly(in_t *in, out_t *out, bool big_endian);
