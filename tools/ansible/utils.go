package ansible

import (
	"strings"
)

func removeOtherString(input string) string {
	noFileName := strings.Split(string(input), " ")[0]
	noSpaces := strings.ReplaceAll(noFileName, " ", "")
	noNewlines := strings.ReplaceAll(noSpaces, "\n", "")
	return noNewlines
}
