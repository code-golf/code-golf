#include <unistd.h>
#include <errno.h>
#include "asm.h"

// Function to read data handling EINTR and short reads
static ssize_t read_all(int fd, void *buf, size_t count) {
	size_t total_read = 0;
	ssize_t bytes_read;

	while (total_read < count) {
		bytes_read = read(fd, (char *)buf + total_read, count - total_read);

		if (bytes_read == 0) { // EOF
			break;
		}

		if (bytes_read < 0) {
			if (errno == EINTR) {
				continue; // Interrupted by signal, retry
			}
			return -1; // Error
		}

		total_read += (size_t)bytes_read;
	}

	return (ssize_t)total_read;
}

// Function to write data handling EINTR and short writes
static ssize_t write_all(int fd, const void *buf, size_t count) {
	size_t total_written = 0;
	ssize_t bytes_written;

	while (total_written < count) {
		bytes_written = write(fd, (const char *)buf + total_written, count - total_written);

		if (bytes_written < 0) {
			if (errno == EINTR) {
				continue; // Interrupted by signal, retry
			}
			return -1; // Error
		}

		total_written += (size_t)bytes_written;
	}

	return (ssize_t)total_written;
}

size_t read_fd(void* fd, char* buf, size_t len) {
	ssize_t bytes_read = read_all((int)(intptr_t)fd, buf, len);

	if (bytes_read < 0) {
		write_all(2, "error reading input\n", sizeof("error reading input\n"));
		_exit(1);
	}

	return (size_t)bytes_read;
}

void write_fd(void* fd, const char* buf, size_t len) {
	ssize_t written = write_all((int)(intptr_t)fd, buf, len);
	if (written < 0 || (size_t)written != len) {
		write_all(2, "out_flush failed\n", sizeof("out_flush failed\n"));
		_exit(1);
	}
}

int main(int argc, char *argv[]) {
	char buffer[4096 + 64];
	char out_buffer[4096];

	in_t in;
	out_t out;

	in_init(&in, buffer, 0, 4096, sizeof(buffer), read_fd, (void*)(intptr_t)STDIN_FILENO);
	out_init(&out, out_buffer, sizeof(out_buffer), write_fd, (void*)(intptr_t)STDOUT_FILENO);

	transform_assembly(&in, &out, false);
	return 0;
}
