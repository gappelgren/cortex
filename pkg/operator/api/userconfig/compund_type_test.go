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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompoundTypeFromString(t *testing.T) {
	var err error

	_, err = CompoundTypeFromString("STRING")
	require.Nil(t, err)
	_, err = CompoundTypeFromString("INT")
	require.Nil(t, err)
	_, err = CompoundTypeFromString("BOOL|FLOAT")
	require.Nil(t, err)

	_, err = CompoundTypeFromString("STRING_COLUMN")
	require.Nil(t, err)
	_, err = CompoundTypeFromString("INT_COLUMN")
	require.Nil(t, err)
	_, err = CompoundTypeFromString("INT_COLUMN|FLOAT_COLUMN")
	require.Nil(t, err)

	_, err = CompoundTypeFromString("INT|INT")
	require.Error(t, err)

	_, err = CompoundTypeFromString("INT|FLOAT_COLUMN")
	require.Error(t, err)

	_, err = CompoundTypeFromString("test")
	require.Error(t, err)
	_, err = CompoundTypeFromString("|")
	require.Error(t, err)
	_, err = CompoundTypeFromString("")
	require.Error(t, err)
	_, err = CompoundTypeFromString(2)
	require.Error(t, err)

	require.Error(t, err)
}
