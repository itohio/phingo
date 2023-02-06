package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func validateChain(v ...func(string) error) func(string) error {
	return func(s string) error {
		for _, val := range v {
			if err := val(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func validateEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("must not be empty")
	}
	return nil
}

func validateUniqueKey[T any](kv map[string]T) func(string) error {
	return func(s string) error {
		if _, ok := kv[s]; ok {
			return errors.New("must be unique")
		}
		return nil
	}
}

func validateOneOfKey[T any](kv map[string]T) func(string) error {
	return func(s string) error {
		if _, ok := kv[s]; !ok {
			return errors.New("unknown")
		}
		return nil
	}
}

func validateUnique(arr ...string) func(string) error {
	kv := make(map[string]struct{}, len(arr))
	for _, v := range arr {
		kv[v] = struct{}{}
	}
	return validateUniqueKey(kv)
}

func validateOneOf(arr ...string) func(string) error {
	kv := make(map[string]struct{}, len(arr))
	for _, v := range arr {
		kv[v] = struct{}{}
	}
	return validateOneOfKey(kv)
}

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
