package cassa

import "github.com/gocql/gocql"

func ArrayOfInt(ids []int) []interface{} {
	slice := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		slice = append(slice, id)
	}
	return slice
}

func ArrayOfUUID(ids []string) ([]interface{}, error) {
	slice := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		val, err := gocql.ParseUUID(id)
		if err != nil {
			return nil, err
		}
		slice = append(slice, val)
	}
	return slice, nil
}
