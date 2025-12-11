package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWithSrcFile           = errors.New("source file problem")
	ErrWithDestFile          = errors.New("destination file problem")
	ErrSrcEqualsDst          = errors.New("destination file equals source file")
)

const (
	readChunkSize = 256
)

type FileCopier struct {
	fromPath, toPath string
	offset, limit    int64

	progress *pb.ProgressBar
}

func NewFileCopier(fromPath, toPath string, offset, limit int64, pb *pb.ProgressBar) *FileCopier {
	return &FileCopier{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
		progress: pb,
	}
}

func (fc *FileCopier) Copy() error {
	filesEqualErr := fc.validateFilesNotEqual()
	if filesEqualErr != nil {
		return filesEqualErr
	}

	srcFile, openingErr := fc.openSrcFile()
	if openingErr != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", openingErr)
	}
	defer srcFile.Close()

	validationErr := fc.validateSrcFile(srcFile)
	if validationErr != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", validationErr)
	}

	dstFile, fileCreateErr := fc.createDstFile()
	if fileCreateErr != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", fileCreateErr)
	}
	defer dstFile.Close()

	buffer, readErr := fc.readSrcFile(srcFile)
	if readErr != nil {
		return fmt.Errorf(ErrWithSrcFile.Error()+": %w", readErr)
	}

	_, writeErr := dstFile.Write(*buffer)
	if writeErr != nil {
		return fmt.Errorf(ErrWithDestFile.Error()+": %w", writeErr)
	}

	return nil
}

func (fc *FileCopier) validateFilesNotEqual() error {
	absFromPath, _ := filepath.Abs(fc.fromPath)
	absToPath, _ := filepath.Abs(fc.toPath)

	if absFromPath == absToPath {
		return ErrSrcEqualsDst
	}

	return nil
}

func (fc *FileCopier) openSrcFile() (*os.File, error) {
	file, err := os.Open(fc.fromPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fc *FileCopier) createDstFile() (*os.File, error) {
	toFile, err := os.Create(fc.toPath)
	if err != nil {
		return nil, err
	}

	return toFile, nil
}

func (fc *FileCopier) validateSrcFile(file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 { // for files with unknown size (e.g. /dev/urandom)
		return ErrUnsupportedFile
	} else if fileSize < fc.offset {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func (fc *FileCopier) readSrcFile(file *os.File) (*[]byte, error) {
	buffer := make([]byte, 0)
	fileInfo, _ := file.Stat() // error suppressed because file already validated
	fileSize := fileInfo.Size()

	var bytesToRead int64
	if fc.limit <= 0 || fc.limit > fileSize || fc.offset+fc.limit > fileSize {
		bytesToRead = fileSize - fc.offset
	} else {
		bytesToRead = fc.limit
	}

	chunkSize := int64(readChunkSize)
	if chunkSize > bytesToRead {
		chunkSize = bytesToRead
	}

	fc.startProgressBar(int64(math.Ceil(float64(bytesToRead) / float64(chunkSize))))

	stepOffset := fc.offset
	for (stepOffset - fc.offset) < bytesToRead {
		remainingBytes := bytesToRead - (stepOffset - fc.offset)
		if remainingBytes < chunkSize {
			chunkSize = remainingBytes
		}

		// use tmpBuf to make correct step amount and process by exactly readChunkSize bytes
		tmpBuf := make([]byte, 0, chunkSize)

		readBytes, readErr := file.ReadAt(tmpBuf[len(tmpBuf):cap(tmpBuf)], stepOffset)
		tmpBuf = tmpBuf[:len(tmpBuf)+readBytes]
		buffer = append(buffer, tmpBuf...)

		if readErr != nil {
			if readErr == io.EOF {
				break
			}

			return nil, fmt.Errorf(ErrWithSrcFile.Error()+": %w", readErr)
		}

		stepOffset += int64(readBytes)
		fc.incProgressBar()
	}

	fc.finishProgressBar()

	return &buffer, nil
}

func (fc *FileCopier) startProgressBar(steps int64) {
	if fc.progress == nil {
		return
	}

	fc.progress.SetTotal(steps).Start()
}

func (fc *FileCopier) incProgressBar() {
	if fc.progress == nil {
		return
	}

	fc.progress.Increment()
}

func (fc *FileCopier) finishProgressBar() {
	if fc.progress == nil {
		return
	}

	fc.progress.Finish()
}
