package chibisnes

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/edsrzf/mmap-go"
)

type SRAM struct {
	file *os.File
	mmap mmap.MMap

	filePath string
	size     int
}

func NewSRAM(filePath string, size int) *SRAM {
	_, err := os.Stat(filePath)
	var sram *SRAM
	if err == nil {
		sram = readRAMFile(filePath, size)
		log.Printf("SRAM: file loaded. Path: %s\n", filePath)
	} else if os.IsNotExist(err) {
		err := createRAMFile(filePath, size)
		if err != nil {
			log.Printf("SRAM: file create failed. Path: %s Error: %s\n", filePath, err)
			return nil
		}
		sram = readRAMFile(filePath, size)
		log.Printf("SRAM: file created & loaded. Path: %s\n", filePath)
	}

	return sram
}

func (s *SRAM) Read(addr uint32) byte {
	return s.mmap[addr]
}

func (s *SRAM) Write(addr uint32, value byte) {
	s.mmap[addr] = value
}

func (s *SRAM) Close() {
	s.mmap.Unmap()
	s.file.Close()
}

func createRAMFile(ramFilePath string, size int) error {
	file, err := os.Create(ramFilePath)
	if err != nil {
		return err
	}

	_, err = file.Seek(int64(size-1), 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{0})
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func readRAMFile(ramFilePath string, size int) *SRAM {
	ramFile, err := os.OpenFile(ramFilePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	ramMMap, err := mmap.Map(ramFile, mmap.RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	return &SRAM{
		file:     ramFile,
		mmap:     ramMMap,
		filePath: ramFilePath,
		size:     size,
	}
}

func getFileNameWithoutExtension(filePath string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
