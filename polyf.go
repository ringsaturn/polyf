// Generic point-in-polygon finder for large data sets.
package polyf

import (
	"errors"

	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
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

type RF[T any] struct {
	Tree *rtree.RTreeG[*Item[T]]
	F    *F[T]
}

func NewRFFromF[T any](f *F[T]) *RF[T] {
	tree := &rtree.RTreeG[*Item[T]]{}
	for _, item := range f.Items {
		minP := item.Poly.Rect().Min
		maxP := item.Poly.Rect().Max
		tree.Insert([2]float64{minP.X, minP.Y}, [2]float64{maxP.X, maxP.Y}, item)
	}
	return &RF[T]{
		Tree: tree,
		F:    f,
	}
}

func (rf *RF[T]) FindOne(x float64, y float64, xDiff float64, yDiff float64) (T, error) {
	p := geometry.Point{
		X: x,
		Y: y,
	}
	var res T
	hit := false
	rf.Tree.Search(
		[2]float64{x - xDiff, y - yDiff},
		[2]float64{x + xDiff, y + yDiff},
		func(min, max [2]float64, data *Item[T]) bool {
			if data.Poly.ContainsPoint(p) {
				res = data.V
				hit = true
				return false
			}
			return true
		},
	)
	if !hit {
		return *new(T), ErrNotFound
	}
	return res, nil
}

func (rf *RF[T]) FindAll(x float64, y float64, xDiff float64, yDiff float64) ([]T, error) {
	p := geometry.Point{
		X: x,
		Y: y,
	}
	res := make([]T, 0)
	rf.Tree.Search(
		[2]float64{x - xDiff, y - yDiff},
		[2]float64{x + xDiff, y + yDiff},
		func(min, max [2]float64, data *Item[T]) bool {
			if data.Poly.ContainsPoint(p) {
				res = append(res, data.V)
			}
			return true
		},
	)
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	return res, nil
}
