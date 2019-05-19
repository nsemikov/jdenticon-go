package jdenticon

import (
	"encoding/xml"
	"fmt"
	"math"
	"strings"
)

// See https://developer.mozilla.org/ru/docs/Web/SVG/Tutorial/Paths

// Want render something like this:
//
// <svg height="200" preserveAspectRatio="xMidYMid meet" viewBox="0 0 200 200" width="200" xmlns="http://www.w3.org/2000/svg">
//   <path fill="#3d73b7" d="M100 58L58 58L58 37ZM100 58L100 16L121 16ZM100 142L142 142L142 163ZM100 142L100 184L79 184ZM58 100L16 100L16 79ZM142 100L142 58L163 58ZM142 100L184 100L184 121ZM58 100L58 142L37 142Z"/>
//   <path fill="#84a8d6" d="M58 16L58 58L16 58ZM184 58L142 58L142 16ZM142 184L142 142L184 142ZM16 142L58 142L58 184ZM73 83.5a10.5,10.5 0 1,1 21,0a10.5,10.5 0 1,1 -21,0M106 83.5a10.5,10.5 0 1,1 21,0a10.5,10.5 0 1,1 -21,0M106 116.5a10.5,10.5 0 1,1 21,0a10.5,10.5 0 1,1 -21,0M73 116.5a10.5,10.5 0 1,1 21,0a10.5,10.5 0 1,1 -21,0"/>
// </svg>
//

// M73 83.5a10.5,10.5 0 1,1 21,0a10.5,10.5 0 1,1 -21,0

// M73 83.5
// a 10.5 10.5, 0 1 1, 21 0
// a 10.5 10.5, 0 1 1, -21 0

type SVG struct {
	XMLName             xml.Name `xml:"svg"`
	Width               int      `xml:"width,attr"`
	Height              int      `xml:"height,attr"`
	PreserveAspectRatio string   `xml:"preserveAspectRatio,attr"`
	ViewBox             string   `xml:"viewBox,attr"`
	Namespace           string   `xml:"xmlns,attr"`
	Paths               `xml:"path"`
}

// -----------------------------------------------------------------------------

type Paths []Path

// -----------------------------------------------------------------------------

type Path struct {
	XMLName    xml.Name `xml:"path"`
	Shapes     Shapes   `xml:"d,attr"`
	Fill       string   `xml:"fill,attr,omitempty"`
	Opacity    float64  `xml:"opacity,attr,omitempty"`
	Stroke     string   `xml:"stroke,attr,omitempty"`
	UseOpacity bool     `xml:"-"`
}

func (p Path) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if p.UseOpacity && p.Opacity == 0 {
		return nil
	}
	if len(p.Fill) > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "fill"}, Value: p.Fill})
	}
	if len(p.Stroke) > 0 {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "stroke"}, Value: p.Stroke})
	}
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "d"}, Value: p.Shapes.String()})
	if p.UseOpacity {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "opacity"}, Value: fmt.Sprintf("%f", p.Opacity)})
	}
	start.End()
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	// flush to ensure tokens are written
	return e.Flush()
}

// -----------------------------------------------------------------------------

type Shapes []Shape

func (shapes Shapes) String() string {
	paths := []string{}
	for _, shape := range shapes {
		paths = append(paths, shape.Path())
	}
	return strings.Join(paths, "")
}

// -----------------------------------------------------------------------------

type Shape interface {
	Path() string
	Rotate(deg float64, center *Point)
	Translate(dx, dy float64)
	Copy() Shape
}

// -----------------------------------------------------------------------------

type Point struct {
	X float64
	Y float64
}

func (p *Point) Path() string {
	return fmt.Sprintf("%.f,%.f", p.X, p.Y)
}

