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
	"github.com/cortexlabs/yaml"

	"github.com/cortexlabs/cortex/pkg/lib/sets/strset"
	"github.com/cortexlabs/cortex/pkg/operator/api/resource"
)

func (ctx *Context) AllComputedResourceDependencies(resourceID string) strset.Set {
	dependencies := ctx.DirectComputedResourceDependencies(resourceID)
	for dependency := range dependencies.Copy() {
		for subDependency := range ctx.AllComputedResourceDependencies(dependency) {
			dependencies.Add(subDependency)
		}
	}
	return dependencies
}

func (ctx *Context) DirectComputedResourceDependencies(resourceID string) strset.Set {
	for _, pythonPackage := range ctx.PythonPackages {
		if pythonPackage.GetID() == resourceID {
			return ctx.pythonPackageDependencies(pythonPackage)
		}
	}
	for _, rawColumn := range ctx.RawColumns {
		if rawColumn.GetID() == resourceID {
			return ctx.rawColumnDependencies(rawColumn)
		}
	}
	for _, aggregate := range ctx.Aggregates {
		if aggregate.ID == resourceID {
			return ctx.aggregatesDependencies(aggregate)
		}
	}
	for _, transformedColumn := range ctx.TransformedColumns {
		if transformedColumn.ID == resourceID {
			return ctx.transformedColumnDependencies(transformedColumn)
		}
	}
	for _, model := range ctx.Models {
		if model.ID == resourceID {
			return ctx.modelDependencies(model)
		}
		if model.Dataset.ID == resourceID {
			return ctx.trainingDatasetDependencies(model)
		}
	}
	for _, api := range ctx.APIs {
		if api.ID == resourceID {
			return ctx.apiDependencies(api)
		}
	}
	return strset.New()
}

func (ctx *Context) pythonPackageDependencies(pythonPackage *PythonPackage) strset.Set {
	return strset.New()
}

func (ctx *Context) rawColumnDependencies(rawColumn RawColumn) strset.Set {
	// Currently python packages are a dependency on raw features because raw features share
	// the same workload as transformed features and aggregates.
	dependencies := strset.New()
	for _, pythonPackage := range ctx.PythonPackages {
		dependencies.Add(pythonPackage.GetID())
	}
	return dependencies
}

func (ctx *Context) aggregatesDependencies(aggregate *Aggregate) strset.Set {
	dependencies := strset.New()

	for _, pythonPackage := range ctx.PythonPackages {
		dependencies.Add(pythonPackage.GetID())
	}

	inputs := ExtractCortexResources(aggregate.Input)
	for _, res := range inputs {
		dependencies.Add(res.GetID())
	}

	return dependencies
}

func (ctx *Context) transformedColumnDependencies(transformedColumn *TransformedColumn) strset.Set {
	dependencies := strset.New()

	for _, pythonPackage := range ctx.PythonPackages {
		dependencies.Add(pythonPackage.GetID())
	}

	inputs := ExtractCortexResources(transformedColumn.Input)
	for _, res := range inputs {
		dependencies.Add(res.GetID())
	}

	return dependencies
}

func (ctx *Context) trainingDatasetDependencies(model *Model) strset.Set {
	dependencies := strset.New()
	for _, columnName := range model.AllColumnNames() {
		column := ctx.GetColumn(columnName)
		dependencies.Add(column.GetID())
	}
	return dependencies
}

func (ctx *Context) modelDependencies(model *Model) strset.Set {
	dependencies := strset.New()

	for _, pythonPackage := range ctx.PythonPackages {
		dependencies.Add(pythonPackage.GetID())
	}

	dependencies.Add(model.Dataset.ID)

	inputs := ExtractCortexResources(model.Input)
	for _, res := range inputs {
		dependencies.Add(res.GetID())
	}

	targetColumn := ctx.GetColumn(model.TargetColumn)
	dependencies.Add(targetColumn.GetID())

	return dependencies
}

func (ctx *Context) apiDependencies(api *API) strset.Set {
	model := ctx.Models[api.ModelName]
	return strset.New(model.ID)
}

// resourceTypes indicates which resource types to include in the query; if none are passed in, no filter is applied
func (ctx *Context) ExtractCortexResources(input interface{}, resourceTypes ...resource.Type) []Resource {
	resourceTypeMap := make(map[resource.Type]bool)
	for _, resourceType := range resourceTypes {
		resourceTypeMap[resourceType] = true
	}

	resources := make(map[string]Resource)
	ctx.extractCortexResourcesHelper(input, resourceTypeMap, resources)

	// convert to slice and sort by ID
	var resourceIDs []string
	for resourceId := range resources {
		resourceIDs = append(resourceIDs, resourceId)
	}
	sort.Strings(resourceIDs)
	resoucesSlice := make([]Resource, len(resources))
	for i, resourceID := range resourceIDs {
		resoucesSlice[i] = resources[resourceID]
	}

	return resoucesSlice
}

func (ctx *Context) extractCortexResourcesHelper(
	input interface{},
	resourceTypeMap map[resource.Type]bool,
	resources map[string]Resource,
) {

	if resourceName, ok := yaml.ExtractAtSymbolText(input); ok {
		for _, res := range AllResourcesByName(resourceName) {
			if len(resourceTypeMap) == 0 || resourceTypeMap[res.GetType()] == true {
				resources[res.GetID()] = res
			}
		}
		return
	}

	if inputSlice, ok := cast.InterfaceToInterfaceSlice(input); ok {
		for _, elem := range inputSlice {
			extractCortexResourcesHelper(elem, resourceTypeMap, resources)
		}
		return
	}

	if inputMap, ok := cast.InterfaceToInterfaceInterfaceMap(input); ok {
		for key, val := range inputMap {
			extractCortexResourcesHelper(key, resourceTypeMap, resources)
			extractCortexResourcesHelper(val, resourceTypeMap, resources)
		}
		return
	}
}
