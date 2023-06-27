package util;

import (
    "errors"

    "github.com/eriq-augustine/comic-server/config"
)

const MIN_CHAN_SIZE = 256;

type PoolConsumer[T any, R any] func(value T) (R, error);

type poolInput[T any] struct {
    index int
    value T
}

type poolResult[R any] struct {
    index int
    result R
    err error
}

// A pool should only be accessed from a single thread.
type ParallelPool[T any, R any] struct {
    workCount int
    inputChan chan poolInput[T]
    outputChan chan poolResult[R]
    openForWork bool
    closed bool
}

func NewParallelPool[T any, R any](consumer PoolConsumer[T, R]) (*ParallelPool[T, R]) {
    return NewParallelPoolWithSize(0, consumer);
}

func NewParallelPoolWithSize[T any, R any](numValues int, consumer PoolConsumer[T, R]) (*ParallelPool[T, R]) {
    poolSize := config.GetInt("parallel.maxpoolsize");

    inputChan := make(chan poolInput[T], getChanSize(numValues));
    outputChan := make(chan poolResult[R], getChanSize(numValues));

    exitedConsumers := make(chan int, poolSize);

    // Create the pool consumers.
    for i := 0; i < poolSize; i++ {
        go func(id int) {
            for input := range inputChan {
                result, err := consumer(input.value);
                outputChan <- poolResult[R]{input.index, result, err};
            }

            exitedConsumers <- id;
        }(i);
    }

    // Close the output channel once all the consumers have stopped.
    go func() {
        for i := 0; i < poolSize; i++ {
            <- exitedConsumers;
        }

        close(outputChan);
    }();

    pool := ParallelPool[T, R]{
        workCount: 0,
        inputChan: inputChan,
        outputChan: outputChan,
        openForWork: true,
        closed: false,
    };

    return &pool;
}

// Will panic if ParallelPool.NoMoreWork() is called first.
// Should be called on the same thread as NoMoreWork.
func (this *ParallelPool[T, R]) AddWork(value T) {
    if (!this.openForWork) {
        panic("Cannot add more work to a ParallelPool after NoMoreWork() is called.");
    }

    this.inputChan <- poolInput[T]{this.workCount, value};
    this.workCount++;
}

// Should be called on the same thread as AddWork.
func (this *ParallelPool[T, R]) NoMoreWork() {
    if (!this.openForWork) {
        return;
    }

    this.openForWork = false
    close(this.inputChan);
}

// The returned boolean will be false when the output channel has been closed and there is no work
// work on it (i.e., all work is consumed).
func (this *ParallelPool[T, R]) GetResult() (R, bool, error) {
    output, ok := <- this.outputChan;
    if (!ok) {
        var empty R;
        return empty, false, nil;
    }

    return output.result, true, output.err;
}

func (this *ParallelPool[T, R]) Close() {
    if (this.closed) {
        return;
    }

    this.NoMoreWork();
    this.closed = true;
}

func PoolMap[T any, R any](values []T, consumer PoolConsumer[T, R]) ([]R, error) {
    pool := NewParallelPoolWithSize(len(values), consumer);
    defer pool.Close();

    // Put the work in the queue.
    for _, value := range values {
        pool.AddWork(value);
    }
    pool.NoMoreWork();

    results := make([]R, len(values));
    var allErrors error = nil;

    for i := 0; i < len(values); i++ {
        result := <- pool.outputChan;
        results[result.index] = result.result;
        allErrors = errors.Join(allErrors, result.err);
    }

    return results, allErrors;
}

func getChanSize(numValues int) int {
    return Max(MIN_CHAN_SIZE, numValues);
}
