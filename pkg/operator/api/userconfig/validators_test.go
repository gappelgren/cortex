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

	cr "github.com/cortexlabs/cortex/pkg/lib/configreader"
	// "github.com/cortexlabs/cortex/pkg/lib/debug"
)

// func TestExtractResourceNames(t *testing.T) {
// 	var resourceNames strset.Set
// 	var err error

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`@test`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`[@test1, @test2]`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test1", "test2"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`{@test1: @test2}`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test1", "test2"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`{test1: @test2}`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test2"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`{1: @test2}`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test2"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`[@test1, @test2, @test1, "@testStr"]`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test1", "test2"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`["@testStr"]`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New())

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`"@testStr"`))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New())

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`
//      map: {a: @test1, b: 3}
//      str: testStr
//      floats: [2.2, 3.3, @test2]
//      map2:
//        testA:
//          @test3: @test4
//          lon:
//            a: @test5
//            1: [@test6, @test7, "@testStr"]
//            @test8: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
//          bools: [@test9]
//     `))
// 	require.NoError(t, err)
// 	require.Equal(t, resourceNames, strset.New("test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9"))

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`@test?test`))
// 	require.Error(t, err)

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`@test test`))
// 	require.Error(t, err)

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`[@test test]`))
// 	require.Error(t, err)

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`{@test test: test}`))
// 	require.Error(t, err)

// 	resourceNames, err = ExtractResourceNames(cr.MustReadYAMLStr(
// 		`{test: @test test}`))
// 	require.Error(t, err)
// }

func TestValidateOutputTypeSchema(t *testing.T) {
	var err error

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`STRING`))
	require.NoError(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`STRING|INT`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`STRING_COLUMN`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`bad`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`[STRING]`))
	require.NoError(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`[STRING, INT]`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`[STRING|INT]`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`[STRING_COLUMN]`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`[bad]`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{mean: FLOAT, stddev: FLOAT}`))
	require.NoError(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{_type: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{_mean: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT: FLOAT}`))
	require.NoError(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT: INT|FLOAT}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT|FLOAT: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT: FLOAT_COLUMN}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT_COLUMN: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateOutputTypeSchema(cr.MustReadYAMLStr(
		`{INT: FLOAT, FLOAT: FLOAT}`))
	require.Error(t, err)
}

func checkCastConstantEqual(t *testing.T, outputSchemaYAML string, valueYAML string, expected interface{}) {
	outputSchema, err := ValidateOutputTypeSchema(cr.MustReadYAMLStr(outputSchemaYAML))
	require.NoError(t, err)
	casted, err := CastConstant(cr.MustReadYAMLStr(valueYAML), outputSchema)
	require.NoError(t, err)
	require.Equal(t, casted, expected)
}

func checkCastConstantError(t *testing.T, outputSchemaYAML string, valueYAML string) {
	outputSchema, err := ValidateOutputTypeSchema(cr.MustReadYAMLStr(outputSchemaYAML))
	require.NoError(t, err)
	_, err = CastConstant(cr.MustReadYAMLStr(valueYAML), outputSchema)
	require.Error(t, err)
}

