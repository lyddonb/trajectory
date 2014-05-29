package db

import "github.com/garyburd/redigo/redis"

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

		values, err = redis.Scan(values, &value)

		if err != nil {
			return nil, err
		}

		if len(values) > 0 {
			results[key] = value
		} else {
			results[key] = ""
		}
	}

	return results, nil
}
