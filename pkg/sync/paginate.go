package sync

import (
	"context"
	"iter"
)

// Paginate returns an iterator over items fetched page by page with limit/offset
// pagination. It calls fetch with increasing offset until a page shorter than
// limit is returned, so when the total count is a multiple of limit the final
// fetch returns an empty page. If limit is not positive, the sequence is empty
// and fetch is never called.
//
// On fetch failure the error is yielded as the last pair (with zero T) and the
// sequence ends, so callers must check the error on every iteration. Breaking
// out of the loop stops the iterator without fetching further pages.
func Paginate[T any](ctx context.Context, limit int,
	fetch func(ctx context.Context, limit, offset int) ([]T, error),
) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if limit <= 0 {
			return
		}
		var offset int
		for {
			res, err := fetch(ctx, limit, offset)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}

			for i := range res {
				if !yield(res[i], nil) {
					return
				}
			}

			if len(res) < limit {
				return
			}
			offset += len(res)
		}
	}
}
