package effect

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/MehdiZouareg/ambituya/config"
	"github.com/kbinani/screenshot"
	"github.com/rs/zerolog/log"
	"github.com/tuya/tuya-connector-go/connector"

	colors "gitlab.com/ethanbakerdev/colors"
)

func Ambilight(cfg *config.Config) {

	devices := cfg.TuyaRegisteredDevices
	refreshRate := cfg.Ambilight.RefreshRate

	for {
		// Input Screen Number here
		bounds := screenshot.GetDisplayBounds(0)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}

		imgSize := img.Bounds().Size()

		var redSum float64
		var greenSum float64
		var blueSum float64

		for x := 0; x < imgSize.X; x++ {
			for y := 0; y < imgSize.Y; y++ {
				pixel := img.At(x, y)
				col := color.RGBAModel.Convert(pixel).(color.RGBA)

				redSum += float64(col.R)
				greenSum += float64(col.G)
				blueSum += float64(col.B)
			}
		}

		imgArea := float64(imgSize.X * imgSize.Y)

		redAverage := math.Round(redSum / imgArea)
		greenAverage := math.Round(greenSum / imgArea)
		blueAverage := math.Round(blueSum / imgArea)

		avgColor := colors.RGBtoHSV(colors.RGB{
			R: int(redAverage),
			G: int(greenAverage),
			B: int(blueAverage),
		})

		command := fmt.Sprintf(`{
			"commands": [
			  {
				"code": "colour_data",
				"value": {
				  "h": %v,
				  "s": 255,
				  "v": 255
				}
			  }
			]
		  }`, avgColor.H)

		for _, device := range devices {
			err = connector.MakePostRequest(
				context.Background(),
				connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", device.ID)),
				connector.WithPayload([]byte(command)))
			if err != nil {
				log.Error().Err(err).Msg("got error while reaching tuya api")
				return
			}
		}

		time.Sleep(time.Duration(refreshRate) * time.Millisecond)
	}
}
