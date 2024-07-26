package args

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Parse(args []string, options interface{}) error {
	optsValue := reflect.ValueOf(options)
	if optsValue.Kind() != reflect.Pointer || optsValue.Elem().Kind() != reflect.Struct {
		return errors.New("options parameter must be a pointer to a struct")
	}

	optsElem := optsValue.Elem()
	optsType := optsElem.Type()

	optionsMap := make(map[string]reflect.Value)
	shortKeysMap := make(map[string]string)
	defaultValuesMap := make(map[string]string)
	notSetValuesMap := make(map[string]bool)

	for i := 0; i < optsElem.NumField(); i++ {
		field := optsType.Field(i)

		tag := field.Tag.Get("cli")
		if tag == "" {
			continue
		}

		if !isSupportedType(field.Type) {
			if tag != "" {
				return fmt.Errorf("field %q has an unsupported type %q, but it has a tag %q", field.Name, field.Type.Kind(), "cli")
			}
			continue
		}

		parsedTag := parseTagData(tag)

		optionsMap[parsedTag.name] = optsElem.Field(i)
		notSetValuesMap[parsedTag.name] = true

		var ok bool
		var val string
		var isBool bool

		ok, val, isBool = parsedTag.getOption("short")
		if ok {
			if isBool {
				return fmt.Errorf("%q is bool", "short")
			}
			if val == "" {
				return fmt.Errorf("%q is empty", "short")
			}
			shortKeysMap[val] = parsedTag.name
		}

		ok, val, isBool = parsedTag.getOption("default")
		if ok {
			if isBool {
				return fmt.Errorf("%q is bool", "default")
			}
			if val == "" {
				return fmt.Errorf("%q is empty", "default")
			}
			defaultValuesMap[parsedTag.name] = val
		}
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-") {
			isSet := false
			key := arg
			value := ""
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				key = parts[0]
				value = parts[1]
				isSet = true
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				value = args[i+1]
				isSet = true
				i++
			}

			key = strings.TrimLeft(key, "-")

			mappedKey, ok := shortKeysMap[key]
			if ok {
				key = mappedKey
			}

			field, ok := optionsMap[key]
			if ok {
				if !isSet && field.Type().Kind() != reflect.Bool {
					return fmt.Errorf("value for the option %q is not set", key)
				}

				err := setValue(field, value)
				if err != nil {
					return fmt.Errorf("value setting error: %s", err)
				}
				delete(notSetValuesMap, key)
			} else {
				return fmt.Errorf("unknown option: %q", key)
			}
		} else {
			return fmt.Errorf("unknown option: %q", arg)
		}
	}

	for key, _ := range notSetValuesMap {
		val, ok1 := defaultValuesMap[key]
		field, ok2 := optionsMap[key]
		if ok1 && ok2 {
			err := setValue(field, val)
			if err != nil {
				return fmt.Errorf("value setting error: %s", err)
			}
		}
	}

	return nil
}

func parseTagData(tag string) *tagData {
	name, options, _ := strings.Cut(tag, ",")
	return &tagData{
		name:    name,
		options: options,
	}
}

type tagData struct {
	name    string
	options string
}

// ok, val, isBool
func (t *tagData) getOption(name string) (bool, string, bool) {
	if len(t.options) == 0 {
		return false, "", false
	}

	var option string
	var value string
	s := t.options
	for s != "" {
		option, s, _ = strings.Cut(s, ",")
		if option == name {
			return true, "", true
		}

		option, value, _ = strings.Cut(option, "=")
		if option == name {
			return true, value, false
		}
	}

	return false, "", false
}

func isSupportedType(typ reflect.Type) bool {
	return isScalarKind(typ.Kind())
}

func isScalarKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func setValue(val reflect.Value, value string) error {
	if !val.IsValid() {
		return errors.New("invalid reflect.Value")
	}

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if !val.CanSet() {
		return errors.New("value cannot be set")
	}

	switch val.Kind() {
	case reflect.Bool:
		if value == "" {
			val.SetBool(true)
		} else {
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("failed to convert %q to bool: %w", value, err)
			}
			val.SetBool(boolValue)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			val.SetInt(0)
		} else {
			intValue, err := strconv.ParseInt(value, 10, val.Type().Bits())
			if err != nil {
				return fmt.Errorf("failed to convert %q to int: %w", value, err)
			}
			val.SetInt(intValue)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			val.SetUint(0)
		} else {
			uintValue, err := strconv.ParseUint(value, 10, val.Type().Bits())
			if err != nil {
				return fmt.Errorf("failed to convert %q to uint: %w", value, err)
			}
			val.SetUint(uintValue)
		}

	case reflect.Float32, reflect.Float64:
		if value == "" {
			val.SetFloat(0)
		} else {
			floatValue, err := strconv.ParseFloat(value, val.Type().Bits())
			if err != nil {
				return fmt.Errorf("failed to convert %q to float: %w", value, err)
			}
			val.SetFloat(floatValue)
		}

	case reflect.String:
		val.SetString(value)

	default:
		return fmt.Errorf("unsupported type: %q", val.Type())
	}

	return nil
}