func TestCastConstant(t *testing.T) {
	checkCastConstantEqual(t, `INT`, `2`, int64(2))
	checkCastConstantError(t, `INT`, `test`)
	checkCastConstantError(t, `INT`, `2.2`)
	checkCastConstantEqual(t, `FLOAT`, `2`, float64(2))
	checkCastConstantError(t, `FLOAT`, `test`)
	checkCastConstantEqual(t, `BOOL`, `true`, true)

	checkCastConstantEqual(t, `{STRING: FLOAT}`, `{test: 2.2, test2: 4.4}`,
		map[interface{}]interface{}{"test": 2.2, "test2": 4.4})
	checkCastConstantError(t, `{STRING: FLOAT}`, `{test: test2}`)
	checkCastConstantEqual(t, `{STRING: FLOAT}`, `{test: 2}`,
		map[interface{}]interface{}{"test": float64(2)})
	checkCastConstantEqual(t, `{STRING: INT}`, `{test: 2}`,
		map[interface{}]interface{}{"test": int64(2)})
	checkCastConstantError(t, `{STRING: INT}`, `{test: 2.0}`)

	checkCastConstantEqual(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2, sum: 4}`,
		map[interface{}]interface{}{"mean": float64(2.2), "sum": int64(4)})
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2, sum: test}`)
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: false, sum: 4}`)
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2, 2: 4}`)
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2, sum: Null}`)
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2}`)
	checkCastConstantError(t, `{mean: FLOAT, sum: INT}`, `{mean: 2.2, sum: 4, stddev: 2}`)

	checkCastConstantEqual(t, `[INT]`, `[1, 2]`,
		[]interface{}{int64(1), int64(2)})
	checkCastConstantError(t, `[INT]`, `[1.0, 2]`)
	checkCastConstantEqual(t, `[FLOAT]`, `[1.0, 2]`,
		[]interface{}{float64(1), float64(2)})

	outputSchemaYAML :=
		`
     map: {STRING: FLOAT}
     str: STRING
     floats: [FLOAT]
     map2:
       STRING:
         lat: FLOAT
         lon:
           a: INT
           b: [STRING]
           c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
         bools: [BOOL]
    `
	outputSchema, err := ValidateOutputTypeSchema(cr.MustReadYAMLStr(outputSchemaYAML))
	require.NoError(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [true, false, true]
    `), outputSchema)
	require.NoError(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [true, false, true]
    `), outputSchema)
	require.Error(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88.8
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [true, false, true]
    `), outputSchema)
	require.Error(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88
           b: [testX, testY, 2]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [true, false, true]
    `), outputSchema)
	require.Error(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: test}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [true, false, true]
    `), outputSchema)
	require.Error(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: true
    `), outputSchema)
	require.Error(t, err)

	_, err = CastConstant(cr.MustReadYAMLStr(
		`
     map: {a: 2.2, b: 3}
     str: test1
     floats: [2.2, 3.3, 4.4]
     map2:
       testA:
         lat: 9.9
         lon:
           a: 17
           b: [test1, test2, test3]
           c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
         bools: [true]
       testB:
         lat: 3.14
         lon:
           a: 88
           b: [testX, testY, testZ]
           c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         bools: [1, 2, 3]
    `), outputSchema)
	require.Error(t, err)
}

