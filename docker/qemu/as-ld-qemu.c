#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <elf.h>
#include <errno.h>
#include <endian.h>
#include <dirent.h>
#include <ctype.h>
extern char **environ;

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

// Function to spawn a command without waiting
static pid_t spawn_command(const char *command, const char *const args[]) {
	pid_t pid = vfork();
	if (pid == -1) {
		perror("vfork failed");
		exit(1);
	}

	if (pid == 0) { // Child process
		execve(command, (char* const*)args, environ);
		perror("execve failed");
		_exit(1); // Use _exit with vfork
	}

	return pid;
}

// Function to wait for a command and check for errors
static void wait_command(pid_t pid) {
	int status;
	waitpid(pid, &status, 0);

	if (!WIFEXITED(status)) {
		fprintf(stderr, "Command did not exit normally\n");
		exit(1);
	}

	if (WEXITSTATUS(status) != 0) {
		fprintf(stderr, "Command failed with status %d\n", WEXITSTATUS(status));
		exit(1);
	}
}

// Function to execute a command and check for errors
static void execute_command(const char *command, const char *const args[]) {
	pid_t pid = spawn_command(command, args);
	wait_command(pid);
}

// Function to fix endianness for 16-bit ELF values
static inline uint16_t fix_endian16(uint16_t value, uint8_t ei_data) {
	if ((ei_data == ELFDATA2LSB && __BYTE_ORDER == __BIG_ENDIAN) ||
		(ei_data == ELFDATA2MSB && __BYTE_ORDER == __LITTLE_ENDIAN)) {
		return __builtin_bswap16(value);
	}
	return value;
}

// Function to fix endianness for 32-bit ELF values
static inline uint32_t fix_endian32(uint32_t value, uint8_t ei_data) {
	if ((ei_data == ELFDATA2LSB && __BYTE_ORDER == __BIG_ENDIAN) ||
		(ei_data == ELFDATA2MSB && __BYTE_ORDER == __LITTLE_ENDIAN)) {
		return __builtin_bswap32(value);
	}
	return value;
}

// Function to fix endianness for 64-bit ELF values
static inline uint64_t fix_endian64(uint64_t value, uint8_t ei_data) {
	if ((ei_data == ELFDATA2LSB && __BYTE_ORDER == __BIG_ENDIAN) ||
		(ei_data == ELFDATA2MSB && __BYTE_ORDER == __LITTLE_ENDIAN)) {
		return __builtin_bswap64(value);
	}
	return value;
}

