package utils

import "os"

func CreateTemp() (*os.File, error) {
	file, err := os.CreateTemp("", "mdconv_file")
	if err != nil {
		return nil, err
	}
	return file, nil
}
