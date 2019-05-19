package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	colorful "github.com/lucasb-eyer/go-colorful"

	jdenticon "github.com/stdatiks/jdenticon-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.CSRF())
	e.Use(middleware.RequestID())
	e.File("/", "test.html")
	e.GET("/icon/", icon)
	e.GET("/icon/:identity", icon)
	e.Logger.Fatal(e.Start(":1323"))
}

func icon(c echo.Context) error {
	identity := c.Param("identity")

	hues := jdenticon.DefaultConfig.Hues
	colorLightness1 := jdenticon.DefaultConfig.Colored.Lightness[0]
	colorLightness2 := jdenticon.DefaultConfig.Colored.Lightness[1]
	colorSaturation := jdenticon.DefaultConfig.Colored.Saturation
	grayscaleLightness1 := jdenticon.DefaultConfig.Grayscale.Lightness[0]
	grayscaleLightness2 := jdenticon.DefaultConfig.Grayscale.Lightness[1]
	grayscaleSaturation := jdenticon.DefaultConfig.Grayscale.Saturation
	width := jdenticon.DefaultConfig.Width
	height := jdenticon.DefaultConfig.Height
	padding := jdenticon.DefaultConfig.Padding
	background := jdenticon.DefaultConfig.Background
	var config *jdenticon.Config

	for name := range c.QueryParams() {
		switch name {
		case "hues":
			if val, err := strconv.Atoi(c.QueryParam(name)); err == nil {
				hues = val
			}
		case "colorLightness1":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				colorLightness1 = val
			}
		case "colorLightness2":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				colorLightness2 = val
			}
		case "colorSaturation":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				colorSaturation = val
			}
		case "grayscaleLightness1":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				grayscaleLightness1 = val
			}
		case "grayscaleLightness2":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				grayscaleLightness1 = val
			}
		case "grayscaleSaturation":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				colorSaturation = val
			}
		case "width":
			if val, err := strconv.Atoi(c.QueryParam(name)); err == nil {
				width = int(val)
			}
		case "height":
			if val, err := strconv.Atoi(c.QueryParam(name)); err == nil {
				height = int(val)
			}
		case "padding":
			if val, err := strconv.ParseFloat(c.QueryParam(name), 8); err == nil {
				padding = val
			}
		case "background":
			if val, err := colorful.Hex(c.QueryParam(name)); err == nil {
				background = val
			}
		case "config":
			if val, err := jdenticon.ConfigFromString(c.QueryParam(name)); err == nil {
				config = val
			} else {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	// icon := jdenticon.New(identity)
	if config == nil {
		config = &jdenticon.Config{}
		config.Hues = hues
		config.Background = background
		config.Colored.Lightness = []float64{colorLightness1, colorLightness2}
		config.Colored.Saturation = colorSaturation
		config.Grayscale.Lightness = []float64{grayscaleLightness1, grayscaleLightness2}
		config.Grayscale.Saturation = grayscaleSaturation
		config.Width = width
		config.Height = height
		config.Padding = padding
	}
	icon := jdenticon.NewWithConfig(identity, config)
	svg, err := icon.SVG()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.Error(err)
	}
	return c.Blob(http.StatusOK, "image/svg+xml", svg)
}
