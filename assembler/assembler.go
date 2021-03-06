// Package assembler contains data structures and functions that
// read, decrypt and reassembel the original data files from the chunks.
package assembler

import (
	"bytes"

	"github.com/thee-engineer/cryptor/archive"
	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
)

// DefaultAssembler ...
type DefaultAssembler struct {
	Tail  []byte
	Cache cachedb.Database
}

// NewDefaultAssembler ...
func NewDefaultAssembler(tail []byte, cache cachedb.Database) Assembler {
	return &DefaultAssembler{
		Tail:  tail,
		Cache: cache,
	}
}

// getChunk returns an encrypted chunk from the attached cache.
func (a *DefaultAssembler) getChunk(hash []byte) (EChunk, error) {
	eChunk, err := a.Cache.Get(hash)
	if err != nil {
		return nil, err
	}
	return eChunk, nil
}

// Assemble starts decrypting the tail chunk with the given AES Key. The
// process extracts the next chunk's data from the current header. If a chunk
// is not found during the assembly process, a network request will be sent
// to the known peers, asking for the missing chunk.
func (a *DefaultAssembler) Assemble(key aes.Key, destination string) error {
	var cBuffer bytes.Buffer // Chunk buffer, content (no header)
	var aBuffer bytes.Buffer // Assembly buffer, final package

	// Memory zeroing
	defer crypt.ZeroBytes(key[:])

	// Request chunk from cache
	eChunk, err := a.getChunk(a.Tail)
	if err != nil {
		return err
	}

	// Decrypt given chunk with given key
	chunk, err := eChunk.Decrypt(key)
	if err != nil {
		return err
	}

	// Store decrypted chunk (including header)
	cBuffer.Write(chunk.Content)

	// Process chunks until a final chunk is passed
	for !chunk.IsLast() {
		// Get the next chunk by using the header.Next hash
		eChunk, err = a.getChunk(chunk.Header.Next)
		if err != nil {
			return err
		}

		// Decrypt the next chunk
		chunk, err = eChunk.Decrypt(chunk.Header.NKey)
		if err != nil {
			return err
		}

		// Store chunk content
		cBuffer.Write(chunk.Content)
	}

	// Fix single chunk size error being 0 by adding chunk header padding
	chunkSize := len(chunk.Content) + int(chunk.Header.Padd)
	bufferLen := cBuffer.Len()
	bufferData := cBuffer.Bytes()

	// Walk trough all processed chunks, place the chunks in the right order
	// inside the assembly buffer
	for index := bufferLen; index > chunkSize; index -= chunkSize {
		aBuffer.Write(bufferData[index-chunkSize : index])
	}

	// Write the final chunk content
	aBuffer.Write(bufferData[0 : bufferLen%chunkSize])

	// Start extracting the .tar.gz archive
	err = archive.UnTarGz(destination, &aBuffer)
	if err != nil {
		return err
	}

	return nil
}
