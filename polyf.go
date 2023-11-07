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
	RTree *rtree.RTreeG[*Item[T]] // RTree
}

func (f *F[T]) SetupRTreeIndex() {
	if f.RTree != nil {
		return
	}
	tr := &rtree.RTreeG[*Item[T]]{}
	for _, item := range f.Items {
		minP := item.Poly.Exterior.Rect().Min
		maxP := item.Poly.Exterior.Rect().Max
		tr.Insert([2]float64{minP.X, minP.Y}, [2]float64{maxP.X, maxP.Y}, item)
	}
	f.RTree = tr
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

func (f *F[T]) FindAllWithRTRee(x float64, y float64) ([]T, error) {

}
