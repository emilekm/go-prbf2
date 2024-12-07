package prism2_test

import (
	"testing"

	"github.com/emilekm/go-prbf2/prism2"
	"github.com/stretchr/testify/require"
)

type testBasicType struct {
	Integer      int
	Uinteger     uint
	Float        float32
	Str          string
	SliceOfBytes []byte
}

type testSimpleType struct {
	Str string
}

type testComplexType struct {
	Arr           [3]string
	StructValue   testSimpleType
	StructPointer *testSimpleType
}

type testSliceSimpleType struct {
	SliceOfStrings []string
}

func TestUnmarshalMessage(t *testing.T) {
	tests := []struct {
		name   string
		rawMsg []byte
		output any
		into   any
	}{
		{
			name:   "basic type success",
			rawMsg: []byte("-123\x03123\x032.0\x03test-string\x03testbytes"),
			output: &testBasicType{
				Integer:      -123,
				Uinteger:     123,
				Float:        2.0,
				Str:          "test-string",
				SliceOfBytes: []byte("testbytes"),
			},
			into: &testBasicType{},
		},
		{
			name:   "complex type success",
			rawMsg: []byte("test1\x03test2\x03test3\x03firstSimple\x03secondSimple"),
			output: &testComplexType{
				Arr:           [3]string{"test1", "test2", "test3"},
				StructValue:   testSimpleType{"firstSimple"},
				StructPointer: &testSimpleType{"secondSimple"},
			},
			into: &testComplexType{},
		},
		{
			name:   "slice type success",
			rawMsg: []byte("test1\x03test2\x03test3"),
			output: &testSliceSimpleType{
				SliceOfStrings: []string{"test1", "test2", "test3"},
			},
			into: &testSliceSimpleType{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := prism2.Unmarshal(test.rawMsg, test.into)
			require.NoError(t, err)

			require.Equal(t, test.output, test.into)
		})
	}
}
