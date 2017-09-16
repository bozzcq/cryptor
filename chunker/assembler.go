package chunker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

// AssembleData ...
func AssembleData(chunkDirPath string) (data []byte, err error) {
	var buffer bytes.Buffer

	// Get list of chunks
	chunks, err := ioutil.ReadDir(chunkDirPath)
	if err != nil {
		return nil, nil
	}

	// Go trough each chunk
	for _, chunk := range chunks {
		chunkPath := fmt.Sprintf("%s/%s", chunkDirPath, chunk.Name())

		if strings.HasSuffix(chunkPath, ppExtension) {
			continue
		}

		// Read content in each chunk
		chunkContent, err := ioutil.ReadFile(chunkPath)
		if err != nil {
			return nil, err
		}

		// Store content in a buffer
		buffer.Write(chunkContent)
	}

	// Read bytes from buffer
	packageBytes := buffer.Bytes()

	// Check for gzip compression
	if packageBytes[0] == 31 && packageBytes[1] == 139 {

		// Decompress package data
		packageData, err := decompress(&buffer)
		if err != nil {
			return nil, nil
		}
		return packageData, nil
	}

	// Return bytes as they are
	return packageBytes, nil
}

// AssembleFile ...
func AssembleFile(chunkDirPath string) (fileName string, err error) {
	return "", nil
}