package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

func (s *Status) IsValid() bool {
	if *s != StatusOpen && *s != StatusMerged {
		return false
	}

	return true
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*s))
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	status := Status(str)
	if !status.IsValid() {
		return fmt.Errorf("invalid status value: %s", status)
	}

	*s = status
	return nil
}

func (s Status) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid status value: %s", s)
	}

	return string(s), nil
}

func (s *Status) Scan(value any) error {
	if value == nil {
		return fmt.Errorf("cannot scan null into Status")
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Status: expected []byte, got %T", value)
	}

	status := Status(bytes)
	if !status.IsValid() {
		return fmt.Errorf("invalid status value: %s", status)
	}

	*s = status
	return nil
}
