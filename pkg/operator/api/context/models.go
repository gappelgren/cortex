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
	"github.com/cortexlabs/cortex/pkg/operator/api/resource"
	"github.com/cortexlabs/cortex/pkg/operator/api/userconfig"
)

type Models map[string]*Model
type TrainingDatasets map[string]*TrainingDataset

type Model struct {
	*userconfig.Model
	*ComputedResourceFields
	Key     string           `json:"key"`
	Dataset *TrainingDataset `json:"dataset"`
}

type TrainingDataset struct {
	userconfig.ResourceFields
	*ComputedResourceFields
	ModelName string `json:"model_name"`
	TrainKey  string `json:"train_key"`
	EvalKey   string `json:"eval_key"`
}

func (trainingDataset *TrainingDataset) GetResourceType() resource.Type {
	return resource.TrainingDatasetType
}

func (models Models) OneByID(id string) *Model {
	for _, model := range models {
		if model.ID == id {
			return model
		}
	}
	return nil
}

func (ctx *Context) OneTrainingDatasetByID(id string) *TrainingDataset {
	for _, model := range ctx.Models {
		if model.Dataset.ID == id {
			return model.Dataset
		}
	}
	return nil
}

func (models Models) GetTrainingDatasets() TrainingDatasets {
	trainingDatasets := make(map[string]*TrainingDataset, len(models))
	for _, model := range models {
		trainingDatasets[model.Dataset.Name] = model.Dataset
	}
	return trainingDatasets
}

func (model *Model) GetRawColumnInputs() []RawColumn {
	combinedInput := []interface{}{model.Input, model.TrainingInput}
	rawColResources := ExtractCortexResources(combinedInput, resource.RawColumnType)
	rawCols := make([]RawColumn, rawColResources)
	for i, rawColResource := range rawColResources {
		rawCols[i] = rawColResource.(RawColumn)
	}
	return rawCols
}
