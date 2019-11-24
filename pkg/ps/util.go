package ps

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readSymlink(filePath string) (string, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != os.ModeSymlink {
		return "", fmt.Errorf("%v is not symlink.", filePath)
	}
	realPath, err := os.Readlink(filePath)
	if err != nil {
		return "", err
	}
	return realPath, nil
}

func parseTSInt32(value string) ([]int32, error) {
	pvals := strings.Split(value, "\t")
	s := make([]int32, len(pvals))
	for i, str := range pvals {
		pval, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		s[i] = int32(pval)
	}
	return s, nil
}