static void read_elf(const char* filename, uint64_t* out_entry, uint64_t* out_bss_size)
{
	int fd = open(filename, O_RDONLY | O_CLOEXEC);
	if (fd == -1) {
		perror("Failed to open file");
		exit(1);
	}

	struct stat st;
	if (fstat(fd, &st) == -1) {
		perror("Failed to get file size");
		exit(1);
	}

	size_t file_size = (size_t)st.st_size;
	if((off_t)file_size != st.st_size) {
		fprintf(stderr, "ELF file too large\n");
		exit(1);
	}

	if(file_size < EI_NIDENT) {
		fprintf(stderr, "ELF file too short\n");
		exit(1);
	}

	void *mapped_file = mmap(NULL, file_size, PROT_READ, MAP_PRIVATE, fd, 0);
	if (mapped_file == MAP_FAILED) {
		perror("Failed to map file into memory");
		exit(1);
	}

	// Check ELF header
	if (memcmp(mapped_file, ELFMAG, SELFMAG) != 0) {
		fprintf(stderr, "Not a valid ELF file\n");
		exit(1);
	}

	uint8_t ei_class = ((const uint8_t*)mapped_file)[EI_CLASS];

	uint64_t bss_size = 0;
	uint64_t entry;

	if (ei_class == ELFCLASS64) {
		if(file_size < sizeof(Elf64_Ehdr)) {
			fprintf(stderr, "ELF file too short\n");
			exit(1);
		}

		Elf64_Ehdr *ehdr = (Elf64_Ehdr *)mapped_file;
		uint8_t ei_data = ehdr->e_ident[EI_DATA];

		// Check if program headers would exceed the file size
		// Multiplication can't overflow because e_phnum is 16-bit
		uint64_t phdr_size = (size_t)fix_endian16(ehdr->e_phnum, ei_data) * sizeof(Elf64_Phdr);
		uint64_t phoff = fix_endian64(ehdr->e_phoff, ei_data);
		if (phoff + phdr_size < phoff || phoff + phdr_size > file_size) {
			fprintf(stderr, "Invalid program headers offset and size\n");
			exit(1);
		}

		const Elf64_Phdr *phdr = (const Elf64_Phdr *)((const char *)mapped_file + phoff);

		// Iterate through program headers
		for (int i = 0; i < ehdr->e_phnum; i++) {
			// Pass ei_data to fix_endian functions
			Elf64_Word p_type = fix_endian32(phdr[i].p_type, ei_data);
			Elf64_Xword memsz = fix_endian64(phdr[i].p_memsz, ei_data);
			Elf64_Xword filesz = fix_endian64(phdr[i].p_filesz, ei_data);

			// Only count loadable segments
			if (p_type == PT_LOAD && memsz > filesz) {
				uint64_t new_bss_size = bss_size + memsz - filesz;
				if(new_bss_size < bss_size) {
					fprintf(stderr, "Integer overflow calculating BSS size\n");
					exit(1);
				}
				bss_size = new_bss_size;
			}
		}

		entry = fix_endian32(ehdr->e_entry, ei_data);
	} else if (ei_class == ELFCLASS32) {
		if(file_size < sizeof(Elf32_Ehdr)) {
			fprintf(stderr, "ELF file too short\n");
			exit(1);
		}

		Elf32_Ehdr *ehdr = (Elf32_Ehdr *)mapped_file;
		uint8_t ei_data = ehdr->e_ident[EI_DATA];

		// Check if program headers would exceed the file size
		// Multiplication can't overflow because e_phnum is 16-bit
		uint32_t phdr_size = (size_t)fix_endian16(ehdr->e_phnum, ei_data) * sizeof(Elf32_Phdr);
		uint32_t phoff = fix_endian32(ehdr->e_phoff, ei_data);
		if (phoff + phdr_size < phoff || phoff + phdr_size > file_size) {
			fprintf(stderr, "Invalid program headers offset and size\n");
			exit(1);
		}

		const Elf32_Phdr *phdr = (const Elf32_Phdr *)((const char *)mapped_file + phoff);

		// Iterate through program headers
		for (int i = 0; i < ehdr->e_phnum; i++) {
			// Pass ei_data to fix_endian functions
			Elf32_Word p_type = fix_endian32(phdr[i].p_type, ei_data);
			Elf32_Word memsz = fix_endian32(phdr[i].p_memsz, ei_data);
			Elf32_Word filesz = fix_endian32(phdr[i].p_filesz, ei_data);

			// Only count loadable segments
			if (p_type == PT_LOAD && memsz > filesz) {
				uint32_t new_bss_size = bss_size + memsz - filesz;
				if(new_bss_size < bss_size) {
					fprintf(stderr, "Integer overflow calculating BSS size\n");
					exit(1);
				}
				bss_size = new_bss_size;
			}
		}

		entry = fix_endian64(ehdr->e_entry, ei_data);
	} else {
		fprintf(stderr, "Unsupported ELF class\n");
		exit(1);
	}

	munmap(mapped_file, file_size);
	close(fd);

	*out_entry = entry;
	*out_bss_size = bss_size;
}

static void write_linker_script(const char* filename, bool single_phdr, uint64_t extra_bss_size) {
	char buffer[1024];
	int file_len = sprintf(buffer,
		single_phdr ?
		"ENTRY(_start)\n"
		"\n"
		"SECTIONS\n"
		"{\n"
		"\t. = 0x400000;\n"
		"\n"
		"\t.text : ALIGN(0x1000) {\n"
		"\t\t*(.text)\n"
		"\t\t*(.data)\n"
		"\t\t*(.rodata*)\n"
		"\t\t. = . + 0x%lx;\n"
		"\t\t*(.bss)\n"
		"\t} :text\n"
		"}\n"
		"\n"
		"PHDRS\n"
		"{\n"
		"\ttext PT_LOAD FLAGS(7) AT(0x400000);\n"
		"}\n"
		:
		"ENTRY(_start)\n"
		"\n"
		"SECTIONS\n"
		"{\n"
		"\t. = 0x400000;\n"
		"\n"
		"\t.text : ALIGN(0x1000) {\n"
		"\t\t*(.text)\n"
		"\t} :text\n"
		"\n"
		"\t.data : ALIGN(0x1000) {\n"
		"\t\t*(.data)\n"
		"\t\t*(.rodata*)\n"
		"\t\t*(.*)\n"
		"\t} :data\n"
		"\n"
		"\t.bss : ALIGN(0x1000) {\n"
		"\t\t. = . + 0x%lx;\n"
		"\t\t*(.bss)\n"
		"\t} :bss\n"
		"}\n"
		"\n"
		"PHDRS\n"
		"{\n"
		"\ttext PT_LOAD FLAGS(7) AT(0x400000);\n"
		"\tdata PT_LOAD FLAGS(7);\n"
		"\tbss PT_LOAD FLAGS(7);\n"
		"}\n", extra_bss_size);

	if(file_len < 0) {
		fprintf(stderr, "Failed to format linked script");
		exit(1);
	}

	int fd = open(filename, O_WRONLY | O_CREAT | O_TRUNC | O_CLOEXEC, 0644);
	if(fd < 0) {
		perror("open linker script");
		exit(1);
	}

	if(write_all(fd, buffer, file_len) != file_len) {
		perror("write linker script");
		exit(1);
	}

	if(close(fd) < 0) {
		perror("close linker script");
		exit(1);
	}
}

