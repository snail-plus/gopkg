// Copyright 2024 eve.  All rights reserved.

package reflect

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func ToGormDBMap(obj any, fields []string) (map[string]any, error) {
	reflectType := reflect.ValueOf(obj).Type()
	reflectValue := reflect.ValueOf(obj)
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflect.ValueOf(obj).Elem()
	}

	ret := make(map[string]any, 0)
	for _, f := range fields {
		fs, exist := reflectType.FieldByName(f)
		if !exist {
			return nil, fmt.Errorf("unknow field " + f)
		}

		tagMap := parseTagSetting(fs.Tag)
		gormfiled, exist := tagMap["COLUMN"]
		if !exist {
			return nil, fmt.Errorf("undef gorm field " + f)
		}

		ret[gormfiled] = reflectValue.FieldByName(f)
	}
	return ret, nil
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func GetObjFieldsMap(obj any, fields []string) map[string]any {
	ret := make(map[string]any)

	modelReflect := reflect.ValueOf(obj)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()
	var fieldData any
	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)
		if len(fields) != 0 && !findString(fields, modelRefType.Field(i).Name) {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			fieldData = GetObjFieldsMap(field.Interface(), []string{})
		default:
			fieldData = field.Interface()
		}

		ret[modelRefType.Field(i).Name] = fieldData
	}

	return ret
}

func CopyObj(from any, to any, fields []string) (changed bool, err error) {
	fromMap := GetObjFieldsMap(from, fields)
	toMap := GetObjFieldsMap(to, fields)
	if reflect.DeepEqual(fromMap, toMap) {
		return false, nil
	}

	t := reflect.ValueOf(to).Elem()
	for k, v := range fromMap {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
	return true, nil
}

// CopyObjViaYaml marshal "from" to yaml data, then unMarshal data to "to".
func CopyObjViaYaml(to any, from any) error {
	if from == nil || to == nil {
		return nil
	}

	data, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, to)
}

// StructName used to get the struct name from the obj.
func StructName(obj any) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

// findString return true if target in slice, return false if not.
func findString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

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
