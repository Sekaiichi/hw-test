package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyErrors(t *testing.T) {
	limit := int64(0)
	offset := int64(0)
	from := "testdata/input.txt"
	to := "out.txt"
	fromFileStat, err := os.Stat(from)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("offset more than filesize", func(t *testing.T) {
		fileSize := fromFileStat.Size()

		fileCopier := NewFileCopier(from, to, fileSize+100, limit, nil)
		copierErr := fileCopier.Copy()

		require.ErrorIs(t, copierErr, ErrOffsetExceedsFileSize)
	})

	t.Run("limit more than filesize", func(t *testing.T) {
		fileSize := fromFileStat.Size()

		fileCopier := NewFileCopier(from, to, offset, fileSize+100, nil)
		copierErr := fileCopier.Copy()

		require.NoError(t, copierErr)
	})

	t.Run("unsupported file", func(t *testing.T) {
		fileCopier := NewFileCopier("/dev/urandom", to, offset, limit, nil)
		copierErr := fileCopier.Copy()

		require.ErrorIs(t, copierErr, ErrUnsupportedFile)
	})

	t.Run("negative offset", func(t *testing.T) {
		fileCopier := NewFileCopier(from, to, -1, limit, nil)
		copierErr := fileCopier.Copy()

		require.ErrorContains(t, copierErr, ErrWithSrcFile.Error())
	})

	t.Run("destination file problem", func(t *testing.T) {
		fileCopier := NewFileCopier(from, "/root/out.txt", offset, limit, nil)
		copierErr := fileCopier.Copy()

		require.ErrorContains(t, copierErr, ErrWithDestFile.Error())
	})

	t.Run("destination file equals source file", func(t *testing.T) {
		fileCopier := NewFileCopier(from, from, offset, limit, nil)
		copierErr := fileCopier.Copy()

		require.ErrorContains(t, copierErr, ErrSrcEqualsDst.Error())
	})

	err = os.Remove(to)
	if err != nil {
		t.Fatal(err)
	}
}
