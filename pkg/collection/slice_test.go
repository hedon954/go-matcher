package collection

import (
	"testing"
)

type Int = int
type MyString = string
type MyStruct struct {
	A int
	B string
}

func TestInSlice(t *testing.T) {
	structPtr := &MyStruct{A: 1, B: "a"}
	structPtr2 := &MyStruct{A: 2, B: "b"}
	structPtr3 := &MyStruct{A: 3, B: "c"}

	tests := []struct {
		name     string
		target   interface{}
		slice    interface{}
		expected bool
	}{
		// basic type
		{name: "FloatInSlice", target: 1.1, slice: []float64{1.1, 2.2, 3.3}, expected: true},
		{name: "FloatNotInSlice", target: 4.4, slice: []float64{1.1, 2.2, 3.3}, expected: false},
		// type alias
		{name: "IntInSlice", target: Int(1), slice: []Int{1, 2, 3}, expected: true},
		{name: "IntNotInSlice", target: Int(4), slice: []Int{1, 2, 3}, expected: false},
		{name: "StringInSlice", target: MyString("a"), slice: []MyString{"a", "b", "c"}, expected: true},
		{name: "StringNotInSlice", target: MyString("d"), slice: []MyString{"a", "b", "c"}, expected: false},
		// Pointer
		{
			name: "StructPtrInSlice", target: structPtr, slice: []*MyStruct{structPtr, structPtr2, structPtr3},
			expected: true,
		},
		{
			name: "StructPtrNotInSlice", target: structPtr3, slice: []*MyStruct{structPtr, structPtr2},
			expected: false,
		},
		// value
		{
			name: "StructInSlice", target: MyStruct{A: 1, B: "x"}, slice: []MyStruct{{A: 1, B: "x"}, {A: 2, B: "y"}},
			expected: true,
		},
		{
			name: "StructNotInSlice", target: MyStruct{A: 3, B: "z"}, slice: []MyStruct{{A: 1, B: "x"}, {A: 2, B: "y"}},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			switch v := test.slice.(type) {
			case []Int:
				result = InSlice(test.target.(Int), v)
			case []MyString:
				result = InSlice(test.target.(MyString), v)
			case []float64:
				result = InSlice(test.target.(float64), v)
			case []*Int:
				result = InSlice(test.target.(*Int), v)
			case []*MyStruct:
				result = InSlice(test.target.(*MyStruct), v)
			case []MyStruct:
				result = InSlice(test.target.(MyStruct), v)
			default:
				t.Fatalf("unsupported type: %T", v)
			}
			if result != test.expected {
				t.Errorf("InSlice(%v, %v) = %v; expected %v", test.target, test.slice, result, test.expected)
			}
		})
	}
}
