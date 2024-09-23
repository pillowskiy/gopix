package snowflake

import (
	"regexp"
)

const minEpochLen = 13

var numericRegex = regexp.MustCompile(`^\d{1,19}$`)

func IsStrSnowflakeID(str string) bool {
	if len(str) < minEpochLen {
		return false
	}

	return numericRegex.MatchString(str)
}
