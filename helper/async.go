package helper

import (
	"sort"
	"sync"
)

var wg sync.WaitGroup

type AsyncResult[T any] struct {
	i   int
	Res T
	Err error
}

func Async(funcs ...func() (any, error)) any {
	resCh := make(chan AsyncResult[any], len(funcs))

	for i, f := range funcs {
		wg.Add(1)

		go func(f func() (any, error), i int) {
			defer wg.Done()

			res, err := f()
			resCh <- AsyncResult[any]{
				i:   i,
				Res: res,
				Err: err,
			}
		}(f, i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	results := []AsyncResult[any]{}

	for res := range resCh {
		results = append(results, res)
	}

	sort.Slice(results, func(a, b int) bool {
		return results[a].i < results[b].i
	})

	return results
}
