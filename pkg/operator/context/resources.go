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

package context

import (
	"strings"

	"github.com/cortexlabs/cortex/pkg/lib/cast"
	"github.com/cortexlabs/cortex/pkg/lib/configreader"
	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/lib/regex"
	s "github.com/cortexlabs/cortex/pkg/lib/strings"
	"github.com/cortexlabs/cortex/pkg/operator/api/context"
	"github.com/cortexlabs/cortex/pkg/operator/api/resource"
	"github.com/cortexlabs/cortex/pkg/operator/api/userconfig"
)

func ValidateInput(
	input interface{},
	schema *userconfig.InputSchema,
	validResourceTypes []resource.Type,
	validResources map[string]context.Resource,
	allResources map[string]userconfig.Resource,
	userAggregators map[string]*context.Aggregator,
	userTransformers map[string]*context.Transformer,
) (interface{}, string, string, error) {
	inputWithIDs, inputWithTagIDs, err := replaceResourceIDsAndValidateResourceTypes(input, validResourceTypes, validResources, allResources)
	if err != nil {
		return "", "", err
	}

	castedInput, err = validateInputRuntimeTypes(input, validResources, userAggregators, userTransformers)
	if err != nil {
		return "", "", err
	}

	return castedInput, hash.Any(inputWithIDs), hash.Any(inputWithTagIDs), nil
}

// don't set default (python has to), check optional and min/max
// resource references have already been validated to exist in validResources
// schema is already a validated InputSchema
func validateInputRuntimeTypes(
	input interface{},
	schema *userconfig.InputSchema,
	validResources map[string]context.Resource,
	userAggregators map[string]*context.Aggregator,
	userTransformers map[string]*context.Transformer,
) (interface{}, error) {

	// Skip validation (schema is nil if user didn't define the aggregator/transformer/estimator)
	if schema == nil {
		return input, nil
	}

	// Check for missing
	if input == nil {
		if schema.Optional {
			return nil, nil
		}
		return nil, ErrorMustBeDefined(schema)
	}

	// Check if input is Cortex resource
	if resourceName, ok := configreader.ExtractResourceName(input); ok {
		res := validResources[resourceName]
		if res == nil {
			return nil, errors.New(resourceName, "missing resource") // unexpected
		}
		switch res.GetResourceType() {
		case resource.ConstantType:
			constant := res.(context.Constant)
			err, inputCasted := validateInputRuntimeTypes(constant.Value, schema, validResources, userAggregators, userTransformers)
			if err != nil {
				return nil, errors.Wrap(err, userconfig.Identify(constant), userconfig.ValueKey)
			}
			return inputCasted, nil
		case resource.RawColumnType:
			rawColumn := res.(context.RawColumn)
			if rawColumn.GetType() != nil {
				if err := validateInputRuntimeOutputTypes(rawColumn.GetType(), schema); err != nil {
					return nil, errors.Wrap(err, userconfig.Identify(rawColumn), userconfig.TypeKey)
				}
			}
			return input, nil
		case resource.AggregateType:
			aggregate := res.(context.Aggregate)
			aggregator, _ := getAggregator(aggregate.Aggregator, userAggregators)
			if aggregator.OutputType != nil {
				if err := validateInputRuntimeOutputTypes(aggregator.OutputType, schema); err != nil {
					return nil, errors.Wrap(err, userconfig.Identify(aggregate), userconfig.Identify(aggregator), userconfig.OutputTypeKey)
				}
			}
			return input, nil
		case resource.TransformedColumnType:
			transformedColumn := res.(context.TransformedColumn)
			transformer, _ := getTransformer(transformedColumn.Transformer, userTransformers)
			if transformer.OutputType != nil {
				if err := validateInputRuntimeOutputTypes(transformer.OutputType, schema); err != nil {
					return nil, errors.Wrap(err, userconfig.Identify(transformedColumn), userconfig.Identify(transformer), userconfig.OutputTypeKey)
				}
			}
			return input, nil
		default:
			return nil, errors.New(res.GetResourceType().String(), "unsupported resource type") // unexpected
		}
	}

	typeSchema := schema.Type

	// CompoundType
	if compoundType, ok := typeSchema.(userconfig.CompoundType); ok {
		return compoundType.CastValue(input)
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
			valueItemCasted, err := validateInputRuntimeTypes(valueItem, inputSchemas[0].(*InputSchema), validResources, userAggregators, userTransformers)
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
				valueKeyCasted, err := validateInputRuntimeTypes(valueKey, &userconfig.InputSchema{Type: genericKey}, validResources, userAggregators, userTransformers)
				if err != nil {
					return nil, err
				}
				valueValCasted, err := validateInputRuntimeTypes(valueVal, genericValue, validResources, userAggregators, userTransformers)
				if err != nil {
					return nil, errors.Wrap(err, s.UserStrStripped(valueKey))
				}
				valueMapCasted[valueKeyCasted] = valueValCasted
			}
			return valueMapCasted, nil
		}

		// Fixed map
		for typeSchemaKey, typeSchemaValue := range typeSchemaMap {
			valueValCasted, err := validateInputRuntimeTypes(valueMap[typeSchemaKey], typeSchemaValue.(*InputSchema), validResources, userAggregators, userTransformers)
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

	return nil, userconfig.ErrorInvalidInputType(typeSchema) // unexpected
}

