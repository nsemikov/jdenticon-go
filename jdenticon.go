package jdenticon

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Jdenticon interface {
	SVG() ([]byte, error)
}

type jdenticon struct {
	config *Config
	svg    *SVG
	hash   string
	baHash [20]byte

	geometry Point
	paddings Point
	zero     Point
}

func New(identity string) Jdenticon {
	return NewWithConfig(identity, DefaultConfig)
}

func NewWithConfig(identity string, c *Config) Jdenticon {
	w := float64(c.Width)
	h := float64(c.Height)
	j := &jdenticon{
		config: c,
		svg: &SVG{
			Width:  c.Width,
			Height: c.Height,
		},
		hash:   sha1hash2string(sha1.Sum([]byte(identity))),
		baHash: sha1.Sum([]byte(identity)),
	}

	j.geometry = Point{
		X: float64(j.config.Width),
		Y: float64(j.config.Height),
	}
	j.paddings = Point{
		X: j.geometry.X * (c.Padding * 2),
		Y: j.geometry.Y * (c.Padding * 2),
	}
	if j.geometry.X > j.geometry.Y {
		j.paddings.X += j.geometry.X - j.geometry.Y
	}
	if j.geometry.Y > j.geometry.X {
		j.paddings.Y += j.geometry.Y - j.geometry.X
	}
	j.zero = Point{
		X: j.paddings.X / 2,
		Y: j.paddings.Y / 2,
	}

	if opacity(c.Background) != 0.0 {
		j.svg.Paths = append(j.svg.Paths, Path{
			Fill:       toHex(c.Background),
			UseOpacity: true,
			Opacity:    opacity(c.Background),
			Shapes: Shapes{
				&Polygon{[]Point{Point{0, 0}, Point{w, 0}, Point{w, h}, Point{0, h}}, false},
			},
		})
	}
	colors := j.colors()
	shapes := map[string]Shapes{}
	shapes[colors[0]] = Shapes{}
	shapes[colors[1]] = Shapes{}
	shapes[colors[2]] = Shapes{}
	shapes[colors[0]] = append(shapes[colors[0]], j.renderShapes(shapeOuter, 2, 3, [][2]float64{
		[2]float64{1, 0},
		[2]float64{2, 0},
		[2]float64{2, 3},
		[2]float64{1, 3},
		[2]float64{0, 1},
		[2]float64{3, 1},
		[2]float64{3, 2},
		[2]float64{0, 2},
	})...)
	shapes[colors[1]] = append(shapes[colors[1]], j.renderShapes(shapeOuter, 4, 5, [][2]float64{
		[2]float64{0, 0},
		[2]float64{3, 0},
		[2]float64{3, 3},
		[2]float64{0, 3},
	})...)
	shapes[colors[2]] = append(shapes[colors[2]], j.renderShapes(shapeInner, 1, 0, [][2]float64{
		[2]float64{1, 1},
		[2]float64{2, 1},
		[2]float64{2, 2},
		[2]float64{1, 2},
	})...)

	for color := range shapes {
		j.svg.Paths = append(j.svg.Paths, Path{
			Fill:   color,
			Shapes: shapes[color],
		})
	}
	return j
}

func (j *jdenticon) SVG() ([]byte, error) {
	t, err := template.New("svg").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	err = t.Execute(buf, j.svg)
	return buf.Bytes(), err
}

func sha1hash2string(hash [20]byte) string {
	h := make([]byte, len(hash))
	for i, v := range hash {
		h[i] = v
	}
	return hex.EncodeToString(h)
}
