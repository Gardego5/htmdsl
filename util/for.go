package util

import html "github.com/Gardego5/htmdsl"

func For[T any](iterable []T, f func(int, T) any) html.Fragment {
	res := make(html.Fragment, len(iterable))
	for i, v := range iterable {
		res[i] = f(i, v)
	}
	return res
}

func ForMap[K comparable, V any](m map[K]V, f func(K, V) any) html.Fragment {
	res := make(html.Fragment, len(m))
	i := 0
	for k, v := range m {
		res[i] = f(k, v)
		i++
	}
	return res
}
