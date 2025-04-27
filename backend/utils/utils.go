package utils

import (
	"io"
	"strings"
)

// Shellescape escapes a string for shell usage
func Shellescape(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// CleanOutput normalizes program output for comparison
func CleanOutput(output string) string {
	output = strings.ReplaceAll(output, "\r\n", "\n")
	output = strings.ReplaceAll(output, "\r", "\n")

	startIndex := 0
	for i, c := range output {
		if c >= 32 && c <= 126 {
			startIndex = i
			break
		}
	}
	if startIndex > 0 {
		output = output[startIndex:]
	}

	return strings.TrimSpace(output)
}

// FormatForDisplay formats a string for display in error messages
func FormatForDisplay(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

// StdCopy handles copying from a Docker container's multiplexed output
func StdCopy(stdout, stderr io.Writer, src io.Reader) (written int64, err error) {
	var buf [8]byte
	var totalN int64

	for {
		_, err = io.ReadFull(src, buf[:])
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return totalN, err
		}

		frameSize := int64(buf[4])<<24 | int64(buf[5])<<16 | int64(buf[6])<<8 | int64(buf[7])

		var dst io.Writer
		if buf[0] == 1 {
			dst = stdout
		} else {
			dst = stderr
		}

		n, err := io.CopyN(dst, src, frameSize)
		totalN += n
		if err != nil {
			return totalN, err
		}
	}
}
