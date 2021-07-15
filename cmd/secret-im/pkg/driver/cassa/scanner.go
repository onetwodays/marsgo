package cassa

import (
	"errors"
	"reflect"

	"github.com/gocql/gocql"
)

// 更好的扫描器
type Scanner struct {
	table    []int
	typeName string
}

// 创建扫描器
func NewScanner(data gocql.RowData, value interface{}) (*Scanner, error) {
	mapper := make(map[string]int)
	for idx, field := range data.Columns {
		mapper[field] = idx
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("non struct object")
	}

	scanner := Scanner{
		typeName: t.String(),
		table:    make([]int, len(data.Columns)),
	}
	for idx := range scanner.table {
		scanner.table[idx] = -1
	}

	for idx := 0; idx < t.NumField(); idx++ {
		tag, ok := t.Field(idx).Tag.Lookup("cql")
		if !ok {
			continue
		}

		pos, ok := mapper[tag]
		if !ok {
			continue
		}
		scanner.table[pos] = idx
	}
	return &scanner, nil
}

// 扫描多行数据
func (scanner *Scanner) ScanRows(iter *gocql.Iter, data gocql.RowData, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		return errors.New("non ptr object")
	}

	results := v.Elem()
	if results.Type().Kind() != reflect.Slice {
		return errors.New("non slice object")
	}
	if results.Type().Elem().String() != scanner.typeName {
		return errors.New("elem type mismatch")
	}
	elemType := results.Type().Elem()
	results.Set(reflect.MakeSlice(results.Type(), 0, iter.NumRows()))

	for iter.Scan(data.Values...) {
		elem := reflect.New(elemType).Elem()
		for pos, val := range data.Values {
			idx := scanner.table[pos]
			if idx == -1 {
				continue
			}

			value := reflect.ValueOf(val)
			if value.IsNil() {
				continue
			}

			if value.Elem().Kind() == reflect.Slice {
				n := value.Elem().Len()
				dst := reflect.MakeSlice(value.Elem().Type(), n, n)
				reflect.Copy(dst, value.Elem())
				elem.Field(idx).Set(dst)
			} else {
				elem.Field(idx).Set(value.Elem().Convert(elem.Field(idx).Type()))
			}
		}
		results.Set(reflect.Append(results, elem))
	}
	return nil
}
