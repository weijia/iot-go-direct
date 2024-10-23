package bsp

import (
	"reflect"
	"testing"
)

func TestInitRequestingColors(t *testing.T) {
	tests := []struct {
		name      string
		nodeStates []NodeState
		want      []NodeState
	}{
		{
			name: "测试用例1",
			nodeStates: []NodeState{
				{NodeRequestingColor: [8]int{1, 2, 3, 4, 5, 6, 7, 8}},
				{NodeRequestingColor: [8]int{9, 10, 11, 12, 13, 14, 15, 16}},
			},
			want: []NodeState{
				{NodeRequestingColor: [8]int{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}},
				{NodeRequestingColor: [8]int{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}},
			},
		},
		{
			name: "测试用例2",
			nodeStates: []NodeState{
				{NodeRequestingColor: [8]int{0, 0, 0, 0, 0, 0, 0, 0}},
			},
			want: []NodeState{
				{NodeRequestingColor: [8]int{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}},
			},
		},
		{
			name: "测试用例3",
			nodeStates: []NodeState{
				{NodeRequestingColor: [8]int{100, 200, 300, 400, 500, 600, 700, 800}},
				{NodeRequestingColor: [8]int{-1, -2, -3, -4, -5, -6, -7, -8}},
			},
			want: []NodeState{
				{NodeRequestingColor: [8]int{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}},
				{NodeRequestingColor: [8]int{0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF, 0xF}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BspConfigInstance.NodeStates = tt.nodeStates
			initRequestingColors()
			if !reflect.DeepEqual(BspConfigInstance.NodeStates, tt.want) {
				t.Errorf("initRequestingColors() = %v, want %v", BspConfigInstance.NodeStates, tt.want)
			}
		})
	}
}
