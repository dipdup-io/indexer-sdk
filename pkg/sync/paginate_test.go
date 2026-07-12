package sync

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// fetcher returns pages from items and records every (limit, offset) call.
type fetcher struct {
	items []int
	calls [][2]int
	errAt int // offset at which to return err (-1 to disable)
	err   error
}

func (f *fetcher) fetch(_ context.Context, limit, offset int) ([]int, error) {
	f.calls = append(f.calls, [2]int{limit, offset})
	if f.err != nil && offset == f.errAt {
		return nil, f.err
	}
	if offset >= len(f.items) {
		return []int{}, nil
	}
	return f.items[offset:min(offset+limit, len(f.items))], nil
}

func makeItems(count int) []int {
	items := make([]int, count)
	for i := range items {
		items[i] = i
	}
	return items
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		count     int
		wantCalls [][2]int
	}{
		{
			name:      "partial last page",
			limit:     10,
			count:     25,
			wantCalls: [][2]int{{10, 0}, {10, 10}, {10, 20}},
		}, {
			name:      "count is a multiple of limit",
			limit:     10,
			count:     20,
			wantCalls: [][2]int{{10, 0}, {10, 10}, {10, 20}},
		}, {
			name:      "single short page",
			limit:     10,
			count:     3,
			wantCalls: [][2]int{{10, 0}},
		}, {
			name:      "empty first page",
			limit:     10,
			count:     0,
			wantCalls: [][2]int{{10, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fetcher{items: makeItems(tt.count), errAt: -1}

			got := make([]int, 0)
			for item, err := range Paginate(t.Context(), tt.limit, f.fetch) {
				require.NoError(t, err)
				got = append(got, item)
			}

			require.Equal(t, makeItems(tt.count), got)
			require.Equal(t, tt.wantCalls, f.calls)
		})
	}
}

func TestPaginateErrorOnFirstPage(t *testing.T) {
	testErr := errors.New("test error")
	f := fetcher{items: makeItems(25), errAt: 0, err: testErr}

	var count int
	var gotErr error
	for item, err := range Paginate(t.Context(), 10, f.fetch) {
		if err != nil {
			gotErr = err
			require.Zero(t, item, "error must be yielded with zero value")
			break
		}
		count++
	}

	require.ErrorIs(t, gotErr, testErr)
	require.Zero(t, count, "no items expected before the error")
	require.Len(t, f.calls, 1)
}

func TestPaginateErrorOnSecondPage(t *testing.T) {
	testErr := errors.New("test error")
	f := fetcher{items: makeItems(25), errAt: 10, err: testErr}

	var got []int
	var gotErr error
	for item, err := range Paginate(t.Context(), 10, f.fetch) {
		if err != nil {
			gotErr = err
			break
		}
		got = append(got, item)
	}

	require.ErrorIs(t, gotErr, testErr)
	require.Equal(t, makeItems(10), got, "first page must be yielded before the error")
	require.Len(t, f.calls, 2)
}

func TestPaginateEarlyBreak(t *testing.T) {
	f := fetcher{items: makeItems(25), errAt: -1}

	var got []int
	for item, err := range Paginate(t.Context(), 10, f.fetch) {
		require.NoError(t, err)
		got = append(got, item)
		if len(got) == 5 {
			break
		}
	}

	require.Equal(t, makeItems(5), got)
	require.Len(t, f.calls, 1, "no fetch calls expected after break")
}

func TestPaginateNonPositiveLimit(t *testing.T) {
	for _, limit := range []int{0, -1} {
		f := fetcher{items: makeItems(25), errAt: -1}

		for item, err := range Paginate(t.Context(), limit, f.fetch) {
			require.NoError(t, err)
			require.Fail(t, "unexpected item", "item: %d", item)
		}

		require.Empty(t, f.calls, "fetch must not be called with limit %d", limit)
	}
}
