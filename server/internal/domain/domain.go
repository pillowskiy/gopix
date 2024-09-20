package domain

import "github.com/pillowskiy/gopix/pkg/snowflake"

type ID = snowflake.SnowflakeID

func ParseID(id string) (res ID, err error) {
	snowflakeID, err := snowflake.Parse(id)
	if err != nil {
		return
	}

	res = ID(snowflakeID)
	return
}
