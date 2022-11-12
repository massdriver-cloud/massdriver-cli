package common

import (
	"os"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func WriteFile(filePath string, data []byte, errToBytes error) error {
	if errToBytes != nil {
		return errToBytes
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, errWrite := file.Write(data); errWrite != nil {
		return errWrite
	}
	return nil
}

func RemoveDuplicateValues(stringSlice []string) []string {
	keysSeen := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if !keysSeen[entry] {
			list = append(list, entry)
			keysSeen[entry] = true
		}
	}

	return list
}

const bytesInMB = 1024 * 1024

func FileSizeMB(bytes int64) float64 {
	return (float64)(bytes / bytesInMB)
}
