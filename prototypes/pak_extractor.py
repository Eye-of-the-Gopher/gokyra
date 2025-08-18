import os.path
import struct
import sys

def extract_pakfile(pakfile):
    fnames = []
    offsets = []
    
    with open(pakfile, "rb") as f:
        while True:

            if offsets and f.tell() +4 >= offsets[0]:
                break

            data = f.read(4)
            if len(data) == 4:
                value = struct.unpack('<I', data)[0] 
                offsets.append(value)
            else:
                print ("Short read : quitting")
                break

            data = f.read(1)
            if not data:
                print ("Short read : quitting")
                break

            fname = []
            while data != b'\0':
                value = struct.unpack('c', data)[0] 
                fname.append(value.decode('ascii'))
                # print(value.decode('ascii'), end="", flush=True)
                data = f.read(1)
                if not data:
                    print ("Short read : quitting")
                    break

            fname = "".join(fname)
            fnames.append(fname)
    
        
    
    base_dir = os.path.basename(pakfile)+"_pieces"
    if not os.path.exists(base_dir):
        os.makedirs(base_dir)

    with open(pakfile, "rb") as i:
        for idx, name in enumerate(fnames):
            print (" Writing ",name)
            start = offsets[idx]
            op = os.path.join(base_dir, name)
            with open(op, "wb") as o:
                try:
                    end = offsets[idx+1] 
                except Exception as e:
                    end = i.seek(0, 2)
                size = end - start
                i.seek(start)
                data = i.read(size)
                o.write(data)
        
        
def main(pakfiles):
    for i in pakfiles:
        if not i.lower().endswith(".pak"):
            print (f"Skipping {i}")
        print ("Extracting ", i)
        extract_pakfile(i)


if __name__ == '__main__':
    main(sys.argv[1:])
