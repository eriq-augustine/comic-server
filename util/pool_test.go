package util;

import (
    "reflect"
    "testing"
)

type poolMapTestCase[T any, R any] struct {
    in []T
    out []R
    consumer PoolConsumer[T, R]
}

var poolMapCases = []poolMapTestCase[int, int]{
    poolMapTestCase[int, int]{[]int{1, 2, 3}, []int{-1, -2, -3}, negateConsumer},
    poolMapTestCase[int, int]{[]int{}, []int{}, negateConsumer},
    poolMapTestCase[int, int]{[]int{0, 1, 2, 3}, []int{0, -1, -2, -3}, negateConsumer},
};

func TestPoolMapTable(test *testing.T) {
    for i, testCase := range poolMapCases {
        result, err := PoolMap(testCase.in, testCase.consumer);
        if (err != nil) {
            test.Errorf("Error in pool map test case (index %d) of %v: %v.", i, testCase.in, err);
            continue;
        }

        if (!reflect.DeepEqual(result, testCase.out)) {
            test.Errorf("Failed pool map test case (index %d) of %v. Expected %v, got %v.", i, testCase.in, testCase.out, result);
            continue;
        }
    }
}

func negateConsumer(value int) (int, error) {
    return -value, nil;
}
