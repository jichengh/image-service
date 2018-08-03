package main

import (
	"github.com/k0kubun/pp"

	//"io"

	"bytes"

	//"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", home)
	e.GET("/image", getImage)
	e.POST("/analyze_image", analyzeImage)
	e.POST("/resize_image", resizeImage)
	e.POST("/create_thumbnail_by_width", createThumbnailByWidth)
	e.Logger.Fatal(e.Start(":8888"))
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func getImage(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func analyzeImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	image, format, err := image.DecodeConfig(src)
	if err != nil {
		return err
	}

	var response struct {
		Filename string `json:"filename"`
		FileSize string `json:"filesize"`
		Format   string `json:"format"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
	}
	response.Filename = file.Filename
	response.FileSize = strconv.FormatInt(file.Size, 10) + " bytes"
	response.Format = format
	response.Width = image.Width
	response.Height = image.Height

	return c.JSON(http.StatusOK, response)
}

func resizeImage(c echo.Context) error {
	width, err := strconv.Atoi(c.FormValue("width"))
	if err != nil {
		log.Printf("Failed to parse data: " + err.Error())
		return errors.New("Failed to parse data: " + err.Error())
	}

	height, err := strconv.Atoi(c.FormValue("height"))
	if err != nil {
		log.Printf("Failed to parse data: " + err.Error())
		return errors.New("Failed to parse data: " + err.Error())
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Failed to get file: " + err.Error())
		return errors.New("Failed to get file: " + err.Error())
	}

	src, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: " + err.Error())
		return errors.New("Failed to open file: " + err.Error())
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		log.Printf("Failed to decode image: " + err.Error())
		return errors.New("Failed to decode image: " + err.Error())
	}

	resizeImg := imaging.Resize(img, width, height, imaging.Lanczos)

	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, resizeImg, nil)
	if err != nil {
		log.Printf("Failed to encode image: " + err.Error())
		return errors.New("Failed to encode image: " + err.Error())
	}

	return c.Blob(http.StatusOK, "image/jpeg", imageBuf.Bytes())
}

func createThumbnailByWidth(c echo.Context) error {
	width, err := strconv.Atoi(c.FormValue("width"))
	if err != nil {
		log.Printf("Failed to parse data: " + err.Error())
		return errors.New("Failed to parse data: " + err.Error())
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Failed to get file: " + err.Error())
		return errors.New("Failed to get file: " + err.Error())
	}

	src, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: " + err.Error())
		return errors.New("Failed to open file: " + err.Error())
	}
	defer src.Close()

	imgConfig, format, err := image.DecodeConfig(src)
	if err != nil {
		log.Printf("Failed to get image config: " + err.Error())
		return errors.New("Failed to get image config: " + err.Error())
	}

	// rewind src (io.Reader)
	src.Seek(0, 0)

	img, err := imaging.Decode(src)
	if err != nil {
		log.Printf("Failed to decode image: " + err.Error())
		return errors.New("Failed to decode image: " + err.Error())
	}

	ratio := float64(width) / float64(imgConfig.Width)
	height := int(ratio * float64(imgConfig.Height))

	pp.Println(ratio, width, height)

	resizeImg := imaging.Thumbnail(img, width, height, imaging.Linear)

	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, resizeImg, nil)
	if err != nil {
		log.Printf("Failed to encode image: " + err.Error())
		return errors.New("Failed to encode image: " + err.Error())
	}
	pp.Println(format)
	return c.Blob(http.StatusOK, "image/jpeg", imageBuf.Bytes())
}