func (p *Point) Translate(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// -----------------------------------------------------------------------------

type Circle struct {
	Center    Point
	Radius    float64
	Clockwise bool
}

func (c *Circle) Path() string {
	var path string
	arc1 := fmt.Sprintf("a%.1f,%.1f 0 1,1 %.1f,0", c.Radius, c.Radius, c.Radius*2)
	arc2 := fmt.Sprintf("a%.1f,%.1f 0 1,1 -%.1f,0", c.Radius, c.Radius, c.Radius*2)
	if !c.Clockwise {
		position := fmt.Sprintf("M%.f,%.f", c.Center.X-c.Radius, c.Center.Y)
		path = position + arc1 + arc2
	} else {
		position := fmt.Sprintf("M%.f,%.f", c.Center.X+c.Radius, c.Center.Y)
		path = position + arc2 + arc1
	}
	return path
}

func (c *Circle) Translate(dx, dy float64) {
	c.Center.Translate(dx, dy)
}

func (c *Circle) Rotate(deg float64, center *Point) {
	x0 := 0.0
	y0 := 0.0
	if center != nil {
		x0 = center.X
		y0 = center.Y
	}
	c.Center.X, c.Center.Y = rotate(deg, c.Center.X, c.Center.Y, x0, y0)
}

func (c *Circle) Copy() Shape {
	result := &Circle{}
	*result = *c
	return result
}

func newCircle(x, y, size float64, cw bool) *Circle {
	return &Circle{
		Center: Point{
			X: x + size/2,
			Y: y + size/2,
		},
		Radius:    size / 2,
		Clockwise: cw,
	}
}

// -----------------------------------------------------------------------------

type Polygon struct {
	Points    []Point
	Clockwise bool
}

func (s *Polygon) Path() string {
	if len(s.Points) == 0 {
		return ""
	}
	points := []string{}
	if s.Clockwise {
		for idx := 0; idx < len(s.Points); idx++ {
			points = append(points, s.Points[idx].Path())
		}
	} else {
		for idx := len(s.Points) - 1; idx >= 0; idx-- {
			points = append(points, s.Points[idx].Path())
		}
	}
	return "M" + strings.Join(points, "L") + "Z"
}

func (s *Polygon) Translate(dx, dy float64) {
	for idx := range s.Points {
		s.Points[idx].Translate(dx, dy)
	}
}

func (s *Polygon) Rotate(deg float64, center *Point) {
	x0 := 0.0
	y0 := 0.0
	if center != nil {
		x0 = center.X
		y0 = center.Y
	}
	for idx := range s.Points {
		p := &(s.Points[idx])
		p.X, p.Y = rotate(deg, p.X, p.Y, x0, y0)
	}
}

func (s *Polygon) Copy() Shape {
	result := &Polygon{}
	*result = *s
	return result
}

func newPolygon(points []Point, cw bool) *Polygon {
	return &Polygon{
		points,
		cw,
	}
}

// -----------------------------------------------------------------------------

func newTriangle(x, y, w, h float64, r int, cw bool) Shape {
	points := []Point{
		Point{x + w, y},
		Point{x + w, y + h},
		Point{x, y + h},
		Point{x, y},
	}
	idx := r % 4
	if len(points) == idx {
		points = points[:idx]
	} else {
		points = append(points[:idx], points[idx+1:]...)
	}
	return newPolygon(points, cw)
}

// -----------------------------------------------------------------------------

func newRectangle(x, y, w, h float64, cw bool) Shape {
	return newPolygon([]Point{
		Point{x, y},
		Point{x + w, y},
		Point{x + w, y + h},
		Point{x, y + h},
	}, cw)
}

// -----------------------------------------------------------------------------

func newRhombus(x, y, w, h float64, cw bool) Shape {
	return newPolygon([]Point{
		Point{x + w/2, y},
		Point{x + w, y + h/2},
		Point{x + w/2, y + h},
		Point{x, y + h/2},
	}, cw)
}

// -----------------------------------------------------------------------------

func rotate(deg float64, x, y, x0, y0 float64) (float64, float64) {
	rad := deg * math.Pi / 180
	xn := x0 + (x-x0)*math.Cos(rad) - (y-y0)*math.Sin(rad)
	yn := y0 + (y-y0)*math.Cos(rad) + (x-x0)*math.Sin(rad)
	return xn, yn
}
