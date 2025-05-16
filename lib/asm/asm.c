#include <stdint.h>
#include <stdio.h>

#include "asm.h"


#define memcpy __builtin_memcpy
#define memmove __builtin_memmove
#define memcmp __builtin_memcmp

static inline int isdigit(int c) {
	return c >= '0' && c <= '9';
}

static inline int isxdigit(int c) {
	return (c >= '0' && c <= '9') ||
		   (c >= 'a' && c <= 'f') ||
		   (c >= 'A' && c <= 'F');
}

static inline int tolower(int c) {
	if (c >= 'A' && c <= 'Z') {
		return c + ('a' - 'A');
	}
	return c;
}

// Flush output buffer to fd
static void out_flush(out_t *out) {
	if (out->len == 0)
		return;

	out->write(out->arg, out->buffer, out->len);
	out->len = 0;
}

// Write char to output buffer, flushing if needed
static inline void out_putc(out_t *out, char c) {
	// If the string is too large for the buffer, flush first
	if (out->len >= out->capacity) {
		out->write(out->arg, out->buffer, out->capacity);
		out->len = 0;
	}

	out->buffer[out->len++] = c;
}

// Write string to output buffer, flushing if needed
static void out_write(out_t *out, const char *str, size_t len) {
	size_t capacity_left = out->capacity - out->len;
	while (len > capacity_left) {
		memcpy(out->buffer + out->len, str, capacity_left);
		len -= capacity_left;
		str += capacity_left;

		out->write(out->arg, out->buffer, out->capacity);
		out->len = 0;
		capacity_left = out->capacity;
	}

	memcpy(out->buffer + out->len, str, len);
	out->len += len;
}

// Fill the input buffer with data from fd
static size_t in_fill(in_t *in) {
	size_t bytes_left = in->len - in->pos;
	if (bytes_left) {
		memmove(in->buffer, in->buffer + in->pos, bytes_left);
	}

	in->pos = 0;
	in->len = bytes_left;

	size_t read_len = in->read_len;
	if(read_len > in->capacity - bytes_left)
		read_len = in->capacity - bytes_left;

	size_t bytes_read = in->read(in->arg, in->buffer + in->len, read_len);

	in->len += bytes_read;
	return bytes_read;
}

// Get a character from the input buffer
static inline int in_getc(in_t *in) {
	if (in->pos >= in->len) {
		if (in_fill(in) == 0) {
			return EOF;
		}
	}

	return (uint8_t)in->buffer[in->pos++];
}

// Peek at the next character without advancing
static inline int in_peek(in_t *in) {
	if (in->pos >= in->len) {
		if (in_fill(in) == 0) {
			return EOF;
		}
	}

	return (uint8_t)in->buffer[in->pos];
}

// Peek at the next character without advancing
static inline size_t in_peekn(in_t *in, size_t len) {
	// Make sure we have enough data
	if (in->pos + len > in->len) {
		in_fill(in);

		size_t left = in->len - in->pos;
		if (left < len)
			return left;
	}

	return len;
}

// Look ahead to check if a sequence matches
static inline bool in_match(in_t *in, const char *seq, size_t len) {
	return in_peekn(in, len) == len && memcmp(in->buffer + in->pos, seq, len) == 0;
}

// Advance the buffer position by n characters
static inline void in_advance(in_t *in, size_t n) {
	in->pos += n;
}

