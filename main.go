package main

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	colors "gitlab.com/ethanbakerdev/colors"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/env/extension"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/example/messaging"
)

func main() {
	// custom init config
	connector.InitWithOptions(env.WithApiHost(httplib.URL_EU),
		env.WithMsgHost(httplib.MSG_EU),
		env.WithAccessID("59d44w87kqge874ckyre"),
		env.WithAccessKey("774f34c2367a45c580a142072c5fc3f5"),
		env.WithAppName("MedmedLand"),
		env.WithDebugMode(false),
	)

	go messaging.Listener()

	// r := handler.NewGinEngin()
	// go r.Run("0.0.0.0:2021")
	inputControl()
	watitSignal()
}

func inputControl() {
	for {
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

		fmt.Printf(
			`Average colour: rgb(%v)\nAverage colour: hsv({"h":%d,"s":%d,"v":%d})\n`,
			avgColor,
			avgColor.H, avgColor.S, avgColor.V,
		)
		// resp := &model.PostDeviceCmdResponse{}

		err = connector.MakePostRequest(
			context.Background(),
			connector.WithAPIUri(fmt.Sprintf("/v1.0/iot-03/devices/%s/commands", "545000508cce4ee25ac7")),
			connector.WithPayload([]byte(fmt.Sprintf(`{
				"commands": [
				  {
					"code": "colour_data",
					"value": {"h":%d,"s":255,"v":255}
				  },
				  {
					"code": "bright",
					"value": 100
				  },
				  {
					"code": "saturation",
					"value": 100
				  }
				]
			  }`, avgColor.H))))
		if err != nil {
			logger.Log.Errorf("err:%s", err.Error())
			return
		}

		time.Sleep(500 * time.Millisecond)

	}

}

func watitSignal() {
	go inputControl()
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		select {
		case c := <-quitCh:
			extension.GetMessage(constant.TUYA_MESSAGE).Stop()
			logger.Log.Infof("receive sig:%v, shutdown the http server...", c.String())
			return
		}
	}
}
