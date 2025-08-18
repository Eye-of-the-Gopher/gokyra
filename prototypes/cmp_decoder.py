import os.path
import struct
import sys

def decode_cmp(cmpfile):
    with open(cmpfile, "rb") as f:
        header = f.read(10)
        data = struct.unpack('<HHLH', header)
        fsize, compression_type, uncompressed_size, palette_size = data
        print (f"Name : {cmpfile:32} | File size : {fsize:8} | Compression type : {compression_type} | Uncompressed size : {uncompressed_size} | Palette size : {palette_size}")


def main():
    for i in sys.argv[1:]:
        if not i.lower().endswith(".cmp"):
            print (f"Skipping {i}")
            continue
        decode_cmp(i)
        

if __name__ == '__main__':
    main()
