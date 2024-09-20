package snowflake

import (
	"encoding/json"
	"strconv"
)

type SnowflakeID uint64

func Parse(str string) (SnowflakeID, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return SnowflakeID(id), nil
}

func (s SnowflakeID) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s SnowflakeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SnowflakeID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	id, err := Parse(str)
	if err != nil {
		return err
	}

	*s = id
	return nil
}
