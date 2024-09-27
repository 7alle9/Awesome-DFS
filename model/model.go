package model

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

type Chunk struct {
	Id          int
	UniqueName  string
	FileId      int
	ChunkNumber int64
	RawData     []byte
	Size        int64
	CheckSum    string
}

type File struct {
	Id           int
	Name         string
	RawData      []byte
	Size         int64
	CheckSum     string
	LastModified time.Time
	Chunks       []*Chunk
}

func FetchLocalFile(filename string) (*File, error) {
	// Open the file
	rawFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Get file info
	info, err := rawFile.Stat()
	if err != nil {
		return nil, err
	}

	// Get file metadata
	name := info.Name()
	size := info.Size()
	lastModified := info.ModTime()

	// Read the file
	rawData := make([]byte, size)
	_, err = rawFile.Read(rawData)
	if err != nil {
		return nil, err
	}

	// Calculate checksum
	checkSum, err := getCheckSum(rawData)
	if err != nil {
		return nil, err
	}

	// Close the file
	if err := rawFile.Close(); err != nil {
		return nil, err
	}

	return &File{
		Name:         name,
		RawData:      rawData,
		Size:         size,
		CheckSum:     checkSum,
		LastModified: lastModified,
	}, nil
}

func (f *File) WriteFile() error {
	if err := os.WriteFile(f.Name, f.RawData, 0644); err != nil {
		return err
	}
	return nil
}

func (f *File) SplitFile(offset int64) error {
	// Check if offset is valid
	if offset <= 0 {
		return fmt.Errorf("invalid offset")
	}

	// Split the file into chunks
	var chunks []*Chunk
	for i := int64(0); i < f.Size; i = i + offset {
		// Calculate the chunk size
		var chunkSize int64
		if i+offset > f.Size {
			chunkSize = f.Size - i
		} else {
			chunkSize = offset
		}

		// Get the chunk data
		chunkData := f.RawData[i : i+chunkSize]

		// Calculate the checksum
		checkSum, err := getCheckSum(chunkData)
		if err != nil {
			return err
		}

		// Create the chunk
		chunk := &Chunk{
			ChunkNumber: i / offset,
			RawData:     chunkData,
			Size:        chunkSize,
			CheckSum:    checkSum,
		}
		chunks = append(chunks, chunk)
	}

	// Update the file data
	f.Chunks = chunks

	return nil
}

func (f *File) ReconstructFile() error {
	// Reconstruct the file
	var rawData []byte
	for _, chunk := range f.Chunks {
		rawData = append(rawData, chunk.RawData...)
	}

	// Calculate the checksum
	checkSum, err := getCheckSum(rawData)
	if err != nil {
		return err
	}
	if checkSum != f.CheckSum || int64(len(rawData)) != f.Size {
		return fmt.Errorf("data discrepancy")
	}

	// Update the file data
	f.RawData = rawData

	return nil
}

func getCheckSum(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}
	checkSumBytes := hasher.Sum(nil)
	return fmt.Sprintf("%x", checkSumBytes), nil
}
