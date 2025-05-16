#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "asm.h"

// Define buffer sizes
#define INPUT_BUFFER_SIZE 4096
#define OUTPUT_BUFFER_SIZE 4096

#define TEST_BUFFER_SIZE 65536

// Static buffers for testing
static char test_input_buffer[TEST_BUFFER_SIZE];
static size_t test_input_pos = 0;
static size_t test_input_left = 0;

static char test_output_buffer[TEST_BUFFER_SIZE + 1];
static size_t test_output_pos = 0;

// Function to read from our test input buffer
static size_t test_read(void* arg, char* buffer, size_t len) {
    size_t to_read = len < test_input_left ? len : test_input_left;
    memcpy(buffer, test_input_buffer + test_input_pos, to_read);
    test_input_pos += to_read;
    test_input_left -= to_read;

    return to_read;
}

// Function to write to our test output buffer
static void test_write(void* arg, const char* buffer, size_t len) {
    if (test_output_pos + len >= TEST_BUFFER_SIZE) {
        fprintf(stderr, "Test output buffer overflow\n");
        exit(1);
    }

    memcpy(test_output_buffer + test_output_pos, buffer, len);
    test_output_pos += len;
    test_output_buffer[test_output_pos] = '\0'; // Ensure null termination
}

static uint64_t failures = 0;

// Test function
void test_endian(const char* name, const char* input, const char* expected_output, bool big_endian) {
    // Reset output buffer
    test_output_pos = 0;

    // Copy input to input buffer
    test_input_pos = 0;
    test_input_left = strlen(input);
    memcpy(test_input_buffer, input, test_input_left);

    in_t in;
    out_t out;

    char in_buffer[INPUT_BUFFER_SIZE + 64];
    char out_buffer[OUTPUT_BUFFER_SIZE];

    in_init(&in, in_buffer, 0, INPUT_BUFFER_SIZE, INPUT_BUFFER_SIZE + 64, test_read, 0);
    out_init(&out, out_buffer, OUTPUT_BUFFER_SIZE, test_write, 0);

    // Run the transformation
    transform_assembly(&in, &out, big_endian); // Use little-endian mode

    test_output_buffer[test_output_pos] = 0;

    // Check if the output matches the expected output
    if (strcmp(test_output_buffer, expected_output) != 0) {
        printf("Test failed: %s!\n", name);
        printf("Input   : %s\n", input);
        printf("Output  : %s\n", test_output_buffer);
        printf("Expected: %s\n", expected_output);


        // Print a detailed comparison for better debugging
        printf("\nDetailed comparison:\n");
        const char *e = expected_output;
        const char *g = test_output_buffer;
        int i = 0;
        while (*e != '\0' || *g != '\0') {
            if (*e != *g) {
                printf("Position %d: expected '%c' (0x%02x), got '%c' (0x%02x)\n",
                       i,
                       *e >= 32 && *e <= 126 ? *e : '.',
                       *e >= 0 ? (unsigned char)*e : 0,
                       *g >= 32 && *g <= 126 ? *g : '.',
                       *g >= 0 ? (unsigned char)*g : 0);
            }
            if (*e != '\0') e++;
            if (*g != '\0') g++;
            i++;
        }
        printf("\n");

        ++failures;
    }
}

void test(const char* name, const char* input, const char* expected_output) {
    test_endian(name, input, expected_output, false);
}

void test_be(const char* name, const char* input, const char* expected_output) {
    test_endian(name, input, expected_output, true);
}

