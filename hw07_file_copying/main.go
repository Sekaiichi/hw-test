package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	fileCopier := NewFileCopier(from, to, offset, limit, pb.New(0))
	err := fileCopier.Copy()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