// outputType should be ValueType|ColumnType, length-one array of <recursive>, or map of {scalar|ValueType -> <recursive>}
func validateInputRuntimeOutputTypes(outputType interface{}, schema *userconfig.InputSchema) error {
	// Check for missing
	if outputType == nil {
		if schema.Optional {
			return nil, nil
		}
		return nil, ErrorMustBeDefined(schema)
	}

	typeSchema := schema.Type

	// CompoundType
	if compoundType, ok := typeSchema.(userconfig.CompoundType); ok {
		if !compoundType.HasType(outputType) {
			return userconfig.ErrorUnsupportedOutputType(outputType, compoundType)
		}
		return nil
	}

	// array of *InputSchema
	if inputSchemas, ok := cast.InterfaceToInterfaceSlice(typeSchema); ok {
		outputTypes, ok := cast.InterfaceToInterfaceSlice(outputType)
		if !ok {
			return userconfig.ErrorUnsupportedOutputType(outputType, inputSchemas)
		}

		err := validateInputRuntimeOutputTypes(outputTypes[0], inputSchemas[0].(*InputSchema))
		if err != nil {
			return errors.Wrap(err, s.Index(0))
		}
		return nil
	}

	// Map
	if typeSchemaMap, ok := cast.InterfaceToInterfaceInterfaceMap(typeSchema); ok {
		outputTypeMap, ok := cast.InterfaceToInterfaceInterfaceMap(outputType)
		if !ok {
			return userconfig.ErrorUnsupportedOutputType(outputType, typeSchemaMap)
		}

		var typeSchemaGenericKey CompoundType
		var typeSchemaGenericValue *InputSchema
		for k, v := range typeSchemaMap {
			ok := false
			if typeSchemaGenericKey, ok = k.(CompoundType); ok {
				typeSchemaGenericValue = v.(*InputSchema)
			}
		}

		var outputTypeGenericKey ValueType
		var outputTypeGenericValue interface{}
		for k, v := range outputTypeMap {
			ok := false
			if outputTypeGenericKey, ok = k.(ValueType); ok {
				outputTypeGenericValue = v
			}
		}

		// Check length if fixed outputType
		if outputTypeGenericValue == nil {
			if schema.MinCount != nil && int64(len(outputTypeMap)) < *schema.MinCount {
				return nil, ErrorTooFewElements(configreader.PrimTypeMap, *schema.MinCount)
			}
			if schema.MaxCount != nil && int64(len(outputTypeMap)) > *schema.MaxCount {
				return nil, ErrorTooManyElements(configreader.PrimTypeMap, *schema.MaxCount)
			}
		}

		// Generic schema map and generic outputType
		if typeSchemaGenericValue != nil && outputTypeGenericValue != nil {
			if err := validateInputRuntimeOutputTypes(outputTypeGenericKey, &userconfig.InputSchema{Type: typeSchemaGenericKey}); err != nil {
				return err
			}
			if err := validateInputRuntimeOutputTypes(outputTypeGenericValue, typeSchemaGenericValue); err != nil {
				return errors.Wrap(err, s.UserStrStripped(outputTypeGenericKey))
			}
			return nil
		}

		// Generic schema map and fixed outputType (we'll check the types of the fixed map)
		if typeSchemaGenericValue != nil && outputTypeGenericValue == nil {
			for outputTypeKey, outputTypeValue := range outputTypeMap {
				if _, err := typeSchemaGenericKey.CastValue(outputTypeKey); err != nil {
					return err
				}
				if err := validateInputRuntimeOutputTypes(outputTypeValue, typeSchemaGenericValue); err != nil {
					return errors.Wrap(err, s.UserStrStripped(outputTypeKey))
				}
			}
			return nil
		}

		// Generic outputType map and fixed schema map (we'll let this slide if the types match, as Python will validate the actual inputs)
		if typeSchemaGenericValue == nil && outputTypeGenericValue != nil {
			for typeSchemaKey, typeSchemaValue := range typeSchemaMap {
				if _, err := outputTypeGenericKey.CastValue(typeSchemaKey); err != nil {
					return err
				}
				if err := validateInputRuntimeOutputTypes(outputTypeGenericValue, typeSchemaValue.(*InputSchema)); err != nil {
					return errors.Wrap(err, s.UserStrStripped(typeSchemaKey))
				}
			}
			return nil
		}

		// Fixed outputType map and fixed schema map
		if typeSchemaGenericValue == nil && outputTypeGenericValue == nil {
			for typeSchemaKey, typeSchemaValue := range typeSchemaMap {
				if err := validateInputRuntimeOutputTypes(outputTypeMap[typeSchemaKey], typeSchemaValue.(*InputSchema)); err != nil {
					return nil, errors.Wrap(err, s.UserStrStripped(typeSchemaKey))
				}
			}
			for valueKey := range outputTypeMap {
				if _, ok := typeSchemaMap[valueKey]; !ok {
					return nil, ErrorUnsupportedLiteralMapKey(valueKey, typeSchemaMap)
				}
			}
			return nil
		}
	}

	return nil, userconfig.ErrorInvalidInputType(typeSchema) // unexpected
}

