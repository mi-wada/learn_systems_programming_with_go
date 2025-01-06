package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func main() {
	f, err := os.Open("test.3.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	chunks := chunks(f)

	// newFile, err := os.Create("test.3.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer newFile.Close()
	// io.WriteString(newFile, "\x89PNG\r\n\x1a\n")
	// io.Copy(newFile, chunks[0])
	// io.Copy(newFile, textChunk("Lambda Note++"))
	// for _, chunk := range chunks[1:] {
	// 	io.Copy(newFile, chunk)
	// }

	for _, chunk := range chunks {
		fmt.Println(dumpChunk(chunk))
	}
}

func textChunk(text string) io.Reader {
	byteText := []byte(text)
	crc := crc32.NewIEEE()
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, int32(len(byteText)))
	writer := io.MultiWriter(&buf, crc)
	io.WriteString(writer, "tEXt")
	writer.Write(byteText)
	binary.Write(&buf, binary.BigEndian, crc.Sum32())
	return &buf
}

func chunks(file *os.File) []io.Reader {
	var chunks []io.Reader
	// Skip first 8 bytes
	file.Seek(8, 0)

	var offset int64 = 8
	for {
		var length int32
		err := binary.Read(file, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		chunks = append(
			chunks,
			// 8(signature) + 4(data length) + length
			io.NewSectionReader(file, offset, 8+4+int64(length)),
		)
		// 4(data kind) + length + 4(CRC)
		offset, _ = file.Seek(int64(4+length+4), 1)
	}
	return chunks
}

func dumpChunk(chunk io.Reader) string {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length)

	kind := make([]byte, 4)
	chunk.Read(kind)

	res := fmt.Sprintf("chunk '%s' (%d bytes)", string(kind), length)
	if string(kind) == "tEXt" {
		text := make([]byte, length)
		chunk.Read(text)
		res += fmt.Sprintf(" (text: %s)", string(text))
	}
	return res
}
