package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
)

type pixel struct {
	r, g, b, a uint8
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("respond:", err)
	}
}

//HandleUploadImage ...
func HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	fmt.Println(r.ContentLength)
	err := r.ParseMultipartForm(2000)
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		fmt.Println("ParseMultipartForm", body)
		respond(w, r, http.StatusBadRequest, body)
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		fmt.Println("FormFile", file, handler, body)
		respond(w, r, http.StatusBadRequest, body)
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		respond(w, r, http.StatusBadRequest, body)
		return
	}

	bounds := img.Bounds()

	slidingWindow(img, 30, 30)

	data := "File processed with success. File name: " + handler.Filename + " " + bounds.String()
	body := SuccessMessage{Message: data}
	respond(w, r, http.StatusOK, body)
}

func slidingWindow(img image.Image, width int, heigth int) {
	bounds := img.Bounds()

	// se a a resolução da imagem for menor que do parametro, fazer o resize.
	// "github.com/nfnt/resize"
	// img = resize.Resize(resX, resY, img, resize.Lanczos3)

	for y := heigth; y <= bounds.Max.Y; y += heigth {
		for x := width; x <= bounds.Max.X; x += width / 2 {
			fmt.Println(x, y)
			imgRect := image.Rect(0, 0, x, y)
			imgNew := image.NewGray(imgRect)
			go processingImage(imgNew)
		}
	}

}

func processingImage(img image.Image) []pixel {
	imgGray := escalaCinza(img)
	imgBeW := escalaPretoBranco(imgGray)
	imgBack := checkBackground(imgBeW)
	pixels := getPixels(imgBack)
	return pixels
}

func getPixels(img image.Image) []pixel {

	bounds := img.Bounds()
	pixels := make([]pixel, bounds.Dx()*bounds.Dy())

	i := 0
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[i].r = uint8(r)
			pixels[i].g = uint8(g)
			pixels[i].b = uint8(b)
			pixels[i].a = uint8(a)
			i++
		}
	}

	return pixels
}

func escalaCinza(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	imgRect := image.Rect(0, 0, w, h)
	gray := image.NewGray(imgRect)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}
	return gray
}

func escalaPretoBranco(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	imgRect := image.Rect(0, 0, w, h)
	gray := image.NewGray(imgRect)
	total := uint32(0)
	media := uint32(0)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			total = total + r
		}
	}

	media = total / uint32(w*h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, _ := img.At(x, y).RGBA()

			if r > media {
				r = 255
			} else {
				r = 0
			}

			gray.Set(x, y, color.Gray{uint8(r)})

		}
	}
	return gray
}

func checkBackground(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	imgRect := image.Rect(0, 0, w, h)
	gray := image.NewGray(imgRect)
	changeBackground := false
	total := uint32(0)
	totalEsquerda := uint32(0)
	totalDireita := uint32(0)
	totalBaixo := uint32(0)
	totalCima := uint32(0)

	for y := 0; y < h; y++ {
		r, _, _, _ := img.At(0, y).RGBA()
		totalEsquerda = totalEsquerda + r
		r, _, _, _ = img.At(w, y).RGBA()
		totalDireita = totalDireita + r
	}

	for x := 0; x < w; x++ {
		r, _, _, _ := img.At(x, 0).RGBA()
		totalBaixo = totalBaixo + r
		r, _, _, _ = img.At(x, h).RGBA()
		totalCima = totalCima + r
	}

	total = totalBaixo + totalCima + totalDireita + totalEsquerda

	if total < 1966050 {
		changeBackground = true
	}

	if changeBackground {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, _, _, _ := img.At(x, y).RGBA()

				if r == 0 {
					r = 255
				} else {
					r = 0
				}

				gray.Set(x, y, color.Gray{uint8(r)})

			}
		}
		return gray
	}
	return img
}
