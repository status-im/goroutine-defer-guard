package utils

import "strconv"

func URI(path string, line int) string {
	return path + ":" + strconv.Itoa(line)
}
