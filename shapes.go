package jdenticon

import (
	"strconv"
)

type shapesGetter func(cell float64, index int) Shapes

// nolint:gochecknoglobals
var shapeInner = []shapesGetter{
	func(cell float64, index int) Shapes {
		k := cell * 0.42
		return Shapes{newPolygon([]Point{
			{0, 0},
			{cell, 0},
			{cell, cell - k*2},
			{cell - k, cell},
			{0, cell},
		}, false)}
	},
	func(cell float64, index int) Shapes {
		w := cell * 0.5
		h := cell * 0.8
		return Shapes{newTriangle(cell-w, 0, w, h, 2, false)}
	},
	func(cell float64, index int) Shapes {
		s := cell / 3
		return Shapes{newRectangle(s, s, cell-s, cell-s, false)}
	},
	func(cell float64, index int) Shapes {
		inner := cell * 0.1
		// Use fixed outer border widths in small icons to ensure the border is drawn
		outer := 1.0
		if cell >= 6 {
			outer = 2
			if cell >= 8 {
				outer = cell * 0.25
			}
		}
		if inner > 1 { // large icon => truncate decimals
			inner = float64(int(inner))
		} else if inner > 0.5 && inner <= 1 { // medium size icon => fixed width
			inner = 1
		} else if inner <= 0.5 { // small icon => anti-aliased border
			inner = 0
		}
		return Shapes{newRectangle(outer, outer, cell-inner-outer, cell-inner-outer, false)}
	},
	func(cell float64, index int) Shapes {
		m := cell * 0.15
		s := cell * 0.5
		return Shapes{newCircle(cell-s-m, cell-s-m, s, false)}
	},
	func(cell float64, index int) Shapes {
		inner := cell * 0.1
		outer := inner * 4

		// Align edge to nearest pixel in large icons
		if outer > 3 {
			outer = float64(int(outer))
		}
		return Shapes{
			newRectangle(0, 0, cell, cell, true),
			newPolygon([]Point{
				{outer, outer},
				{cell - inner, outer},
				{outer + (cell-outer-inner)/2, cell - inner},
			}, true),
		}
	},
	func(cell float64, index int) Shapes {
		return Shapes{newPolygon([]Point{
			{0, 0},
			{cell, 0},
			{cell, cell * 0.7},
			{cell * 0.4, cell * 0.4},
			{cell * 0.7, cell},
			{0, cell},
		}, false)}
	},
	func(cell float64, index int) Shapes {
		return Shapes{newTriangle(cell/2, cell/2, cell/2, cell/2, 3, false)}
	},
	func(cell float64, index int) Shapes {
		return Shapes{
			newRectangle(0, 0, cell, cell/2, false),
			newRectangle(0, cell/2, cell/2, cell/2, false),
			newTriangle(cell/2, cell/2, cell/2, cell/2, 1, false),
		}
	},
	func(cell float64, index int) Shapes {
		inner := cell * 0.14
		// Use fixed outer border widths in small icons to ensure the border is drawn
		outer := 0.0
		if cell < 4 {
			outer = 1
		} else if cell < 6 {
			outer = 2
		} else {
			outer = float64(int(cell * 0.35))
		}
		if cell >= 8 { // large icon => truncate decimals
			inner = float64(int(inner))
		}
		return Shapes{
			newRectangle(0, 0, cell, cell, false),
			newRectangle(outer, outer, cell-outer-inner, cell-outer-inner, true),
		}
	},
	func(cell float64, index int) Shapes {
		inner := cell * 0.12
		outer := inner * 3
		return Shapes{
			newRectangle(0, 0, cell, cell, false),
			newCircle(outer, outer, cell-inner-outer, true),
		}
	},
	func(cell float64, index int) Shapes {
		return Shapes{
			newTriangle(cell/2, cell/2, cell/2, cell/2, 3, false),
		}
	},
	func(cell float64, index int) Shapes {
		m := cell * 0.25
		return Shapes{
			newRectangle(0, 0, cell, cell, false),
			newRhombus(m, m, cell-m, cell-m, true),
		}
	},
	func(cell float64, index int) Shapes {
		m := cell * 0.4
		s := cell * 1.2
		var shapes Shapes
		if index == 0 {
			shapes = Shapes{newCircle(m, m, s, false)}
		}
		return shapes
	},
}

// nolint:gochecknoglobals
var shapeOuter = []shapesGetter{
	func(cell float64, index int) Shapes {
		return Shapes{newTriangle(0, 0, cell, cell, 0, false)}
	},
	func(cell float64, index int) Shapes {
		return Shapes{newTriangle(0, cell/2, cell, cell/2, 0, false)}
	},
	func(cell float64, index int) Shapes {
		return Shapes{newRhombus(0, 0, cell, cell, false)}
	},
	func(cell float64, index int) Shapes {
		m := cell / 6
		return Shapes{newCircle(m, m, cell-2*m, false)}
	},
}

func (j *jdenticon) renderShapes(getters []shapesGetter, index int, rotationIndex int, positions [][2]float64) Shapes {
	r := 0
	if rotationIndex > 0 {
		h, _ := strconv.ParseInt("0x"+j.hash[rotationIndex:rotationIndex+1], 0, 64)
		r = int(h)
	}
	shapeIdx, _ := strconv.ParseInt("0x"+j.hash[index:index+1], 0, 64)
	getter := getters[int(shapeIdx)%len(getters)]
	width := j.geometry.X - j.paddings.X
	// height := j.geometry.Y - j.paddings.Y
	cell := width / 4
	result := Shapes{}
	for i := range positions {
		shapes := getter(cell, index)
		for _, shape := range shapes {
			bottomleft := &Point{
				X: j.zero.X + positions[i][0]*cell,
				Y: j.zero.Y + positions[i][1]*cell,
			}
			center := &Point{
				X: bottomleft.X + cell/2,
				Y: bottomleft.Y + cell/2,
			}
			shape.Translate(bottomleft.X, bottomleft.Y)
			shape.Rotate(float64(r%4)*90, center)
			result = append(result, shape)
		}
		r++
	}
	return result
}
