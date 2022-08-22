package common

import (
	"os"
)

// TODO: use generics
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

// TODO: use generics
func RemoveDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

const bytesInMB = 1024 * 1024

func FileSizeMB(bytes int64) float64 {
	kilobytes := (bytes / tenTwentyFour)
	return (float64)(kilobytes / tenTwentyFour)
}
