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
	"bytes"
	"path/filepath"

	"github.com/cortexlabs/cortex/pkg/consts"
	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/lib/hash"
	"github.com/cortexlabs/cortex/pkg/operator/api/context"
	"github.com/cortexlabs/cortex/pkg/operator/api/resource"
	"github.com/cortexlabs/cortex/pkg/operator/api/userconfig"
)

func getAggregates(
	config *userconfig.Config,
	constants context.Constants,
	rawColumns context.RawColumns,
	userAggregators map[string]*context.Aggregator,
	root string,
) (context.Aggregates, error) {

	aggregates := context.Aggregates{}

	for _, aggregateConfig := range config.Aggregates {
		aggregator, err := getAggregator(aggregateConfig.Aggregator, userAggregators)
		if err != nil {
			return nil, errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.AggregatorKey)
		}

		validResources := make(map[string]context.Resource)
		for name, res := range constants {
			validResources[name] = res
		}
		for name, res := range rawColumns {
			validResources[name] = res
		}

		castedInput, inputID, inputIDWithTags, err := ValidateInput(
			aggregateConfig.Input,
			aggregator.Input,
			[]resource.Type{resource.RawColumnType, resource.ConstantType},
			validResources,
			config.Resources,
			userAggregators,
			nil,
		)
		if err != nil {
			return errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.InputKey)
		}
		aggregateConfig.Input = castedInput

		var buf bytes.Buffer
		buf.WriteString(inputID)
		buf.WriteString(aggregator.ID)
		id := hash.Bytes(buf.Bytes())

		buf.Reset()
		buf.WriteString(inputIDWithTags)
		buf.WriteString(aggregator.IDWithTags)
		buf.WriteString(aggregateConfig.Tags.ID())
		idWithTags := hash.Bytes(buf.Bytes())

		aggregateRootKey := filepath.Join(
			root,
			consts.AggregatesDir,
			id,
		)
		aggregateKey := aggregateRootKey + ".msgpack"
		aggregateMetadataKey := aggregateRootKey + "_metadata.json"

		aggregates[aggregateConfig.Name] = &context.Aggregate{
			ComputedResourceFields: &context.ComputedResourceFields{
				ResourceFields: &context.ResourceFields{
					ID:           id,
					IDWithTags:   idWithTags,
					ResourceType: resource.AggregateType,
					MetadataKey:  aggregateMetadataKey,
				},
			},
			Aggregate: aggregateConfig,
			Type:      aggregator.OutputType,
			Key:       aggregateKey,
		}
	}

	return aggregates, nil
}

func validateAggregateInputs(
	aggregateConfig *userconfig.Aggregate,
	constants context.Constants,
	rawColumns context.RawColumns,
	aggregator *context.Aggregator,
) error {
	if aggregateConfig.AggregatorPath != nil {
		return nil
	}

	columnRuntimeTypes, err := context.GetColumnRuntimeTypes(aggregateConfig.Inputs.Columns, rawColumns)
	if err != nil {
		return errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.InputsKey, userconfig.ColumnsKey)
	}
	err = userconfig.CheckColumnRuntimeTypesMatch(columnRuntimeTypes, aggregator.Inputs.Columns)
	if err != nil {
		return errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.InputsKey, userconfig.ColumnsKey)
	}

	argTypes, err := getAggregateArgTypes(aggregateConfig.Inputs.Args, constants)
	if err != nil {
		return errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.InputsKey, userconfig.ArgsKey)
	}
	err = userconfig.CheckArgRuntimeTypesMatch(argTypes, aggregator.Inputs.Args)
	if err != nil {
		return errors.Wrap(err, userconfig.Identify(aggregateConfig), userconfig.InputsKey, userconfig.ArgsKey)
	}

	return nil
}

func getAggregateArgTypes(
	args map[string]interface{},
	constants context.Constants,
) (map[string]interface{}, error) {

	if len(args) == 0 {
		return nil, nil
	}

	argTypes := make(map[string]interface{}, len(args))
	for argName, constantName := range args {
		constantNameStr := constantName.(string)
		constant, ok := constants[constantNameStr]
		if !ok {
			return nil, errors.Wrap(userconfig.ErrorUndefinedResource(constantNameStr, resource.ConstantType), argName)
		}
		argTypes[argName] = constant.Type
	}
	return argTypes, nil
}
