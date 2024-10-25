package serializer

import (
	"context"
	"github.com/dobyte/due/v2/encoding/json"
	"gorm.io/gorm/schema"
	"reflect"
)

func init() {
	schema.RegisterSerializer("json", JSONSerializer{})
}

type JSONSerializer struct {
}

// Scan implements serializer interface
func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			bytes, err = json.Marshal(v)
			if err != nil {
				return err
			}
		}

		if len(bytes) > 0 {
			err = json.Unmarshal(bytes, fieldValue.Interface())
		}
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (JSONSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	result, err := json.Marshal(fieldValue)
	if err != nil {
		return "", err
	}

	if string(result) == "null" {
		switch field.FieldType.Kind() {
		case reflect.Array, reflect.Slice:
			return "[]", nil
		default:
			return "{}", nil
		}
	}

	return string(result), nil
}
