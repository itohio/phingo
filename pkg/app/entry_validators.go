package app

import (
	"fmt"
	"strconv"
	"time"
)

func validateInt(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func validateIntMin(min int64) func(string) error {
	return func(s string) error {
		v, err := strconv.ParseInt(s, 10, 64)
		if v < min {
			return fmt.Errorf("value must not be less than %d", min)
		}
		return err
	}
}

func validateFloat(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	return err
}

func validateDuration(s string) error {
	_, err := time.ParseDuration(s)
	return err
}

func validateDurationMin(min time.Duration) func(string) error {
	return func(s string) error {
		v, err := time.ParseDuration(s)
		if v < min {
			return fmt.Errorf("value must not be less than %s", min.String())
		}
		return err
	}
}
