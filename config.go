package jdenticon

import (
	"fmt"
	"image/color"
	"strconv"
)

// DefaultConfig var
var DefaultConfig = &Config{ // nolint:gochecknoglobals
	Hues: -1,
	Colored: Color{
		Saturation: 0.5,
		Lightness:  []float64{0.4, 0.8},
	},
	Grayscale: Color{
		Saturation: 0.0,
		Lightness:  []float64{0.3, 0.9},
	},
	Background: color.RGBA{0xff, 0xff, 0xff, 0x00},
	Width:      200,
	Height:     200,
	Padding:    0.08,
}

type Config struct {
	Hues       int
	Colored    Color
	Grayscale  Color
	Background color.Color
	Width      int
	Height     int
	Padding    float64
}

type Color struct {
	Lightness  []float64
	Saturation float64
}

func (c Color) lightness(p float64) float64 {
	if len(c.Lightness) == 0 {
		return 0
	}
	if len(c.Lightness) == 1 || c.Lightness[0] == c.Lightness[1] || p == 0 {
		return c.Lightness[0]
	}
	if p >= 1 {
		return c.Lightness[1]
	}
	return c.Lightness[0] + c.Lightness[0]*p
}

func (c Color) color(hue float64, lightness float64) string {
	return correctedHsl(hue, c.Saturation, c.lightness(lightness))
}

/*
864444000141320028501e5a
^^ R color
864444000141320028501e5a
  ^^ G color
864444000141320028501e5a
    ^^ B color
864444000141320028501e5a
      ^^ A color
864444000141320028501e5a
        ^ single hue
864444000141320028501e5a
         ^^^ hue
864444000141320028501e5a
            ^^ color saturation
864444000141320028501e5a
              ^^ gray saturation
864444000141320028501e5a
                ^^ color lightness 1
864444000141320028501e5a
                  ^^ color lightness 2
864444000141320028501e5a
                    ^^ gray lightness 1
864444000141320028501e5a
                      ^^ gray lightness 2
*/

func ConfigFromString(h string) (c *Config, err error) {
	if len(h) != 24 {
		err = fmt.Errorf("invalid config length")
		return
	}
	var (
		hues            int
		R, G, B, A      int64
		colorSaturation int64
		colorLightness1 int64
		colorLightness2 int64
		graySaturation  int64
		grayLightness1  int64
		grayLightness2  int64
	)
	if R, err = strconv.ParseInt("0x"+h[0:2], 0, 64); err != nil {
		return
	}
	if G, err = strconv.ParseInt("0x"+h[2:4], 0, 64); err != nil {
		return
	}
	if B, err = strconv.ParseInt("0x"+h[4:6], 0, 64); err != nil {
		return
	}
	if A, err = strconv.ParseInt("0x"+h[6:8], 0, 64); err != nil {
		return
	}
	if colorSaturation, err = strconv.ParseInt("0x"+h[12:14], 0, 64); err != nil {
		return
	}
	if graySaturation, err = strconv.ParseInt("0x"+h[14:16], 0, 64); err != nil {
		return
	}
	if colorLightness1, err = strconv.ParseInt("0x"+h[16:18], 0, 64); err != nil {
		return
	}
	if colorLightness2, err = strconv.ParseInt("0x"+h[18:20], 0, 64); err != nil {
		return
	}
	if grayLightness1, err = strconv.ParseInt("0x"+h[20:22], 0, 64); err != nil {
		return
	}
	if grayLightness2, err = strconv.ParseInt("0x"+h[22:24], 0, 64); err != nil {
		return
	}
	if h[8] == '1' {
		var hue int64
		if hue, err = strconv.ParseInt("0x"+h[9:12], 0, 64); err != nil {
			return nil, err
		}
		hues = int(hue) - 1
		if hues > 360 {
			hues = 360
		}
	} else {
		hues = -1
	}

	c = &Config{
		Hues: hues,
		Colored: Color{
			Saturation: float64(colorSaturation) / 100,
			Lightness:  []float64{float64(colorLightness1) / 100, float64(colorLightness2) / 100},
		},
		Grayscale: Color{
			Saturation: float64(graySaturation) / 100,
			Lightness:  []float64{float64(grayLightness1) / 100, float64(grayLightness2) / 100},
		},
		Background: color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)},
		Width:      200,
		Height:     200,
		Padding:    0.08,
	}
	return c, nil
}

func ConfigFromBytes(h []byte) (c *Config, err error) {
	return ConfigFromString(string(h))
}
