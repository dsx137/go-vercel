package vercelkit

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func ReadParamsFromQuery[T any](queryParams url.Values) (*T, error) {
	params := new(T)
	missing := make([]string, 0)
	val := reflect.ValueOf(params).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("query")
		var paramNames []string
		if tag != "" {
			paramNames = strings.Split(tag, ",")
		} else {
			paramNames = []string{CamelToSnake(field.Name)}
		}

		var paramValue string
		for _, name := range paramNames {
			if v := queryParams.Get(name); v != "" {
				paramValue = v
				break
			}
		}

		if paramValue == "" {
			missing = append(missing, paramNames[0])
			continue
		}

		f := val.Field(i)
		switch f.Kind() {
		case reflect.String:
			f.SetString(paramValue)
		case reflect.Int, reflect.Int64, reflect.Int32:
			if v, err := strconv.ParseInt(paramValue, 10, 64); err == nil {
				f.SetInt(v)
			}
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			if v, err := strconv.ParseUint(paramValue, 10, 64); err == nil {
				f.SetUint(v)
			}
		case reflect.Float64, reflect.Float32:
			if v, err := strconv.ParseFloat(paramValue, 64); err == nil {
				f.SetFloat(v)
			}
		case reflect.Slice:
			elemKind := f.Type().Elem().Kind()
			items := strings.Split(paramValue, ",")
			slice := reflect.MakeSlice(f.Type(), len(items), len(items))
			for j, item := range items {
				item = strings.TrimSpace(item)
				switch elemKind {
				case reflect.String:
					slice.Index(j).SetString(item)
				case reflect.Int, reflect.Int64, reflect.Int32:
					if v, err := strconv.ParseInt(item, 10, 64); err == nil {
						slice.Index(j).SetInt(v)
					}
				case reflect.Uint, reflect.Uint64, reflect.Uint32:
					if v, err := strconv.ParseUint(item, 10, 64); err == nil {
						slice.Index(j).SetUint(v)
					}
				case reflect.Float64, reflect.Float32:
					if v, err := strconv.ParseFloat(item, 64); err == nil {
						slice.Index(j).SetFloat(v)
					}
				}
			}
			f.Set(slice)
		}
	}

	if len(missing) != 0 {
		return nil, fmt.Errorf("missing parameters: %v", missing)
	}

	return params, nil
}
