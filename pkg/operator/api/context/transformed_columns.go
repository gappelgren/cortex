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
	"github.com/cortexlabs/cortex/pkg/lib/cast"
	"github.com/cortexlabs/cortex/pkg/lib/sets/strset"
	"github.com/cortexlabs/cortex/pkg/operator/api/userconfig"
)

type TransformedColumns map[string]*TransformedColumn

type TransformedColumn struct {
	*userconfig.TransformedColumn
	*ComputedResourceFields
	Type userconfig.ColumnType `json:"type"`
}

func (column *TransformedColumn) GetType() userconfig.ColumnType {
	return column.Type
}

// Returns map[string]string because after autogen, arg values are constant or aggregate names
func (column *TransformedColumn) Args() map[string]string {
	args, _ := cast.InterfaceToStrStrMap(column.Inputs.Args)
	return args
}

func (column *TransformedColumn) InputAggregateNames(ctx *Context) strset.Set {
	inputAggregateNames := strset.New()
	for _, valueResourceName := range column.Args() {
		if _, ok := ctx.Aggregates[valueResourceName]; ok {
			inputAggregateNames.Add(valueResourceName)
		}
	}
	return inputAggregateNames
}

func (columns TransformedColumns) OneByID(id string) *TransformedColumn {
	for _, transformedColumn := range columns {
		if transformedColumn.ID == id {
			return transformedColumn
		}
	}
	return nil
}

func (column *TransformedColumn) GetInputRawColumnNames() []string {
	return column.InputColumnNames()
}
