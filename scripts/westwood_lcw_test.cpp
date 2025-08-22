#include <iostream>
#include <fstream>
#include <vector>
#include <cstring>
#include <cstdint>

// Westwood's original LCW decompression function
int LCW_Uncomp(void const * source, void * dest, unsigned long length)
{
	unsigned char * source_ptr, * dest_ptr, * copy_ptr, op_code, data;
	unsigned	  count, * word_dest_ptr, word_data;

	/* Copy the source and destination ptrs. */
	source_ptr = (unsigned char*) source;
	dest_ptr   = (unsigned char*) dest;

	while (1 /*TRUE*/) {

		/* Read in the operation code. */
		op_code = *source_ptr++;

		if (!(op_code & 0x80)) {

			/* Do a short copy from destination. */
			count	 = (op_code >> 4) + 3;
			copy_ptr = dest_ptr - ((unsigned) *source_ptr++ + (((unsigned) op_code & 0x0f) << 8));

			while (count--) *dest_ptr++ = *copy_ptr++;

		} else {

			if (!(op_code & 0x40)) {

				if (op_code == 0x80) {

					/* Return # of destination bytes written. */
					return ((unsigned long) (dest_ptr - (unsigned char*) dest));

				} else {

					/* Do a medium copy from source. */
					count = op_code & 0x3f;

					while (count--) *dest_ptr++ = *source_ptr++;
				}

			} else {

				if (op_code == 0xfe) {

					/* Do a long run. */
					count = *source_ptr + ((unsigned) *(source_ptr + 1) << 8);
					word_data = data = *(source_ptr + 2);
					word_data  = (word_data << 24) + (word_data << 16) + (word_data << 8) + word_data;
					source_ptr += 3;

					copy_ptr = dest_ptr + 4 - ((uintptr_t) dest_ptr & 0x3);
					count -= (copy_ptr - dest_ptr);
					while (dest_ptr < copy_ptr) *dest_ptr++ = data;

					word_dest_ptr = (unsigned*) dest_ptr;

					dest_ptr += (count & 0xfffffffc);

					while (word_dest_ptr < (unsigned*) dest_ptr) {
						*word_dest_ptr		= word_data;
						*(word_dest_ptr + 1) = word_data;
						word_dest_ptr += 2;
					}

					copy_ptr = dest_ptr + (count & 0x3);
					while (dest_ptr < copy_ptr) *dest_ptr++ = data;

				} else {

					if (op_code == 0xff) {

						/* Do a long copy from destination. */
						count	 = *source_ptr + ((unsigned) *(source_ptr + 1) << 8);
						copy_ptr = (unsigned char*) dest + *(source_ptr + 2) + ((unsigned) *(source_ptr + 3) << 8);
						source_ptr += 4;

						while (count--) *dest_ptr++ = *copy_ptr++;

					} else {

						/* Do a medium copy from destination. */
						count = (op_code & 0x3f) + 3;
						copy_ptr = (unsigned char*) dest + *source_ptr + ((unsigned) *(source_ptr + 1) << 8);
						source_ptr += 2;

						while (count--) *dest_ptr++ = *copy_ptr++;
					}
				}
			}
		}
	}
}

int main(int argc, char* argv[]) {
    if (argc != 3) {
        std::cerr << "Usage: " << argv[0] << " <input_cmp_file> <output_raw_file>" << std::endl;
        return 1;
    }

    const char* input_filename = argv[1];
    const char* output_filename = argv[2];

    // Read the input file
    std::ifstream input_file(input_filename, std::ios::binary);
    if (!input_file) {
        std::cerr << "Error: Cannot open input file " << input_filename << std::endl;
        return 1;
    }

    // Get file size
    input_file.seekg(0, std::ios::end);
    size_t input_size = input_file.tellg();
    input_file.seekg(0, std::ios::beg);

    // Read compressed data (skip CMP header - first 10 bytes)
    std::vector<unsigned char> compressed_data(input_size - 10);
    input_file.seekg(10); // Skip the 10-byte CMP header
    input_file.read(reinterpret_cast<char*>(compressed_data.data()), input_size - 10);
    input_file.close();

    std::cout << "Read " << compressed_data.size() << " bytes of compressed data" << std::endl;

    // Allocate output buffer (assume max 200KB should be enough)
    std::vector<unsigned char> output_buffer(200000, 0);

    // Decompress using Westwood's function
    int decompressed_size = LCW_Uncomp(compressed_data.data(), output_buffer.data(), 0);
    
    if (decompressed_size <= 0) {
        std::cerr << "Error: Decompression failed" << std::endl;
        return 1;
    }

    std::cout << "Decompressed to " << decompressed_size << " bytes" << std::endl;

    // Write output to file
    std::ofstream output_file(output_filename, std::ios::binary);
    if (!output_file) {
        std::cerr << "Error: Cannot create output file " << output_filename << std::endl;
        return 1;
    }

    output_file.write(reinterpret_cast<const char*>(output_buffer.data()), decompressed_size);
    output_file.close();

    std::cout << "Successfully wrote " << decompressed_size << " bytes to " << output_filename << std::endl;

    return 0;
}

// Compile with: g++ -o test_lcw westwood_lcw_test.cpp
// Run with: ./test_lcw EOBTITLE.CMP output.raw