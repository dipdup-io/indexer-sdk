package modules

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testModule struct {
	Field1 int
	Field2 string
	Field3 chan int

	Output *Output[int]
}

type testDirectionalModule struct {
	Field1 int
	Field2 string
	Field3 chan int
	Field4 <-chan bool

	Output *Output[int]

	Struct    struct{}
	StructPtr *struct{}

	Int *int

	Array  []int
	Array2 []<-chan int
}

func Test_getModulePorts(t *testing.T) {
	tests := []struct {
		name    string
		module  any
		want    modulePorts
		wantErr bool
	}{
		{
			name:   "test 1",
			module: &testModule{},
			want: modulePorts{
				inputs: ports{
					"Field3": reflect.Value{},
				},
				outputs: ports{
					"Output": reflect.Value{},
				},
			},
			wantErr: false,
		}, {
			name:   "test 2",
			module: &testDirectionalModule{},
			want: modulePorts{
				inputs: ports{
					"Field3": reflect.Value{},
					"Field4": reflect.Value{},
				},
				outputs: ports{
					"Output": reflect.Value{},
				},
			},
			wantErr: false,
		}, {
			name:    "test 3",
			module:  1000,
			want:    newModulePorts(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getModulePorts(tt.module)
			require.Equal(t, (err != nil), tt.wantErr, "getModulePorts")
			for key := range tt.want.inputs {
				require.Contains(t, got.inputs, key)
			}
			for key := range got.inputs {
				require.Contains(t, tt.want.inputs, key)
			}
			for key := range tt.want.outputs {
				require.Contains(t, got.outputs, key)
			}
			for key := range got.outputs {
				require.Contains(t, tt.want.outputs, key)
			}
		})
	}
}
