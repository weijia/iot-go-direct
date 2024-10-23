package bsp

import "testing"

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
