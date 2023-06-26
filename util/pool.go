package util;

import (
    "errors"

    "github.com/eriq-augustine/comic-server/config"
)

const MAX_CHAN_SIZE = 256

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

func PoolMap[T any, R any](values []T, consumer PoolConsumer[T, R]) ([]R, error) {
    poolSize := config.GetInt("parallel.maxpoolsize");

    inputChan := make(chan poolInput[T], Min(len(values), MAX_CHAN_SIZE));
    outputChan := make(chan poolResult[R], Min(len(values), MAX_CHAN_SIZE));

    // Create the pool consumers.
    for i := 0; i < poolSize; i++ {
        go func() {
            for input := range inputChan {
                result, err := consumer(input.value);
                outputChan <- poolResult[R]{input.index, result, err};
            }
        }();
    }

    // Put the work in the queue.
    for index, value := range values {
        inputChan <- poolInput[T]{index, value};
    }
    close(inputChan);

    results := make([]R, len(values));
    var consumeErrors error = nil;

    // Consume output.
    for i := 0; i < len(values); i++ {
        output := <- outputChan;

        results[output.index] = output.result;
        consumeErrors = errors.Join(consumeErrors, output.err);
    }
    close(outputChan);

    return results, consumeErrors;
}
