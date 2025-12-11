package main

import (
	"errors"
	"fmt"
	"io"
	_ "math"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWithSrcFile           = errors.New("source file problem")
	ErrWithDestFile          = errors.New("destination file problem")
	ErrSrcEqualsDst          = errors.New("destination file equals source file")
)

const readChunkSize = 256

type FileCopier struct {
	fromPath string
	toPath   string
	offset   int64
	limit    int64
}

func NewFileCopier(fromPath, toPath string, offset, limit int64) *FileCopier {
	return &FileCopier{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
	}
}

func (fc *FileCopier) Copy() error {
	if err := fc.validateFilesNotEqual(); err != nil {
		return err
	}

	srcFile, err := fc.openSrcFile()
	if err != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", err)
	}
	defer srcFile.Close()

	if err := fc.validateSrcFile(srcFile); err != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", err)
	}

	dstFile, err := fc.createDstFile()
	if err != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", err)
	}
	defer dstFile.Close()

	data, err := fc.readSrcFile(srcFile)
	if err != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", err)
	}

	if _, err = dstFile.Write(*data); err != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", err)
	}

	return nil
}

func (fc *FileCopier) validateFilesNotEqual() error {
	absFrom, _ := filepath.Abs(fc.fromPath)
	absTo, _ := filepath.Abs(fc.toPath)

	if absFrom == absTo {
		return ErrSrcEqualsDst
	}
	return nil
}

func (fc *FileCopier) openSrcFile() (*os.File, error) {
	return os.Open(fc.fromPath)
}

func (fc *FileCopier) createDstFile() (*os.File, error) {
	return os.Create(fc.toPath)
}

func (fc *FileCopier) validateSrcFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	size := info.Size()
	if size == 0 { // special files like /dev/urandom return 0
		return ErrUnsupportedFile
	}

	if size < fc.offset {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func (fc *FileCopier) readSrcFile(file *os.File) (*[]byte, error) {
	info, _ := file.Stat()
	fileSize := info.Size()

	var bytesToRead int64
	if fc.limit <= 0 || fc.limit > fileSize || fc.offset+fc.limit > fileSize {
		bytesToRead = fileSize - fc.offset
	} else {
		bytesToRead = fc.limit
	}

	chunk := int64(readChunkSize)
	if chunk > bytesToRead {
		chunk = bytesToRead
	}

	data := make([]byte, 0)
	offsetNow := fc.offset

	for (offsetNow - fc.offset) < bytesToRead {
		remaining := bytesToRead - (offsetNow - fc.offset)
		if remaining < chunk {
			chunk = remaining
		}

		tmp := make([]byte, chunk)

		n, err := file.ReadAt(tmp, offsetNow)
		data = append(data, tmp[:n]...)

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		offsetNow += int64(n)
	}

	return &data, nil
}