func TestValidateInputSchema(t *testing.T) {
	var inputSchema, inputSchema2, inputSchema3, inputSchema4 *InputSchema
	var err error

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`STRING`))
	require.NoError(t, err)
	inputSchema2, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`_type: STRING`))
	require.NoError(t, err)
	require.Equal(t, inputSchema, inputSchema2)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: STRING
     _default: test
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: STRING_COLUMN
     _default: test
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: STRING
     _default: Null
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: STRING
     _default: 2
    `))
	require.Error(t, err)

	// Lists

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`[STRING]`))
	require.NoError(t, err)
	inputSchema2, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`_type: [STRING]`))
	require.NoError(t, err)
	inputSchema3, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       - _type: STRING
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema, inputSchema2)
	require.Equal(t, inputSchema, inputSchema3)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`[STRING|INT_COLUMN]`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1, test2, test3]
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING_COLUMN]
     _default: [test1, test2, test3]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1, 2, test3]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING|INT]
     _default: [test1, 2, test3]
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING|FLOAT]
     _default: [test1, 2, test3]
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: test1
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _min_count: 2
     _max_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _min_count: 2
     _max_count: 1
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1]
     _min_count: 2
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1, test2]
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _min_count: -1
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _min_count: test
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1, test2, test3]
     _max_count: 2
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [STRING]
     _default: [test1, test2]
     _max_count: 2
    `))
	require.NoError(t, err)

	// Maps

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`arg1: STRING`))
	require.NoError(t, err)
	inputSchema2, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: STRING
    `))
	require.NoError(t, err)
	inputSchema3, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {arg1: STRING}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema, inputSchema2)
	require.Equal(t, inputSchema, inputSchema3)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`_arg1: STRING`))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`arg1: test`))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`STRING: test`))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`STRING_COLUMN: test`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {arg1: STRING}
     _min_count: 2
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT}
     _default: {test: 2}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {FLOAT: INT}
     _default: {2: 2}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT}
     _default: {test: test}
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT|STRING}
     _default: {test: test}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT_COLUMN}
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING_COLUMN: INT}
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING_COLUMN: INT}
     _default: {test: 2}
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING_COLUMN: INT_COLUMN}
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING_COLUMN: INT_COLUMN|STRING}
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING_COLUMN: INT_COLUMN|STRING_COLUMN}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT}
     _min_count: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: STRING
       _optional: true
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: STRING
       _default: test
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: STRING
       _default: 2
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: STRING
       _default: Null
    `))
	require.Error(t, err)

	// Mixed

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`[[STRING]]`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [[STRING]]
     _default: [[test1, test2]]
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [[STRING_COLUMN]]
     _default: [[test1, test2]]
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     - arg1: STRING
       arg2: INT
    `))
	require.NoError(t, err)
	inputSchema2, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
     - arg1: STRING
       arg2: INT
    `))
	require.NoError(t, err)
	inputSchema3, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     - arg1: {_type: STRING}
       arg2: {_type: INT}
    `))
	require.NoError(t, err)
	inputSchema4, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
     - arg1:
         _type: STRING
       arg2:
         _type: INT
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema, inputSchema2)
	require.Equal(t, inputSchema, inputSchema3)
	require.Equal(t, inputSchema, inputSchema4)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     - arg1:
         _type: STRING
         _default: test
       arg2:
         _type: INT
         _default: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     - arg1:
         _type:
           arg_a: STRING
           arg_b:
             _type: INT
             _default: 1
       arg2:
         _type: INT
         _default: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg_a:
       arg1:
         _type:
           arg_a: STRING
           arg_b:
             _type: INT
             _default: 1
       arg2:
         _type: INT
         _default: 2
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       INT:
         arg_a: INT
         arg_b:
           _type: STRING
           _default: test
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     - arg1:
         INT:
           arg_a: INT
           arg_b:
             _type: STRING
             _default: test
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
     - INT:
         arg_a: INT
         arg_b:
           _type: STRING
           _default: test
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       2: STRING
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       _type: {2: STRING}
       _default: {2: test}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     arg1:
       2:
         _type: STRING
         _default: test
    `))
	require.NoError(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`[{INT_COLUMN: STRING|INT}]`))
	require.NoError(t, err)
	inputSchema2, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [{INT_COLUMN: STRING|INT}]
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema, inputSchema2)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {BOOL|FLOAT: INT|STRING}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {mean: FLOAT, stddev: FLOAT}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {lat: FLOAT, lon: FLOAT}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {lat: FLOAT, lon: [FLOAT]}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {FLOAT: INT}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {FLOAT: [INT]}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {lat: FLOAT, lon: {lat2: FLOAT, lon2: INT}}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map6: {STRING: {lat: FLOAT, lon: {lat2: FLOAT, lon2: {INT: STRING}}}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map6: {STRING: {lat: FLOAT, lon: {lat2: FLOAT, lon2: {INT: STRING}, mean: BOOL}}}`))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     num: [INT]
     str: STRING
     map1: {STRING: INT}
     map2: {mean: FLOAT, stddev: FLOAT}
     map3: {STRING: {lat: FLOAT, lon: FLOAT}}
     map3: {STRING: {lat: FLOAT, lon: [FLOAT]}}
     map4: {STRING: {FLOAT: INT}}
     map5: {STRING: {BOOL: [INT]}}
     map6: {STRING: {lat: FLOAT, lon: {lat2: FLOAT, lon2: INT}}}
     map6: {STRING: {lat: FLOAT, lon: {lat2: FLOAT, lon2: {INT: STRING}, mean: BOOL}}}
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: INT, INT: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: INT, INT: [FLOAT]}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {mean: FLOAT, INT: FLOAT}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {mean: FLOAT, INT: [FLOAT]}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {lat: FLOAT, STRING: FLOAT}}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`map: {STRING: {STRING: test}}`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: [STRING_COLUMN, INT_COLUMN]`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: [STRING_COLUMNs]`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: [STRING_COLUMN|BAD]`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: Null`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: 1`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: [1]`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`cols: []`))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     float: FLOAT_COLUMN
     int: INT_COLUMN
     str: STRING_COLUMN
     int_list: FLOAT_LIST_COLUMN
     float_list: INT_LIST_COLUMN
     str_list: STRING_LIST_COLUMN
     num1: FLOAT_COLUMN|INT_COLUMN
     num2: INT_COLUMN|FLOAT_COLUMN
     num3: STRING_COLUMN|INT_COLUMN
     num4: INT_COLUMN|FLOAT_COLUMN|STRING_COLUMN
     num5: STRING_COLUMN|INT_COLUMN|FLOAT_COLUMN
     num6: STRING_LIST_COLUMN|INT_LIST_COLUMN|FLOAT_LIST_COLUMN
     num7: STRING_COLUMN|INT_LIST_COLUMN|FLOAT_LIST_COLUMN
     nums1: [INT_COLUMN]
     nums2: [FLOAT_COLUMN]
     nums3: [INT_COLUMN|FLOAT_COLUMN]
     nums4: [FLOAT_COLUMN|INT_COLUMN]
     nums5: [STRING_COLUMN|INT_COLUMN|FLOAT_COLUMN]
     nums6: [INT_LIST_COLUMN]
     nums7: [INT_LIST_COLUMN|STRING_LIST_COLUMN]
     nums8: [INT_LIST_COLUMN|STRING_COLUMN]
     float_scalar: FLOAT
     int_scalar: INT
     str_scalar: STRING
     bool_scalar: BOOL
     num1_scalar: FLOAT|INT
     num2_scalar: INT|FLOAT
     num3_scalar: STRING|INT
     num4_scalar: INT|FLOAT|STRING
     num5_scalar: STRING|INT|FLOAT
     nums1_scalar: [INT]
     nums2_scalar: [FLOAT]
     nums3_scalar: [INT|FLOAT]
     nums4_scalar: [FLOAT|INT]
     nums5_scalar: [STRING|INT|FLOAT]
     nums6_scalar: [STRING|INT|FLOAT|BOOL]
    `))
	require.NoError(t, err)

	// Casting defaults

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: INT
     _default: 2
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, int64(2))

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: INT
     _default: test
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: INT
     _default: 2.2
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: FLOAT
     _default: 2
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, float64(2))

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: FLOAT|INT
     _default: 2
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, int64(2))

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: BOOL
     _default: true
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, true)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: FLOAT}
     _default: {test: 2.2, test2: 4.4}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{"test": 2.2, "test2": 4.4})

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: FLOAT}
     _default: {test: test2}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: FLOAT}
     _default: {test: 2}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{"test": float64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: FLOAT}
     _default: {test: 2.0}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{"test": float64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT}
     _default: {test: 2}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{"test": int64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {STRING: INT}
     _default: {test: 2.0}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2, sum: 4}
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{"mean": float64(2.2), "sum": int64(4)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2, sum: test}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: false, sum: 4}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2, 2: 4}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2, sum: Null}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: {mean: FLOAT, sum: INT}
     _default: {mean: 2.2, sum: 4, stddev: 2}
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [INT]
     _default: [1, 2]
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, []interface{}{int64(1), int64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [INT]
     _default: [1.0, 2]
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [FLOAT]
     _default: [1.0, 2]
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, []interface{}{float64(1), float64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [FLOAT|INT]
     _default: [1.0, 2]
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, []interface{}{float64(1), int64(2)})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [FLOAT|INT|BOOL]
     _default: [1.0, 2, true, test]
    `))
	require.Error(t, err)

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type: [FLOAT|INT|BOOL|STRING]
     _default: [1.0, 2, true, test]
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, []interface{}{float64(1), int64(2), true, "test"})

	inputSchema, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       STRING:
         a:
           _type: INT
           _optional: true
         b:
           _type: [STRING]
           _optional: true
         c:
           _type: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           _optional: true
           _optional: true
         d:
           _type: INT
           _default: 2
     _default:
       testA: {}
       testB:
         a: 88
         b: [testX, testY, testZ]
         c: {mean: 1.7, sum: [1], stddev: {z: 12}}
         d: 17
    `))
	require.NoError(t, err)
	require.Equal(t, inputSchema.Default, map[interface{}]interface{}{
		"testA": map[interface{}]interface{}{
			"a": nil,
			"b": nil,
			"c": nil,
			"d": int64(2),
		},
		"testB": map[interface{}]interface{}{
			"a": int64(88),
			"b": []interface{}{"testX", "testY", "testZ"},
			"c": map[interface{}]interface{}{
				"mean":   float64(1.7),
				"sum":    []interface{}{int64(1)},
				"stddev": map[interface{}]interface{}{"z": int64(12)},
			},
			"d": int64(17),
		},
	})

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [true, false, true]
           anything: [10, 2.2, test, false]
    `))
	require.NoError(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [true, false, true]
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88.8
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [true, false, true]
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88
             b: [testX, testY, 2]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [true, false, true]
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: test}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [true, false, true]
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: true
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)

	_, err = ValidateInputSchema(cr.MustReadYAMLStr(
		`
     _type:
       map: {STRING: FLOAT}
       str: STRING
       floats: [FLOAT]
       map2:
         STRING:
           lat: FLOAT
           lon:
             a: INT
             b: [STRING]
             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT}}
           bools: [BOOL]
           anything: [BOOL|INT|FLOAT|STRING]
     _default:
       map: {a: 2.2, b: 3}
       str: test1
       floats: [2.2, 3.3, 4.4]
       map2:
         testA:
           lat: 9.9
           lon:
             a: 17
             b: [test1, test2, test3]
             c: {mean: 8.8, sum: [3, 2, 1], stddev: {a: 1, b: 2}}
           bools: [true]
           anything: []
         testB:
           lat: 3.14
           lon:
             a: 88
             b: [testX, testY, testZ]
             c: {mean: 1.7, sum: [1], stddev: {z: 12}}
           bools: [1, 2, 3]
           anything: [10, 2.2, test, false]
    `))
	require.Error(t, err)
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//

// func TestCheckColumnRuntimeTypesMatch(t *testing.T) {
// 	var columnTypes map[string]interface{}
// 	var runtimeTypes map[string]interface{}

// 	columnTypes = cr.MustReadYAMLStrMap("in: INT_COLUMN")
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: FLOAT_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in: INT_COLUMN|FLOAT_COLUMN")
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: FLOAT_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: STRING_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in: STRING_COLUMN|INT_COLUMN|FLOAT_COLUMN")
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: FLOAT_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: STRING_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: BAD_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in: [INT_COLUMN]")
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, INT_COLUMN, INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [FLOAT_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, FLOAT_COLUMN, INT_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in: [INT_COLUMN|FLOAT_COLUMN]")
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, INT_COLUMN, INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [FLOAT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, FLOAT_COLUMN, INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, FLOAT_COLUMN, STRING_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in: [STRING_COLUMN|INT_COLUMN|FLOAT_COLUMN]")
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, INT_COLUMN, INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: INT_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [BAD_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [STRING_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in: [INT_COLUMN, FLOAT_COLUMN, STRING_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in1: [INT_COLUMN]\nin2: STRING_COLUMN")
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN]\nin2: STRING_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN, INT_COLUMN]\nin2: STRING_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in2: STRING_COLUMN\nin1: [INT_COLUMN, INT_COLUMN]")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN]")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN]\nin2: STRING_COLUMN\nin3: INT_COLUMN")
// 	require.Error(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))

// 	columnTypes = cr.MustReadYAMLStrMap("in1: [INT_COLUMN|FLOAT_COLUMN|STRING_COLUMN]\nin2: STRING_COLUMN")
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN]\nin2: STRING_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// 	runtimeTypes = readRuntimeTypes("in1: [INT_COLUMN, FLOAT_COLUMN, STRING_COLUMN, FLOAT_COLUMN]\nin2: STRING_COLUMN")
// 	require.NoError(t, CheckColumnRuntimeTypesMatch(runtimeTypes, columnTypes))
// }

// func readRuntimeTypes(yamlStr string) map[string]interface{} {
// 	runtimeTypes := make(map[string]interface{})
// 	runtimeTypesStr := cr.MustReadYAMLStrMap(yamlStr)

// 	for k, v := range runtimeTypesStr {
// 		if runtimeTypeStr, ok := v.(string); ok {
// 			runtimeTypes[k] = ColumnTypeFromString(runtimeTypeStr)
// 		} else if runtimeTypeStrs, ok := cast.InterfaceToStrSlice(v); ok {
// 			runtimeTypesSlice := make([]ColumnType, len(runtimeTypeStrs))
// 			for i, runtimeTypeStr := range runtimeTypeStrs {
// 				runtimeTypesSlice[i] = ColumnTypeFromString(runtimeTypeStr)
// 			}
// 			runtimeTypes[k] = runtimeTypesSlice
// 		}
// 	}

// 	return runtimeTypes
// }

// func TestValidateArgTypes(t *testing.T) {
// 	var argTypes map[string]interface{}

// 	argTypes = cr.MustReadYAMLStrMap("STRING: FLOAT")
// 	require.Error(t, ValidateArgTypes(argTypes))
// }

// func TestCheckValueRuntimeTypesMatch(t *testing.T) {
// 	var schemaType interface{}
// 	var runtimeType interface{}

// 	schemaType = "INT"
// 	runtimeType = "INT"
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "FLOAT"
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "FLOAT|INT"
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = "FLOAT|INT"
// 	runtimeType = "FLOAT|INT"
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "INT|FLOAT"
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "FLOAT"
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "INT"
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "STRING"
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = []string{"INT"}
// 	runtimeType = []string{"INT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"STRING"}
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = []string{"BOOL"}
// 	runtimeType = []string{"BOOL"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"INT"}
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = []string{"FLOAT|INT"}
// 	runtimeType = []string{"INT|FLOAT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"INT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = []string{"BOOL|FLOAT|INT"}
// 	runtimeType = []string{"FLOAT|INT|BOOL"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"INT|FLOAT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"FLOAT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"INT"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"BOOL"}
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = []string{"STRING"}
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = "FLOAT"
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap("STRING: FLOAT")
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: FLOAT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: INT")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap("STRING: [INT|FLOAT]")
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: [FLOAT|INT]")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: [INT]")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: [BOOL]")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("STRING: INT")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap("INT|FLOAT: STRING")
// 	runtimeType = cr.MustReadYAMLStrMap("FLOAT|INT: STRING")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("INT: STRING")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("BOOL: STRING")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap("mean: FLOAT\nsum: INT")
// 	runtimeType = cr.MustReadYAMLStrMap("mean: FLOAT\nsum: INT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("sum: INT\nmean: FLOAT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("sum: INT\nmean: INT")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("mean: FLOAT")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("mean: FLOAT\nsum: INT\nextra: STRING")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap("mean: FLOAT\nsum: INT|FLOAT")
// 	runtimeType = cr.MustReadYAMLStrMap("mean: FLOAT\nsum: FLOAT|INT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("sum: FLOAT\nmean: FLOAT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("sum: INT\nmean: FLOAT")
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap("sum: INT\nmean: INT")
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))

// 	schemaType = cr.MustReadYAMLStrMap(
// 		`
//       map: {STRING: FLOAT}
//       str: STRING
//       floats: [FLOAT]
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             a: INT|FLOAT
//             b: [STRING]
//             c: {mean: FLOAT, sum: [INT], stddev: {STRING: INT|FLOAT}}
//             d: [BOOL]
//      `)
// 	runtimeType = cr.MustReadYAMLStrMap(
// 		`
//       floats: [FLOAT]
//       str: STRING
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             c: {sum: [INT], mean: FLOAT, stddev: {STRING: FLOAT|INT}}
//             b: [STRING]
//             a: FLOAT|INT
//             d: [BOOL]
//       map: {STRING: FLOAT}
//      `)
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap(
// 		`
//       floats: [FLOAT]
//       str: STRING
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             c: {sum: [INT], mean: FLOAT, stddev: {STRING: FLOAT|INT}}
//             b: [STRING]
//             a: INT
//             d: [BOOL]
//       map: {STRING: FLOAT}
//      `)
// 	require.NoError(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap(
// 		`
//       floats: [FLOAT]
//       str: STRING
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             c: {sum: [INT], mean: FLOAT, stddev: {STRING: FLOAT|INT}}
//             b: STRING
//             a: FLOAT|INT
//             d: [BOOL]
//       map: {STRING: FLOAT}
//      `)
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap(
// 		`
//       floats: [FLOAT]
//       str: STRING
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             c: {sum: [INT], stddev: {STRING: FLOAT|INT}}
//             b: [STRING]
//             a: FLOAT|INT
//             d: [BOOL]
//       map: {STRING: FLOAT}
//      `)
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// 	runtimeType = cr.MustReadYAMLStrMap(
// 		`
//       floats: [FLOAT]
//       str: STRING
//       map2:
//         STRING:
//           lat: FLOAT
//           lon:
//             c: {sum: [INT], mean: FLOAT, stddev: {STRING: FLOAT|INT}}
//             b: [STRING]
//             a: FLOAT|INT
//             d: BOOL
//       map: {STRING: FLOAT}
//      `)
// 	require.Error(t, CheckValueRuntimeTypesMatch(runtimeType, schemaType))
// }
