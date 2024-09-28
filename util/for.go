package util

import html "github.com/Gardego5/htmdsl"

func For[T any](iterable []T, f func(int, T) any) html.Fragment {
	res := make(html.Fragment, len(iterable))
	for i, v := range iterable {
		res[i] = f(i, v)
	}
	return res
}