// Function to count non-zero bytes in a file, ignoring "sparse" page tails
static uint64_t count_nonzero_bytes(const char *filename) {
	int fd = open(filename, O_RDONLY | O_CLOEXEC);
	if (fd == -1) {
		perror("Failed to open file for counting non-zero bytes");
		return -1;
	}

	const size_t PAGE_SIZE = 4096;
	uint8_t buffer[PAGE_SIZE];
	uint64_t total_count = 0;
	ssize_t bytes_read;

	while ((bytes_read = read_all(fd, buffer, PAGE_SIZE)) > 0) {
		int last_nonzero = -1;

		for (int i = 0; i < bytes_read; i++) {
			if (buffer[i] != 0) {
				last_nonzero = i;
			}
		}

		if (last_nonzero >= 0) {
			total_count += last_nonzero + 1;
		}
	}

	if(bytes_read < 0) {
		perror("read from binary");
		exit(1);
	}

	close(fd);
	return total_count;
}

size_t read_fd(void* fd, char* buf, size_t len) {
	ssize_t bytes_read = read_all((int)(intptr_t)fd, buf, len);

	if (bytes_read < 0) {
		perror("reading input");
		exit(1);
	}

	return (size_t)bytes_read;
}

void write_fd(void* fd, const char* buf, size_t len) {
	ssize_t written = write_all((int)(intptr_t)fd, buf, len);
	if (written < 0 || (size_t)written != len) {
		fprintf(stderr, "out_flush failed");
		exit(1);
	}
}

