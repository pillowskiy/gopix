package domain

import (
	"errors"

	"github.com/pillowskiy/gopix/pkg/snowflake"
)

type ID = snowflake.SnowflakeID

func ParseID(id string) (res ID, err error) {
	if !snowflake.IsStrSnowflakeID(id) {
		return 0, errors.New("string is not snowflake id")
	}

	snowflakeID, err := snowflake.Parse(id)
	if err != nil {
		return
	}

	res = ID(snowflakeID)
	return
}
