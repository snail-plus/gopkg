package reflects

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

func GetNumberField(obj any, key string) float64 {
	switch obj.(type) {
	case map[string]interface{}:
		return obj.(map[string]interface{})[key].(float64)

	case struct{}:
		val, err := GetFieldValueByName(obj, key)
		if err != nil {
			return 0
		}

		return cast.ToFloat64(val)
	default:
		return 0
	}
}

func GetFieldValueByName(obj interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(obj)

	// 确保传入的是指针，然后获取其指向的实际值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保传入的是结构体
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %T", val)
	}

	// 获取结构体的字段值
	fieldVal := val.FieldByName(fieldName)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("field %s not found in struct", fieldName)
	}

	// 返回字段的值
	return fieldVal.Interface(), nil
}

func ExistFiled(obj interface{}, fieldName string) bool {
	_, err := GetFieldValueByName(obj, fieldName)
	return err == nil
}

func SetFieldValueByName(myInstance interface{}, fieldName string, fieldValue interface{}) error {
	// 获取myInstance的反射对象
	var val reflect.Value

	if reflect.TypeOf(myInstance).Kind() == reflect.Ptr {
		val = reflect.ValueOf(myInstance).Elem()
	} else {
		val = reflect.ValueOf(&myInstance).Elem()
	}

	// 遍历结构体的字段，找到匹配的字段名称
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)

		// 检查字段名称是否匹配
		if typeField.Name == fieldName {
			if !field.CanSet() {
				return errors.New("Cannot set field:" + fieldName)
			}

			// 根据字段类型设置值
			switch field.Kind() {
			case reflect.String:
				field.SetString(fmt.Sprintf("%v", fieldValue))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInt(reflect.ValueOf(fieldValue).Convert(field.Type()).Int())
			case reflect.Bool:
				field.SetBool(fieldValue.(bool))
			case reflect.Float64:
				field.SetFloat(reflect.ValueOf(fieldValue).Convert(field.Type()).Float())
			case reflect.Ptr:
				// 解引用指针
				if field.IsNil() {
					// 如果指针是nil，需要先分配内存
					field.Set(reflect.New(field.Type().Elem()))
				}
				actualField := field.Elem()

				// 确保解引用后的字段可设置
				if !actualField.CanSet() {
					return fmt.Errorf("field %s not found in struct", fieldName)
				}

				// 根据字段类型设置值
				switch actualField.Kind() {
				case reflect.String:
					actualField.SetString(fmt.Sprintf("%v", fieldValue))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					field.SetInt(reflect.ValueOf(fieldValue).Convert(field.Type()).Int())
				// 更多的类型可以在这里添加
				default:
					fmt.Printf("Unsupported field type %s\n", actualField.Kind())
				}
			// 更多的类型可以在这里添加
			default:
				return errors.New("Unsupported field type " + field.Kind().String())
			}

			break
		}
	}

	return nil
}

func GetStructName(myInstance interface{}) string {
	// 获取myInstance的反射对象Value
	val := reflect.ValueOf(myInstance)

	if val.Kind() == reflect.Ptr {
		// 获取指针指向的元素类型
		elemType := val.Elem()
		// 输出结构体的名称（通过指针）
		return elemType.Type().Name()
	}

	// 获取myInstance的反射类型Type
	typeVal := val.Type()

	// 输出结构体的名称
	return typeVal.Name()
}
