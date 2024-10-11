package tools

import (
	"io"
	"os"
)

func ReadFileOrStdin(inputFile string) ([]byte, error) {
	if inputFile == "" {
		return io.ReadAll(os.Stdin)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return io.ReadAll(file)
}