int main(int argc, char *argv[]) {
	fcntl(3, F_SETFD, FD_CLOEXEC);
	chdir("/tmp");

	if (argc < 2) {
		fprintf(stderr, "Usage: %s <architecture>\n", argv[0]);
		return 1;
	}

	// Process architecture and set variables
	const char *prefix;
	const char *machine;
	const char *output_format;
	char *as_extra_args[16] = {NULL}; // Static array with NULL termination
	int as_extra_arg_count = 0;
	int bits = 0;
	const char *qemu_flavor32 = NULL;
	const char *qemu_flavor64 = NULL;
	char march[2048];
	char qemu_cpu[512];

	char buffer[4096 + 64];
	ssize_t bytes_read;

	if((bytes_read = read_all(STDIN_FILENO, buffer, sizeof(buffer))) < 0) {
		perror("Error reading from stdin");
		exit(1);
	}

	//const char* no_bits_error = "Error: file must start with #32 or #64 to select the architecture flavor\n";
	if (bytes_read >= 3 && memcmp(buffer, "#64", 3) == 0) {
		bits = 64;
	} else if (bytes_read >= 3 && memcmp(buffer, "#32", 3) == 0) {
		bits = 32;
	}

	strcpy(march, "-march=");

	if (strcmp(argv[1], "risc-v") == 0) {
		prefix = "riscv64-linux-gnu-";
		qemu_flavor64 = "riscv64";
		qemu_flavor32 = "riscv32";

		if (bits == 32) {
			strcat(march, "rv32");
			machine = "elf32lriscv";
			output_format = "elf32-littleriscv";
		} else { //if (bits == 64) {
			strcat(march, "rv64");
			machine = "elf64lriscv";
			output_format = "elf64-littleriscv";
		}
		/*
			else {
			fprintf(stderr, no_bits_error);
			exit(1);
		}
		*/

		bool prefer_m = false;

		if(bits) {
			for(unsigned i = 3; i < bytes_read; ++i) {
				char c = buffer[i];
				if(c == '\n')
					break;
				if(c == 'M')
					prefer_m = true;
				else {
					fprintf(stderr, "Unsupported character on architecture flavor line. Syntax is #(32|64)M? where M specifies to enable Zcmp+Zcmt instead of Zcd");
					exit(1);
				}
			}
		} else {
			bits = 64;
		}

		strcpy(qemu_cpu, "max,xtheadba=true,xtheadbb=true,xtheadbs=true,xtheadcmo=true,xtheadcondmov=true,xtheadfmemidx=true,xtheadfmv=true,xtheadmac=true,xtheadmemidx=true,xtheadmempair=true,xtheadsync=true,xventanacondops=true");
		strcat(march, "gmafdqcbvh_zic64b_ziccamoa_ziccif_zicclsm_ziccrse_zicbom_zicbop_zicboz_zicond_zicntr_zicsr_zifencei_zihintntl_zihintpause_zihpm_zimop_zicfiss_zicfilp_zmmul_za64rs_za128rs_zaamo_zabha_zacas_zalrsc_zawrs_zfbfmin_zfa_zfh_zfhmin_zbb_zba_zbc_zbs_zbkb_zbkc_zbkx_zk_zkn_zknd_zkne_zknh_zkr_zks_zksed_zksh_zkt_zve32x_zve32f_zve64x_zve64f_zve64d_zvbb_zvbc_zvfbfmin_zvfbfwma_zvfh_zvfhmin_zvkb_zvkg_zvkn_zvkng_zvknc_zvkned_zvknha_zvknhb_zvksed_zvksh_zvks_zvksg_zvksc_zvkt_zvl32b_zvl64b_zvl128b_zvl256b_zvl512b_zvl1024b_zvl2048b_zvl4096b_zvl8192b_zvl16384b_zvl32768b_zvl65536b_ztso_zca_zcb_zcmop_xtheadba_xtheadbb_xtheadbs_xtheadcmo_xtheadcondmov_xtheadfmemidx_xtheadfmv_xtheadmac_xtheadmemidx_xtheadmempair_xtheadsync_xventanacondops");

		if(prefer_m) {
			strcat(qemu_cpu, ",zcd=false,zce=true,zcmp=true,zcmt=true");
			strcat(march, "_zcmp_zcmt");
		} else {
			strcat(qemu_cpu, ",zcd=true");
			strcat(march, "_zcd");
		}

		as_extra_args[0] = "-mno-arch-attr";
		as_extra_arg_count = 1;
	} else if (strcmp(argv[1], "arm64") == 0) {
		bits = 64;
		prefix = "aarch64-linux-gnu-";
		qemu_flavor64 = "aarch64";

		strcpy(qemu_cpu, "max");
		strcat(march, "armv9.5-a+fp+bf16+crypto+crc+f32mm+f64mm+fp16+fp16fml+memtag+rng+sb+simd+sme+sme-f64f64+sme-i16i64+sve+sve2");
		machine = "aarch64elf";
		output_format = "elf64-littleaarch64";
	} else {
		fprintf(stderr, "Error: unknown architecture\n");
		return 1;
	}

	// Create pipe for assembly code just before spawning the assembler
	int pipe_fds[2];
	if (pipe(pipe_fds) == -1) {
		perror("Failed to create pipe");
		return 1;
	}

	// Get pipe read and write ends
	int pipe_read_fd = pipe_fds[0];
	int pipe_write_fd = pipe_fds[1];

	// Set write end to close-on-exec
	fcntl(pipe_write_fd, F_SETFD, FD_CLOEXEC);

	// Create /proc/self/fd path for the pipe's read end
	char fd_path[64];
	sprintf(fd_path, "/proc/self/fd/%d", pipe_read_fd);

	// 1. Assemble code
	char as_cmd[256];
	sprintf(as_cmd, "/usr/bin/%sas", prefix);

	// Calculate needed size for as_args
	const char *as_args[32]; // Static array large enough for all arguments
	int arg_index = 0;

	// Add the command itself
	as_args[arg_index++] = as_cmd;

	// Add all extra_args if present
	for (int i = 0; i < as_extra_arg_count; i++) {
		as_args[arg_index++] = as_extra_args[i];
	}

	// Add other arguments
	as_args[arg_index++] = march;
	as_args[arg_index++] = fd_path;
	as_args[arg_index++] = "-o";
	as_args[arg_index++] = "asm.o";
	as_args[arg_index] = NULL;

	// Spawn the assembler process
	pid_t as_pid = spawn_command(as_cmd, as_args);
	close(pipe_read_fd);

	// Write the rest of stdin to the pipe
	char out_buffer[4096];

	// Initialize buffer structures on the stack
	in_t in;
	out_t out;

	in_init(&in, buffer, bytes_read, 4096, sizeof(buffer), read_fd, (void*)(intptr_t)STDIN_FILENO);
	out_init(&out, out_buffer, sizeof(out_buffer), write_fd, (void*)(intptr_t)pipe_write_fd);

	// Process the data
	transform_assembly(&in, &out, false);

	// Close write end of pipe after writing all data
	close(pipe_write_fd);

	// Now wait for the assembler
	wait_command(as_pid);

	// 2. Link initial ELF executable
	char ld_cmd[256];
	sprintf(ld_cmd, "/usr/bin/%sld", prefix);
	write_linker_script("asm.elf.lds", false, 0);
	const char * const ld_args[] = {ld_cmd, "-m", machine, "-T", "asm.elf.lds", "--no-warn-rwx", "asm.o", "-o", "asm.elf", NULL};
	execute_command(ld_cmd, ld_args);

	// 3. Get entry point directly from ELF header
	uint64_t entry_point;
	uint64_t bss_size;
	read_elf("asm.elf", &entry_point, &bss_size);
	char entry_point_str[32];
	sprintf(entry_point_str, "0x%lx", entry_point);

	// 4. Flatten executable to binary
	char objcopy_cmd[256];
	sprintf(objcopy_cmd, "/usr/bin/%sobjcopy", prefix);
	const char * const objcopy_args1[] = {objcopy_cmd, "-O", "binary", "-S", "asm.elf", "asm.bin", NULL};
	execute_command(objcopy_cmd, objcopy_args1);

	// 5. Count non-zero bytes in code.bin and write to fd 3
	long nonzero_count = count_nonzero_bytes("asm.bin");
	char count_str[32];
	sprintf(count_str, "%ld\n", nonzero_count);
	write(3, count_str, strlen(count_str));
	close(3); // Close fd 3 after writing to it

	// 6. Convert back flat binary to ELF object
	const char * const objcopy_args2[] = {objcopy_cmd, "-I", "binary", "-O", output_format,
							 "--rename-section", ".data=.text", "asm.bin", "asm.bin.o", NULL};
	execute_command(objcopy_cmd, objcopy_args2);

	// 7. Link final ELF executable
	const char *ld_args2[] = {ld_cmd, "-m", machine, "-T", "asm.lds", "-e", entry_point_str, "--no-warn-rwx",
						"asm.bin.o", "-o", "/tmp/asm", NULL};
	write_linker_script("asm.lds", true, bss_size);
	execute_command(ld_cmd, ld_args2);

	char qemu32_path[256] = {0};
	sprintf(qemu32_path, "/usr/bin/qemu-%s", qemu_flavor32);

	char qemu64_path[256] = {0};
	sprintf(qemu64_path, "/usr/bin/qemu-%s", qemu_flavor64);

	const char* qemu_path = bits == 32 ? qemu32_path : qemu64_path;

	int qemu_fd = open(qemu_path, O_RDONLY | O_CLOEXEC);
	if(qemu_fd < 0) {
		perror("Failed to open QEMU executable");
		exit(1);
	}

#ifndef KEEP_FILES
	const char *files[] = {
		"asm.o", "asm.elf.lds", "asm.elf",
		"asm.bin", "asm.bin.o", "asm.lds",
		//as_cmd, ld_cmd, objcopy_cmd, qemu32_path, qemu64_path, argv[0]
	};

	for (size_t i = 0; i < sizeof(files) / sizeof(files[0]); ++i) {
		if(files[i][0] && unlink(files[i]) < 0) {
			perror("unlink");
			exit(1);
		}
	}
#endif

	unsigned base_qemu_argc = 4;
	const char** qemu_argv = malloc(sizeof(char*) * (argc + base_qemu_argc - 2 + 1));
	qemu_argv[0] = qemu_path;
	qemu_argv[1] = "-cpu";
	qemu_argv[2] = qemu_cpu;
	qemu_argv[3] = "/tmp/asm";
	for(int i = 2; i <= argc; ++i) {
		qemu_argv[base_qemu_argc + (i - 2)] = argv[i];
	}
	fexecve(qemu_fd, (char* const*)qemu_argv, environ);

	// If we get here, execution failed
	perror("Failed to execute code with QEMU");
	return 1;
}