func replaceResourceIDsAndValidateResourceTypes(
	input interface{},
	validResourceTypes []resource.Type,
	validResources map[string]context.Resource,
	allResources map[string][]userconfig.Resource,
) (interface{}, interface{}, error) {

	if resourceName, ok := configreader.ExtractResourceName(input); ok {
		if res, ok := validResources[resourceName]; ok {
			return res.GetID(), res.GetIDWithTags(), nil
		}

		if len(allResources[resourceName] > 0) {
			return nil, nil, userconfig.ErrorResourceWrongType(allResources[resourceName], validResourceTypes...)
		}

		return nil, nil, userconfig.ErrorUndefinedResource(resourceName, validResourceTypes...)
	}

	if inputSlice, ok := cast.InterfaceToInterfaceSlice(input); ok {
		sliceWithIDs := make([]interface{}, len(inputSlice))
		sliceWithTagIDs := make([]interface{}, len(inputSlice))
		for i, elem := range inputSlice {
			elemWithIDs, elemWithTagIDs, err := replaceResourceIDs(elem, validResourceTypes, validResources, allResources)
			if err != nil {
				return nil, nil, errors.Wrap(err, s.Index(i))
			}
			sliceWithIDs[i] = elemWithIDs
			sliceWithTagIDs[i] = elemWithTagIDs
		}
		return sliceWithIDs, sliceWithTagIDs, nil
	}

	if inputMap, ok := cast.InterfaceToInterfaceInterfaceMap(input); ok {
		mapWithIDs := make(map[interface{}]interface{}, len(inputMap))
		mapWithTagIDs := make(map[interface{}]interface{}, len(inputMap))
		for key, val := range inputMap {
			keyWithIDs, keyWithTagIDs, err := replaceResourceIDs(key, validResourceTypes, validResources, allResources)
			if err != nil {
				return nil, nil, err
			}
			valWithIDs, valWithTagIDs, err := replaceResourceIDs(val, validResourceTypes, validResources, allResources)
			if err != nil {
				return nil, nil, errors.Wrap(err, s.UserStrStripped(key))
			}
			mapWithIDs[keyWithIDs] = valWithIDs
			mapWithTagIDs[keyWithTagIDs] = valWithTagIDs
		}
		return mapWithIDs, mapWithTagIDs, nil
	}

	return input, input, nil
}