// Parse a single character from a quoted string in a FILE*, handling all C escape sequences, return -1 if unescaped delimiter is read
static int parse_char_in_string(in_t *in, char out[16], char delimiter, uint64_t* in_newlines) {
	for(;;) {
		int c = 0;
		int i, digit, value;

		c = in_getc(in);
		if (c == EOF) return 0;

		if (c != '\\') {
			// Regular character
			if (c == delimiter) {
				return -1; // Signal unescaped delimiter
			}

			if(c == '\n')
				++*in_newlines;

			out[0] = (uint8_t)c;
			return 1;
		}

		// Save original backslash in case we need to return uninterpreted sequence
		out[0] = '\\';
		int pos = 1;

		// Handle escape sequence
		c = in_getc(in);
		if (c == EOF) return 1; // Return just the backslash if EOF after backslash

		switch (c) {
			case '\n': ++*in_newlines; break;
			case 'a': out[0] = '\a'; return 1;
			case 'b': out[0] = '\b'; return 1;
			case 'f': out[0] = '\f'; return 1;
			case 'n': out[0] = '\n'; return 1;
			case 'r': out[0] = '\r'; return 1;
			case 't': out[0] = '\t'; return 1;
			case 'v': out[0] = '\v'; return 1;
			case '\\': out[0] = '\\'; return 1;
			case '\'': out[0] = '\''; return 1;
			case '\"': out[0] = '\"'; return 1;
			case '?': out[0] = '\?'; return 1;

			case 'x': // Hexadecimal escape: \xHH
				out[pos++] = 'x';
				value = 0;
				for (i = 0; i < 2; i++) {
					c = in_peek(in);
					if (c == EOF) return pos; // Return partially parsed sequence

					if (isxdigit(c)) {
						in_advance(in, 1);
						out[pos++] = (uint8_t)c;
						digit = isdigit(c) ? c - '0' : tolower(c) - 'a' + 10;
						value = (value << 4) | digit;
					} else {
						return pos;
					}
				}

				if (i == 2) { // Successfully parsed two hex digits
					out[0] = (uint8_t)value;
					return 1;
				}
				return pos; // Return incomplete sequence

			case 'u': // Unicode escape: \uXXXX or \u{XXXXXX}
				out[pos++] = 'u';
				c = in_peek(in);
				if (c == EOF) return pos;

				uint32_t code_point = 0;

				if (c == '{') { // \u{XXXXXX} format
					in_advance(in, 1);
					out[pos++] = (uint8_t)c;

					code_point = 0;
					int num_xdigits = 0;
					while (1) {
						c = in_peek(in);
						if (c == EOF) return pos;

						if (c == '}') {
							if(num_xdigits == 0)
								return pos;
							in_advance(in, 1);
							out[pos++] = (uint8_t)c;
							break;
						}

						if (isxdigit(c)) {
							if(num_xdigits >= 8)
								return pos;

							in_advance(in, 1);
							out[pos++] = (uint8_t)c;

							digit = isdigit(c) ? c - '0' : tolower(c) - 'a' + 10;
							++num_xdigits;
							code_point = (code_point << 4) | digit;
							if (code_point > 0x10FFFF) {
								return pos; // Invalid Unicode code point
							}
						} else {
							return pos; // Error, return collected bytes
						}
					}
				} else { // \uXXXX format
					code_point = 0;

					for (i = 0; i < 4; i++) {
						c = in_peek(in);
						if (c == EOF) return pos;

						if (isxdigit(c)) {
							in_advance(in, 1);
							out[pos++] = (uint8_t)c;

							digit = isdigit(c) ? c - '0' : tolower(c) - 'a' + 10;
							code_point = (code_point << 4) | digit;
							if (code_point > 0x10FFFF) {
								return pos; // Invalid Unicode code point
							}
						} else {
							return pos; // Error, return collected bytes
						}
					}
				}

				// Convert Unicode code point to UTF-8
				if (code_point < 0x80) {
					out[0] = (uint8_t)code_point;
					return 1;
				} else if (code_point < 0x800) {
					out[0] = (uint8_t)(0xC0 | (code_point >> 6));
					out[1] = (uint8_t)(0x80 | (code_point & 0x3F));
					return 2;
				} else if (code_point < 0x10000) {
					out[0] = (uint8_t)(0xE0 | (code_point >> 12));
					out[1] = (uint8_t)(0x80 | ((code_point >> 6) & 0x3F));
					out[2] = (uint8_t)(0x80 | (code_point & 0x3F));
					return 3;
				} else {
					out[0] = (uint8_t)(0xF0 | (code_point >> 18));
					out[1] = (uint8_t)(0x80 | ((code_point >> 12) & 0x3F));
					out[2] = (uint8_t)(0x80 | ((code_point >> 6) & 0x3F));
					out[3] = (uint8_t)(0x80 | (code_point & 0x3F));
					return 4;
				}

			case '0': case '1': case '2': case '3':
			case '4': case '5': case '6': case '7': // Octal escape: \NNN
				out[pos++] = (uint8_t)c;
				value = c - '0';

				for (i = 1; i < 3; i++) {
					c = in_peek(in);
					if (c == EOF) return pos;

					if (c >= '0' && c <= '7') {
						in_advance(in, 1);
						out[pos++] = (uint8_t)c;

						value = (value << 3) | (c - '0');
					} else {
						return pos;
					}
				}

				out[0] = (uint8_t)value;
				return 1;

			default: // Not a recognized escape sequence, treat as literal
				out[0] = (uint8_t)c;
				return 1;
		}
	}
}

