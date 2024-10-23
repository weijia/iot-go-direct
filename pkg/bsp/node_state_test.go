package bsp

import (
	"reflect"
	"testing"
)

func TestCompareArrays(t *testing.T) {
    tests := []struct {
        name string
        a    [8]int
        b    [8]int
        want bool
    }{
        {
            name: "两个数组完全相同",
            a:    [8]int{1, 2, 3, 4, 5, 6, 7, 8},
            b:    [8]int{1, 2, 3, 4, 5, 6, 7, 8},
            want: true,
        },
        {
            name: "两个数组有一个元素不同",
            a:    [8]int{1, 2, 3, 4, 5, 6, 7, 8},
            b:    [8]int{1, 2, 3, 4, 5, 6, 7, 9},
            want: false,
        },
        {
            name: "两个数组有一个元素为0xF",
            a:    [8]int{1, 2, 3, 4, 5, 6, 7, 0xF},
            b:    [8]int{1, 2, 3, 4, 5, 6, 7, 8},
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := IsEqual(tt.a, tt.b); got != tt.want {
                t.Errorf("CompareArrays() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestSetRequestingColorByNodeStateRef(t *testing.T) {
	tests := []struct {
		name             string
		nodeState        *NodeState
		requestingColor  string
		expectedNodeState *NodeState
	}{
		{
			name:             "测试用例1：正常的十六进制字符串",
			nodeState:        &NodeState{},
			requestingColor:  "1a2b3c4d",
			expectedNodeState: &NodeState{NodeRequestingColor: [8]int{10, 11, 12, 13, 0, 0, 0, 0}},
		},
		{
			name:             "测试用例2：空的十六进制字符串",
			nodeState:        &NodeState{},
			requestingColor:  "",
			expectedNodeState: &NodeState{NodeRequestingColor: [8]int{0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			name:             "测试用例3：非法的十六进制字符串",
			nodeState:        &NodeState{},
			requestingColor:  "zzzz",
			expectedNodeState: &NodeState{NodeRequestingColor: [8]int{0, 0, 0, 0, 0, 0, 0, 0}},
		},
		{
			name:             "测试用例4：部分有效的十六进制字符串",
			nodeState:        &NodeState{},
			requestingColor:  "5e6f",
			expectedNodeState: &NodeState{NodeRequestingColor: [8]int{0, 0, 0, 0, 14, 15, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetRequestingColorByNodeStateRef(tt.nodeState, tt.requestingColor)
			if !reflect.DeepEqual(tt.nodeState, tt.expectedNodeState) {
				t.Errorf("SetRequestingColorByNodeStateRef() = %v, want %v", tt.nodeState, tt.expectedNodeState)
			}
		})
	}
}
