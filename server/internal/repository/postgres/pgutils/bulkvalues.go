package pgutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// BulkUpdateValues generates a SQL VALUES clause for bulk updates.
// It takes a slice of structs and a string specifying the keys for extracting values.
// Returns the constructed VALUES clause, a slice of arguments, and an error if any.
// @param data - slice of structs
// @param keys - string specifying the keys for extracting values from db tag (eg. "id::int, name")
// @returns VALUES syntax string and a slice of arguments for the placeholders
func BulkUpdateValues[T any](data []T, keys string) (string, []interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return "", nil, errors.New("expected a non-empty slice")
	}

	var placeholders []string
	args := make([]interface{}, 0)

	elemType := v.Index(0).Type()
	keyPairs := parseKeys(keys)

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() != reflect.Struct {
			return "", nil, errors.New("expected a struct")
		}

		vals, err := getValues(elem, keyPairs, elemType, &args)
		if err != nil {
			return "", nil, err
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(vals, ", ")))
	}

	valuesClause := fmt.Sprintf("VALUES %s", strings.Join(placeholders, ", "))
	return valuesClause, args, nil
}

func parseKeys(keys string) []string {
	return strings.Split(strings.TrimSpace(keys), ",")
}

func getValues(elem reflect.Value, keyPairs []string, elemType reflect.Type, args *[]interface{}) ([]string, error) {
	var vals []string

	for _, key := range keyPairs {
		parts := strings.Split(strings.TrimSpace(key), "::")
		if len(parts) != 2 {
			return nil, errors.New("keys must be in format 'field::type'")
		}

		fieldName := parts[0]
		fieldType := parts[1]

		dbTag, err := getDBTag(elemType, fieldName)
		if err != nil {
			return nil, err
		}

		fieldValue := elem.FieldByName(dbTag)
		if !fieldValue.IsValid() {
			return nil, fmt.Errorf("invalid field: %s", dbTag)
		}

		placeholder := fmt.Sprintf("$%d::%s", len(*args)+1, fieldType)
		vals = append(vals, placeholder)
		*args = append(*args, fieldValue.Interface())
	}

	return vals, nil
}

func getDBTag(elemType reflect.Type, fieldName string) (string, error) {
	for i := 0; i < elemType.NumField(); i++ {
		f := elemType.Field(i)
		if f.Tag.Get("db") == fieldName {
			return f.Name, nil
		}
	}
	return "", fmt.Errorf("field %s not found in struct", fieldName)
}