static void num_to_hex(uint64_t v, char* buf, size_t len) {
	for(int i = len - 1; i >= 0; --i) {
		uint8_t nibble = (uint8_t)(v & 0xf);
		v >>= 4;

		buf[i] = nibble >= 10 ? 'a' + (nibble - 10) : '0' + nibble;
	}
}

static void write_escaped(out_t* out, char* bytes, size_t len, char delimiter) {
	for(size_t i = 0; i < len; ++i) {
		uint8_t c = bytes[i];
		if(c == '\\' || c == delimiter)
			out_putc(out, '\\');

		if(c >= 32 && c <= 126)
			out_putc(out, c);
		else if(c == '\n')
			out_write(out, "\\n", 2);
		else if(c == '\r')
			out_write(out, "\\r", 2);
		else if(c == '\t')
			out_write(out, "\\t", 2);
		else {
			char buf[4];
			buf[0] = '\\';
			buf[1] = 'x';
			num_to_hex(c, buf + 2, 2);
			out_write(out, buf, 4);
		}
	}
}

// Process quoted characters and perform transformations to emulate DefAsm features
void transform_assembly(in_t *in, out_t *out, bool big_endian) {
	int ch;
	uint64_t in_newlines = 0;

	while ((ch = in_getc(in)) != EOF) {
		if (ch == '#') {
			while((ch = in_getc(in)) != EOF) {
				if(ch == '\n') {
					break;
				}
			}

			if(ch == EOF)
				break;
			// process '\n'
		}

		if (ch == '\'') {
			char buf[32];
			size_t total_bytes = 0;
			int read_bytes;
			uint64_t value = 0;
			int shift = 0;

			while((read_bytes = parse_char_in_string(in, buf, '\'', &in_newlines)) > 0) {
				for (size_t i = 0; i < (size_t)read_bytes; ++i) {
					if(total_bytes < 8) {
						if (big_endian) {
							value = (value << 8) | (uint8_t)buf[i];
						} else {
							value = value | ((uint64_t)(uint8_t)buf[i] << shift);
							shift += 8;
						}
						++total_bytes;
					}
				}
			}

			buf[0] = ' ';
			buf[1] = '0';
			buf[2] = 'x';
			num_to_hex(value, buf + 3, total_bytes * 2);
			buf[3 + total_bytes * 2] = ' ';
			out_write(out, buf, 4 + total_bytes * 2);
			if(read_bytes == 0) {
				// cause an error if we found an EOF in the string
				out_write(out, "+'", 2);
			}
		} else if(ch == '\"') {
			// Collect characters until closing quote
			char bytes[16];
			int read_bytes;

			out_putc(out, ch);
			while((read_bytes = parse_char_in_string(in, bytes, ch, &in_newlines)) > 0) {
				write_escaped(out, bytes, read_bytes, ch);
			}
			if(read_bytes < 0)
				out_putc(out, ch);
		} else if (ch == '.' && in_match(in, "ascii\"", 6)) {
			in_advance(in, 5);
			out_write(out, ".ascii ", 7);
		} else if (ch == '.' && in_match(in, "asciz\"", 6)) {
			in_advance(in, 5);
			out_write(out, ".asciz ", 7);
		} else if (ch == '.' && in_match(in, "string\"", 7)) {
			in_advance(in, 6);
			out_write(out, ".string ", 8);
		} else if(ch == '\n') {
			// Preserve line numbers by inserting newlines after a line for each newline in a string
			for(uint64_t i = 0; i <= in_newlines; ++i) {
				out_putc(out, '\n');
			}
			in_newlines = 0;
		} else {
			out_putc(out, ch);
		}
	}

	for(uint64_t i = 0; i < in_newlines; ++i) {
		out_putc(out, '\n');
	}
	out_flush(out);
}
