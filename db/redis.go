package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func ScanMap(values []interface{}) (map[string]string, error) {
	results := make(map[string]string)
	var err error

	for len(values) > 0 {
		var key string
		var value string

		values, err = redis.Scan(values, &key)

		if err != nil {
			return nil, err
		}

		if len(values) > 0 {
			values, err = redis.Scan(values, &value)

			if err != nil {
				return nil, err
			}

			results[key] = value
		} else {
			fmt.Println("Unable to find value for %s.", key)
			results[key] = ""
		}

	}

	return results, nil
}
