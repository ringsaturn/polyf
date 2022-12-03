package featurecollection

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/ringsaturn/polyf"
	"github.com/tidwall/geojson/geometry"
)

const (
	MultiPolygonType = "MultiPolygon"
	PolygonType      = "Polygon"
	FeatureType      = "Feature"
)

type PolygonCoordinates [][][2]float64
type MultiPolygonCoordinates []PolygonCoordinates

type Feature[T any] struct {
	Geometry struct {
		Coordinates interface{} `json:"coordinates"`
		Type        string      `json:"type"`
	} `json:"geometry"`
	Properties T      `json:"properties"`
	Type       string `json:"type"`
}

type BoundaryFile[T any] struct {
	Features []*Feature[T] `json:"features"`
}

func Do[T any](input *BoundaryFile[T]) (*polyf.F[T], error) {
	f := &polyf.F[T]{}

	for _, item := range input.Features {
		var coordinates MultiPolygonCoordinates

		MultiPolygonTypeHandler := func() error {
			if err := mapstructure.Decode(item.Geometry.Coordinates, &coordinates); err != nil {
				return err
			}
			return nil
		}
		PolygonTypeHandler := func() error {
			var polygonCoordinates PolygonCoordinates
			if err := mapstructure.Decode(item.Geometry.Coordinates, &polygonCoordinates); err != nil {
				return err
			}
			coordinates = append(coordinates, polygonCoordinates)
			return nil
		}

		switch item.Type {
		case MultiPolygonType:
			if err := MultiPolygonTypeHandler(); err != nil {
				return nil, err
			}
		case PolygonType:
			if err := PolygonTypeHandler(); err != nil {
				return nil, err
			}
		case FeatureType:
			switch item.Geometry.Type {
			case MultiPolygonType:
				if err := MultiPolygonTypeHandler(); err != nil {
					return nil, err
				}
			case PolygonType:
				if err := PolygonTypeHandler(); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unknown type %v", item.Type)
			}
		default:
			return nil, fmt.Errorf("unknown type %v", item.Type)
		}

		for _, subcoordinates := range coordinates {
			exterior := []geometry.Point{}
			holes := [][]geometry.Point{}
			for index, geopoly := range subcoordinates {
				if index == 0 {
					for _, rawCoods := range geopoly {
						exterior = append(exterior, geometry.Point{X: rawCoods[0], Y: rawCoods[1]})
					}
					continue
				}
				holepoints := []geometry.Point{}
				for _, rawCoods := range geopoly {
					holepoints = append(holepoints, geometry.Point{X: rawCoods[0], Y: rawCoods[1]})
				}
				holes = append(holes, holepoints)
			}
			newItem := &polyf.Item[T]{
				V:    item.Properties,
				Poly: geometry.NewPoly(exterior, holes, nil),
			}

			f.Items = append(f.Items, newItem)
		}
	}
	return f, nil
}
