/*
Copyright 2019 Cortex Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package userconfig

import (
	"strings"

	"github.com/cortexlabs/cortex/pkg/lib/cast"
	"github.com/cortexlabs/cortex/pkg/lib/configreader"
	cr "github.com/cortexlabs/cortex/pkg/lib/configreader"
	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/lib/pointer"
	s "github.com/cortexlabs/cortex/pkg/lib/strings"
)

// Returns ValueType, length-one array of <recursive>, or map of {scalar|ValueType -> <recursive>}
// Output types cannot have *_COLUMN types, cannot be compound types, and don't have cortex options (e.g. _default)
// This is used for constants and aggregates
func ValidateOutputTypeSchema(outputTypeSchemaInter interface{}) (interface{}, error) {
	// String
	if typeSchemaStr, ok := outputTypeSchemaInter.(string); ok {
		valueType := ValueTypeFromString(typeSchemaStr)
		if valueType == UnknownValueType {
			if colType := ColumnTypeFromString(typeSchemaStr); colType != UnknownColumnType {
				return nil, ErrorColumnTypeInOutputType(typeSchemaStr)
			}
			if _, err := CompoundTypeFromString(typeSchemaStr); err == nil {
				return nil, ErrorCompoundTypeInOutputType(typeSchemaStr)
			}
			return nil, ErrorInvalidOutputType(typeSchemaStr)
		}
		return valueType, nil
	}

	// List
	if typeSchemaSlice, ok := cast.InterfaceToInterfaceSlice(outputTypeSchemaInter); ok {
		if len(typeSchemaSlice) != 1 {
			return nil, ErrorTypeListLength(typeSchemaSlice)
		}
		elementInputSchema, err := ValidateOutputTypeSchema(typeSchemaSlice[0])
		if err != nil {
			return nil, errors.Wrap(err, s.Index(0))
		}
		return []interface{}{elementInputSchema}, nil
	}

	// Map
	if typeSchemaMap, ok := cast.InterfaceToInterfaceInterfaceMap(outputTypeSchemaInter); ok {
		if len(typeSchemaMap) == 0 {
			return nil, ErrorTypeMapZeroLength(typeSchemaMap)
		}

		var typeKey ValueType
		var typeValue interface{}
		for k, v := range typeSchemaMap {
			if kStr, ok := k.(string); ok {
				typeKey = ValueTypeFromString(kStr)
				if typeKey != UnknownValueType {
					typeValue = v
					break
				}
				if colType := ColumnTypeFromString(kStr); colType != UnknownColumnType {
					return nil, ErrorColumnTypeInOutputType(kStr)
				}
				if _, err := CompoundTypeFromString(kStr); err == nil {
					return nil, ErrorCompoundTypeInOutputType(kStr)
				}
			}
		}

		// Generic map
		if typeValue != nil {
			if len(typeSchemaMap) != 1 {
				return nil, ErrorGenericTypeMapLength(typeSchemaMap)
			}
			valueOutputTypeSchema, err := ValidateOutputTypeSchema(typeValue)
			if err != nil {
				return nil, errors.Wrap(err, string(typeKey))
			}
			return map[interface{}]interface{}{typeKey: valueOutputTypeSchema}, nil
		}

		// Fixed map
		castedTypeSchemaMap := map[interface{}]interface{}{}
		for key, value := range typeSchemaMap {
			if !cast.IsScalarType(key) {
				return nil, configreader.ErrorInvalidPrimitiveType(key, configreader.PrimTypeScalars...)
			}
			if keyStr, ok := key.(string); ok {
				if strings.HasPrefix(keyStr, "_") {
					return nil, ErrorUserKeysCannotStartWithUnderscore(keyStr)
				}
			}

			valueOutputTypeSchema, err := ValidateOutputTypeSchema(value)
			if err != nil {
				return nil, errors.Wrap(err, s.UserStrStripped(key))
			}
			castedTypeSchemaMap[key] = valueOutputTypeSchema
		}
		return castedTypeSchemaMap, nil
	}

	return nil, ErrorInvalidOutputType(outputTypeSchemaInter)
}

// typeSchema is a validated output type schema
func CastConstant(value interface{}, outputTypeSchema interface{}) (interface{}, error) {
	// Check for missing
	if value == nil {
		return nil, ErrorMustBeDefined(outputTypeSchema)
	}

	// ValueType
	if valueType, ok := outputTypeSchema.(ValueType); ok {
		return valueType.CastValue(value)
	}

	// Array
	if typeSchemas, ok := cast.InterfaceToInterfaceSlice(outputTypeSchema); ok {
		values, ok := cast.InterfaceToInterfaceSlice(value)
		if !ok {
			return nil, ErrorUnsupportedLiteralType(value, outputTypeSchema)
		}
		valuesCasted := make([]interface{}, len(values))
		for i, valueItem := range values {
			valueItemCasted, err := CastConstant(valueItem, typeSchemas[0])
			if err != nil {
				return nil, errors.Wrap(err, s.Index(i))
			}
			valuesCasted[i] = valueItemCasted
		}
		return valuesCasted, nil
	}

	// Map
	if typeSchemaMap, ok := cast.InterfaceToInterfaceInterfaceMap(outputTypeSchema); ok {
		valueMap, ok := cast.InterfaceToInterfaceInterfaceMap(value)
		if !ok {
			return nil, ErrorUnsupportedLiteralType(value, outputTypeSchema)
		}

		isGeneric := false
		var genericKey ValueType
		var genericValue interface{}
		for k, v := range typeSchemaMap {
			ok := false
			if genericKey, ok = k.(ValueType); ok {
				isGeneric = true
				genericValue = v
			}
		}

		valueMapCasted := make(map[interface{}]interface{}, len(valueMap))

		// Generic map
		if isGeneric {
			for valueKey, valueVal := range valueMap {
				valueKeyCasted, err := CastConstant(valueKey, genericKey)
				if err != nil {
					return nil, err
				}
				valueValCasted, err := CastConstant(valueVal, genericValue)
				if err != nil {
					return nil, errors.Wrap(err, s.UserStrStripped(valueKey))
				}
				valueMapCasted[valueKeyCasted] = valueValCasted
			}
			return valueMapCasted, nil
		}

		// Fixed map
		for typeSchemaKey, typeSchemaValue := range typeSchemaMap {
			valueValCasted, err := CastConstant(valueMap[typeSchemaKey], typeSchemaValue)
			if err != nil {
				return nil, errors.Wrap(err, s.UserStrStripped(typeSchemaKey))
			}
			valueMapCasted[typeSchemaKey] = valueValCasted
		}
		for valueKey := range valueMap {
			if _, ok := typeSchemaMap[valueKey]; !ok {
				return nil, ErrorUnsupportedLiteralMapKey(valueKey, typeSchemaMap)
			}
		}
		return valueMapCasted, nil
	}

	return nil, ErrorInvalidOutputType(outputTypeSchema) // unexpected
}

type InputSchema struct {
	Type     interface{} `json:"_type" yaml:"_type"` // CompundType, length-one array of *InputSchema, or map of {scalar|CompoundType -> *InputSchema}
	Optional bool        `json:"_optional" yaml:"_optional"`
	Default  interface{} `json:"_default" yaml:"_default"`
	MinCount *int64      `json:"_min_count" yaml:"_min_count"`
	MaxCount *int64      `json:"_max_count" yaml:"_max_count"`
}

// Returns CompundType, length-one array of *InputSchema, or map of {scalar|CompoundType -> *InputSchema}
func validateInputTypeSchema(inputTypeSchemaInter interface{}) (interface{}, error) {
	// String
	if typeSchemaStr, ok := inputTypeSchemaInter.(string); ok {
		compoundType, err := CompoundTypeFromString(typeSchemaStr)
		if err != nil {
			return nil, err
		}
		return compoundType, nil
	}

	// List
	if typeSchemaSlice, ok := cast.InterfaceToInterfaceSlice(inputTypeSchemaInter); ok {
		if len(typeSchemaSlice) != 1 {
			return nil, ErrorTypeListLength(typeSchemaSlice)
		}
		elementInputSchema, err := ValidateInputSchema(typeSchemaSlice[0])
		if err != nil {
			return nil, errors.Wrap(err, s.Index(0))
		}
		return []interface{}{elementInputSchema}, nil
	}

	// Map
	if typeSchemaMap, ok := cast.InterfaceToInterfaceInterfaceMap(inputTypeSchemaInter); ok {
		if len(typeSchemaMap) == 0 {
			return nil, ErrorTypeMapZeroLength(typeSchemaMap)
		}

		var typeKey CompoundType
		var typeValue interface{}
		for k, v := range typeSchemaMap {
			var err error
			typeKey, err = CompoundTypeFromString(k)
			if err == nil {
				typeValue = v
				break
			}
		}

		// Generic map
		if typeValue != nil {
			if len(typeSchemaMap) != 1 {
				return nil, ErrorGenericTypeMapLength(typeSchemaMap)
			}
			valueInputSchema, err := ValidateInputSchema(typeValue)
			if err != nil {
				return nil, errors.Wrap(err, string(typeKey))
			}
			return map[interface{}]interface{}{typeKey: valueInputSchema}, nil
		}

		// Fixed map
		castedTypeSchemaMap := map[interface{}]interface{}{}
		for key, value := range typeSchemaMap {
			if !cast.IsScalarType(key) {
				return nil, configreader.ErrorInvalidPrimitiveType(key, configreader.PrimTypeScalars...)
			}
			if keyStr, ok := key.(string); ok {
				if strings.HasPrefix(keyStr, "_") {
					return nil, ErrorUserKeysCannotStartWithUnderscore(keyStr)
				}
			}

			valueInputSchema, err := ValidateInputSchema(value)
			if err != nil {
				return nil, errors.Wrap(err, s.UserStrStripped(key))
			}
			castedTypeSchemaMap[key] = valueInputSchema
		}
		return castedTypeSchemaMap, nil
	}

	return nil, ErrorInvalidInputType(inputTypeSchemaInter)
}

// Returns InputSchema
func ValidateInputSchema(inputSchemaInter interface{}) (*InputSchema, error) {
	// Check for cortex options vs short form
	if inputSchemaMap, ok := cast.InterfaceToStrInterfaceMap(inputSchemaInter); ok {
		foundUnderscore, foundNonUnderscore := false, false
		for key := range inputSchemaMap {
			if strings.HasPrefix(key, "_") {
				foundUnderscore = true
			} else {
				foundNonUnderscore = true
			}
		}

		if foundUnderscore {
			if foundNonUnderscore {
				return nil, ErrorMixedInputArgOptionsAndUserKeys()
			}

			inputSchemaValidation := &cr.StructValidation{
				StructFieldValidations: []*cr.StructFieldValidation{
					{
						StructField: "Type",
						InterfaceValidation: &cr.InterfaceValidation{
							Required:  true,
							Validator: validateInputTypeSchema,
						},
					},
					{
						StructField:    "Optional",
						BoolValidation: &cr.BoolValidation{},
					},
					{
						StructField:         "Default",
						InterfaceValidation: &cr.InterfaceValidation{},
					},
					{
						StructField: "MinCount",
						Int64PtrValidation: &cr.Int64PtrValidation{
							GreaterThanOrEqualTo: pointer.Int64(0),
						},
					},
					{
						StructField: "MaxCount",
						Int64PtrValidation: &cr.Int64PtrValidation{
							GreaterThanOrEqualTo: pointer.Int64(0),
						},
					},
				},
			}
			inputSchema := &InputSchema{}
			errs := cr.Struct(inputSchema, inputSchemaMap, inputSchemaValidation)

			if errors.HasErrors(errs) {
				return nil, errors.FirstError(errs...)
			}

			if err := validateInputSchemaOptions(inputSchema); err != nil {
				return nil, err
			}

			return inputSchema, nil
		}
	}

	typeSchema, err := validateInputTypeSchema(inputSchemaInter)
	if err != nil {
		return nil, err
	}
	inputSchema := &InputSchema{
		Type: typeSchema,
	}

	if err := validateInputSchemaOptions(inputSchema); err != nil {
		return nil, err
	}

	return inputSchema, nil
}

func validateInputSchemaOptions(inputSchema *InputSchema) error {
	if inputSchema.Default != nil {
		inputSchema.Optional = true
	}

	_, isSlice := cast.InterfaceToInterfaceSlice(inputSchema.Type)
	isGenericMap := false
	if interfaceMap, ok := cast.InterfaceToInterfaceInterfaceMap(inputSchema.Type); ok {
		for k := range interfaceMap {
			_, isGenericMap = k.(CompoundType)
			break
		}
	}

	if inputSchema.MinCount != nil {
		if !isGenericMap && !isSlice {
			return ErrorOptionOnNonIterable(MinCountOptKey)
		}
	}

	if inputSchema.MaxCount != nil {
		if !isGenericMap && !isSlice {
			return ErrorOptionOnNonIterable(MaxCountOptKey)
		}
	}

	if inputSchema.MinCount != nil && inputSchema.MaxCount != nil && *inputSchema.MinCount > *inputSchema.MaxCount {
		return ErrorMinCountGreaterThanMaxCount()
	}

	// Validate default against schema
	if inputSchema.Default != nil {
		var err error
		inputSchema.Default, err = CastInputDefault(inputSchema.Default, inputSchema)
		if err != nil {
			return errors.Wrap(err, DefaultOptKey)
		}
	}

	return nil
}

func CastInputDefault(value interface{}, inputSchema *InputSchema) (interface{}, error) {
	// Check for missing
	if value == nil {
		if inputSchema.Optional {
			return inputSchema.Default, nil
		}
		return nil, ErrorMustBeDefined(inputSchema)
	}

	typeSchema := inputSchema.Type

	// CompoundType
	if compoundType, ok := typeSchema.(CompoundType); ok {
		return compoundType.CastValue(value)
	}

	// array of *InputSchema
	if inputSchemas, ok := cast.InterfaceToInterfaceSlice(typeSchema); ok {
		values, ok := cast.InterfaceToInterfaceSlice(value)
		if !ok {
			return nil, ErrorUnsupportedLiteralType(value, typeSchema)
		}

		if inputSchema.MinCount != nil && int64(len(values)) < *inputSchema.MinCount {
			return nil, ErrorTooFewElements(configreader.PrimTypeList, *inputSchema.MinCount)
		}
		if inputSchema.MaxCount != nil && int64(len(values)) > *inputSchema.MaxCount {
			return nil, ErrorTooManyElements(configreader.PrimTypeList, *inputSchema.MaxCount)
		}

		valuesCasted := make([]interface{}, len(values))
		for i, valueItem := range values {
			valueItemCasted, err := CastInputDefault(valueItem, inputSchemas[0].(*InputSchema))
			if err != nil {
				return nil, errors.Wrap(err, s.Index(i))
			}
			valuesCasted[i] = valueItemCasted
		}
		return valuesCasted, nil
	}

	// Map
	if typeSchemaMap, ok := cast.InterfaceToInterfaceInterfaceMap(typeSchema); ok {
		valueMap, ok := cast.InterfaceToInterfaceInterfaceMap(value)
		if !ok {
			return nil, ErrorUnsupportedLiteralType(value, typeSchema)
		}

		var genericKey CompoundType
		var genericValue *InputSchema
		for k, v := range typeSchemaMap {
			ok := false
			if genericKey, ok = k.(CompoundType); ok {
				genericValue = v.(*InputSchema)
			}
		}

		valueMapCasted := make(map[interface{}]interface{}, len(valueMap))

		// Generic map
		if genericValue != nil {
			if inputSchema.MinCount != nil && int64(len(valueMap)) < *inputSchema.MinCount {
				return nil, ErrorTooFewElements(configreader.PrimTypeMap, *inputSchema.MinCount)
			}
			if inputSchema.MaxCount != nil && int64(len(valueMap)) > *inputSchema.MaxCount {
				return nil, ErrorTooManyElements(configreader.PrimTypeMap, *inputSchema.MaxCount)
			}

			for valueKey, valueVal := range valueMap {
				valueKeyCasted, err := CastInputDefault(valueKey, &InputSchema{Type: genericKey})
				if err != nil {
					return nil, err
				}
				valueValCasted, err := CastInputDefault(valueVal, genericValue)
				if err != nil {
					return nil, errors.Wrap(err, s.UserStrStripped(valueKey))
				}
				valueMapCasted[valueKeyCasted] = valueValCasted
			}
			return valueMapCasted, nil
		}

		// Fixed map
		for typeSchemaKey, typeSchemaValue := range typeSchemaMap {
			valueValCasted, err := CastInputDefault(valueMap[typeSchemaKey], typeSchemaValue.(*InputSchema))
			if err != nil {
				return nil, errors.Wrap(err, s.UserStrStripped(typeSchemaKey))
			}
			valueMapCasted[typeSchemaKey] = valueValCasted
		}
		for valueKey := range valueMap {
			if _, ok := typeSchemaMap[valueKey]; !ok {
				return nil, ErrorUnsupportedLiteralMapKey(valueKey, typeSchemaMap)
			}
		}
		return valueMapCasted, nil
	}

	return nil, ErrorInvalidInputType(typeSchema) // unexpected
}