int main() {
    test("Basic string literals",
        "mov rax, \"hello world\"",
        "mov rax, \"hello world\""
    );

    test("Escaped characters in string literals",
        "mov rax, \"hello\\nworld\\t!\"",
        "mov rax, \"hello\\nworld\\t!\""
    );

    test("Special ASCII commands",
        ".ascii\"Hello\"\n.asciz\"World\"\n.string\"Test\"",
        ".ascii \"Hello\"\n.asciz \"World\"\n.string \"Test\""
    );

    test("Character literals (little endian)",
        "mov rax, 'A'",
        "mov rax,  0x41 "
    );

    test("Multi-byte character literals (little endian)",
        "mov rax, 'AB'",
        "mov rax,  0x4241 "
    );

    test("Multi-byte character literals (up to 8 bytes, little endian)",
        "mov rax, 'ABCDEFGH'",
        "mov rax,  0x4847464544434241 "
    );

    test("Character literals exceeding 8 bytes (should truncate, little endian)",
        "mov rax, 'ABCDEFGHIJ'",
        "mov rax,  0x4847464544434241 "
    );

    test_be("Multi-byte character literals (big endian)",
        "mov rax, 'AB'",
        "mov rax,  0x4142 "
    );

    test_be("Multi-byte character literals (up to 8 bytes, big endian)",
        "mov rax, 'ABCDEFGH'",
        "mov rax,  0x4142434445464748 "
    );

    test_be("Character literals exceeding 8 bytes (should truncate, big endian)",
        "mov rax, 'ABCDEFGHIJ'",
        "mov rax,  0x4142434445464748 "
    );

    test("Escape sequences in character literals",
        "mov rax, '\\n'",
        "mov rax,  0x0a "
    );

    test("Hexadecimal escape sequences",
        "mov rax, '\\x41'",
        "mov rax,  0x41 "
    );

    test("Octal escape sequences",
        "mov rax, '\\101'",
        "mov rax,  0x41 "
    );

    test("Unicode escape sequences with \\u",
        "mov rax, '\\u0041'",
        "mov rax,  0x41 "
    );

    test("Unicode escape sequences with \\u{}",
        "mov rax, '\\u{41}'",
        "mov rax,  0x41 "
    );

    test("UTF-8 multi-byte encoding",
        "mov rax, '\\u{1F600}'", // Emoji grinning face
        "mov rax,  0x80989ff0 " // Little endian representation of UTF-8 bytes
    );

    test_be("UTF-8 multi-byte encoding (big endian)",
        "mov rax, '\\u{1F600}'", // Emoji grinning face
        "mov rax,  0xf09f9880 " // Big endian representation of UTF-8 bytes
    );

    test("Comments",
        "mov rax, 'A' # This is a comment\nmov rbx, 'B'",
        "mov rax,  0x41  \nmov rbx,  0x42 "
    );

    test("Incomplete string (EOF in string)",
        "mov rax, \"incomplete",
        "mov rax, \"incomplete"
    );

    test("Incomplete character literal (EOF in char literal)",
        "mov rax, 'incomplete",
        "mov rax,  0x656c706d6f636e69 +'"
    );

    test("Newlines inside string literals",
        "mov rax, \"hello\nworld\"",
        "mov rax, \"hello\\nworld\"\n"
    );

    test("All escape characters",
        "mov rax, '\\a\\b\\f\\n\\r\\t\\v\\\\'\nmov rax, '\\'\\\"\\?'",
        "mov rax,  0x5c0b090d0a0c0807 \nmov rax,  0x3f2227 "
    );


    test("Incomplete escape sequences",
        "mov rax, '\\x4'",
        "mov rax,  0x34785c "  // Interprets as '\' and '4'
    );

    test("Invalid Unicode code points",
        "mov rax, '\\u{110000}'",  // Beyond valid Unicode range
        "mov rax,  0x30303031317b755c " // Should just output the literal characters
    );

    test("Nested quotes in strings",
        "mov rax, \"He said \\\"Hello\\\"\"",
        "mov rax, \"He said \\\"Hello\\\"\""
    );

    test("All non-printable ASCII characters in string",
        "mov rax, \"\\x00\\x01\\x02\\x03\\x04\\x05\\x06\\x07\\x08\\x09\\x0A\\x0B\\x0C\\x0D\\x0E\\x0F\\x10\\x11\\x12\\x13\\x14\\x15\\x16\\x17\\x18\\x19\\x1A\\x1B\\x1C\\x1D\\x1E\\x1F\"",
        "mov rax, \"\\x00\\x01\\x02\\x03\\x04\\x05\\x06\\x07\\x08\\t\\n\\x0b\\x0c\\r\\x0e\\x0f\\x10\\x11\\x12\\x13\\x14\\x15\\x16\\x17\\x18\\x19\\x1a\\x1b\\x1c\\x1d\\x1e\\x1f\""
    );

    test("Consecutive strings",
        "mov rax, \"hello\"\"world\"",
        "mov rax, \"hello\"\"world\""
    );

    test("Consecutive character literals",
        "mov rax, 'A''B'",
        "mov rax,  0x41  0x42 "
    );

    test("Empty strings",
        "mov rax, \"\"",
        "mov rax, \"\""
    );

    test("Empty character literals",
        "mov rax, ''",
        "mov rax,  0x "
    );

    test("Multiple newlines in strings",
        "mov rax, \"hello\n\n\nworld\"",
        "mov rax, \"hello\\n\\n\\nworld\"\n\n\n"
    );

    test("Buffer boundary conditions - long strings",
        "mov rax, \"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\"",
        "mov rax, \"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\""
    );

    test("Buffer boundary conditions - long character literals",
        "mov rax, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA'",
        "mov rax,  0x4141414141414141 " // Truncated to 8 bytes
    );

    test("Mixing features",
        "mov rax, 'A'\nmov rbx, \"Hello\\nWorld\"\n.ascii\"Test\"# Comment",
        "mov rax,  0x41 \nmov rbx, \"Hello\\nWorld\"\n.ascii \"Test\""
    );

    test("Malformed Unicode",
        "mov rax, '\\u{ZZZZ}'",
        "mov rax,  0x7d5a5a5a5a7b755c "  // Should just output the literal characters
    );

    test("Unfinished Unicode",
        "mov rax, '\\u{41'",
        "mov rax,  0x31347b755c "  // Should output the literal characters
    );

    test("Complex nested case with various escapes",
        ".ascii\"Line1\\n\\\"Quoted\\\"\\n\"\n.asciz\"Line2\\n\\t\\u0041\"",
        ".ascii \"Line1\\n\\\"Quoted\\\"\\n\"\n.asciz \"Line2\\n\\tA\""
    );

    test("Edge case - backslash at the end of a string",
        "mov rax, \"backslash\\\\\"",
        "mov rax, \"backslash\\\\\""
    );

    test("Edge case - incomplete hex escape",
        "mov rax, \"\\x4\"",
        "mov rax, \"\\\\x4\""  // Doesn't get converted
    );

    test("Edge case - overflow in unicode handling",
        "mov rax, '\\u{FFFFFFFF}'",
        "mov rax,  0x46464646467b755c "  // Should just output the literal characters
    );

    // Buffer handling - input buffer refill
    // Generate a string larger than the buffer size to test in_fill function
    char large_input[4096 + 64];
    char large_output[4096 + 64];
    memset(large_input, 'A', 4093);
    strcpy(large_input + 4093, ".ascii\"foo");
    memset(large_output, 'A', 4093);
    strcpy(large_output + 4093, ".ascii \"foo");
    test("Large input/output with sequence crossing load boundary", large_input, large_output);

    // Output buffer flushing
    // Generate a string larger than the output buffer to test out_flush
    // (Same as above, reusing the large buffers)

    test("Mixed input with all features",
        "mov rax, 'A'    # Character literal\n"
        "mov rbx, \"Hello\\nWorld\"  # String literal\n"
        ".ascii\"Test\" # ASCII directive\n"
        "mov rcx, '\\u{1F600}'  # Unicode emoji\n"
        "mov rdx, '\\x41\\x42'  # Hex escapes\n"
        "# Comment line\n"
        ".string\"Final test\"",

        "mov rax,  0x41     \n"
        "mov rbx, \"Hello\\nWorld\"  \n"
        ".ascii \"Test\" \n"
        "mov rcx,  0x80989ff0   \n"
        "mov rdx,  0x4241   \n"
        "\n"
        ".string \"Final test\""
    );

    test("Character literals at beginning of line",
        "'A' is ASCII 65",
        " 0x41  is ASCII 65"
    );

    test("Character literals at end of line",
        "ASCII 65 is 'A'",
        "ASCII 65 is  0x41 "
    );

    test("Multiple character literals in one line",
        "Characters: 'A', 'B', 'C'",
        "Characters:  0x41 ,  0x42 ,  0x43 "
    );

    test("Character literals before and after strings",
        "'A' \"middle\" 'B'",
        " 0x41  \"middle\"  0x42 "
    );

    test("Character literals",
        "mov rax, 'A' + 'B'",
        "mov rax,  0x41  +  0x42 "
    );

    if(!failures) {
        printf("All tests succeeded!\n");
        return 0;
    } else {
        printf("%lu tests failed\n", failures);
        return 1;
    }
}