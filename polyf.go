package polyf

import "github.com/tidwall/geojson/geometry"

type item[T any] struct {
	v    T
	poly *geometry.Poly
}

type F[T any] struct {
	items []*item[T]
}

func (f *F[T]) Insert(poly *geometry.Poly, v T) {
	f.items = append(f.items, &item[T]{
		v:    v,
		poly: poly,
	})
}

func (f *F[T]) Contains(x float64, y float64) []T {
	res := make([]T, 0)
	p := geometry.Point{
		X: x,
		Y: y,
	}
	for _, item := range f.items {
		if item.poly.ContainsPoint(p) {
			res = append(res, item.v)
		}
	}
	return res
}
