// Generic point-in-polygon finder for large data sets.
package polyf

import (
	"errors"

	"github.com/tidwall/geojson/geometry"
)

var ErrNotFound = errors.New("polyf: not found")

type Item[T any] struct {
	V    T
	Poly *geometry.Poly
}

type F[T any] struct {
	Items []*Item[T]
}

func (f *F[T]) Insert(poly *geometry.Poly, v T) {
	f.Items = append(f.Items, &Item[T]{
		V:    v,
		Poly: poly,
	})
}

func (f *F[T]) FindOne(x float64, y float64) (T, error) {
	p := geometry.Point{
		X: x,
		Y: y,
	}
	for _, item := range f.Items {
		if item.Poly.ContainsPoint(p) {
			return item.V, nil
		}
	}
	return *new(T), ErrNotFound
}

func (f *F[T]) FindAll(x float64, y float64) ([]T, error) {
	res := make([]T, 0)
	p := geometry.Point{
		X: x,
		Y: y,
	}
	for _, item := range f.Items {
		if item.Poly.ContainsPoint(p) {
			res = append(res, item.V)
		}
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	return res, nil
}
