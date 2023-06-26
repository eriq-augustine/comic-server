package util;

import (
    "errors"

    "github.com/eriq-augustine/comic-server/config"
)

const MAX_CHAN_SIZE = 256
const MIN_CHAN_SIZE = 4

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
    fetchedResult bool
    closed bool
}

func NewPool[T any, R any](consumer PoolConsumer[T, R]) (*ParallelPool[T, R]) {
    return NewPoolWithSize(0, consumer);
}

func NewPoolWithSize[T any, R any](numValues int, consumer PoolConsumer[T, R]) (*ParallelPool[T, R]) {
    poolSize := config.GetInt("parallel.maxpoolsize");

    inputChan := make(chan poolInput[T], getChanSize(numValues));
    outputChan := make(chan poolResult[R], getChanSize(numValues));

    // Create the pool consumers.
    for i := 0; i < poolSize; i++ {
        go func() {
            for input := range inputChan {
                result, err := consumer(input.value);
                outputChan <- poolResult[R]{input.index, result, err};
            }
        }();
    }

    pool := ParallelPool[T, R]{
        workCount: 0,
        inputChan: inputChan,
        outputChan: outputChan,
        openForWork: true,
        fetchedResult: false,
        closed: false,
    };

    return &pool;
}

// Will panic if ParallelPool.NoMoreWork() is called first.
func (this *ParallelPool[T, R]) AddWork(value T) {
    if (!this.openForWork) {
        panic("Cannot add more work ParallelPool to after NoMoreWork() of GetAllResults called.");
    }

    this.inputChan <- poolInput[T]{this.workCount, value};
    this.workCount++;
}

func (this *ParallelPool[T, R]) NoMoreWork() {
    if (!this.openForWork) {
        return;
    }

    this.openForWork = false
    close(this.inputChan);
}

// Cannot be used with GetAllResults().
func (this *ParallelPool[T, R]) GetResult() (R, error) {
    this.fetchedResult = true;

    output := <- this.outputChan;
    return output.result, output.err;
}

// Will also call NoMoreWork() and Close().
// Cannot be used with GetResult().
func (this *ParallelPool[T, R]) GetAllResults() ([]R, error) {
    if (this.fetchedResult) {
        panic("Cannot use both GetResult() and GetAllResults() in ParallelPool.");
    }

    this.NoMoreWork();

    results := make([]R, this.workCount);
    var consumeErrors error = nil;

    // Consume output.
    for i := 0; i < this.workCount; i++ {
        output := <- this.outputChan;

        results[output.index] = output.result;
        consumeErrors = errors.Join(consumeErrors, output.err);
    }

    this.Close();

    return results, consumeErrors;
}

func (this *ParallelPool[T, R]) Close() {
    if (this.closed) {
        return;
    }

    this.closed = true;
    close(this.outputChan);
}

func PoolMap[T any, R any](values []T, consumer PoolConsumer[T, R]) ([]R, error) {
    pool := NewPoolWithSize(len(values), consumer);
    defer pool.Close();

    // Put the work in the queue.
    for _, value := range values {
        pool.AddWork(value);
    }
    pool.NoMoreWork();

    return pool.GetAllResults();
}

func getChanSize(numValues int) int {
    return Max(MIN_CHAN_SIZE, Min(MAX_CHAN_SIZE, numValues));
}
