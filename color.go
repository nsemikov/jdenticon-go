package jdenticon

import (
	"image/color"
	"strconv"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func (j *jdenticon) hue() float64 {
	var hue float64
	if j.config.Hues == -1 {
		h, _ := strconv.ParseInt("0x"+j.hash[len(j.hash)-7:], 0, 64)
		hue = float64(h) / 0xfffffff
	} else {
		hue = float64(j.config.Hues) / 360
	}
	return hue
}

func (j *jdenticon) theme() []string {
	hue := j.hue()
	darkgray := j.config.Grayscale.color(hue, 0)
	midcolor := j.config.Colored.color(hue, 0.5)
	lightgray := j.config.Grayscale.color(hue, 1)
	lightcolor := j.config.Colored.color(hue, 1)
	darkcolor := j.config.Colored.color(hue, 0)
	return []string{darkgray, midcolor, lightgray, lightcolor, darkcolor}
}

func correctedHsl(h, s, l float64) string {
	correctors := []float64{0.55, 0.5, 0.5, 0.46, 0.6, 0.55, 0.55}
	corrector := correctors[int(h*6+0.5)]
	// Adjust the input lightness relative to the corrector
	if l < 0.5 {
		l = l * corrector * 2
	} else {
		l = corrector + (l-0.5)*(1-corrector)*2
	}
	return colorful.Hsl(h*360, s, l).Hex()
}

func (j *jdenticon) colors() []string {
	theme := j.theme()
	available := []string{}
	var (
		dark  bool
		light bool
	)
	for i := 0; i < 3; i++ {
		s := j.hash
		s = "0x" + s[i+8:i+9]
		h, _ := strconv.ParseInt(s, 0, 64)
		idx := int(h) % len(theme)
		if idx == 0 || idx == 4 {
			if dark {
				idx = 1
			}
			dark = true
		}
		if idx == 2 || idx == 3 {
			if light {
				idx = 1
			}
			light = true
		}
		available = append(available, theme[idx])
	}
	return available
}

func toHex(c color.Color) string {
	cc, _ := colorful.MakeColor(c)
	return cc.Hex()
}

func opacity(c color.Color) float64 {
	_, _, _, a := c.RGBA()
	return float64(uint8(a)) / 255
}
