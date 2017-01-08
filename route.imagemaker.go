package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	png "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"self/imagemaker/util"
	"strconv"
	"time"

	"image/jpeg"

	"github.com/julienschmidt/httprouter"
)

func UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//[post]
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		err := util.RenderPage(w, "imagemaker", map[string]string{
			"Token": token,
		})
		if err != nil {
			log.Println(err)
		}

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		destFile, err := os.Create("assets/images/" + handler.Filename)
		defer destFile.Close()

		io.Copy(destFile, file)
		destFile.Sync()

		err = OverlayImage("assets/images/default.png", "assets/images/"+handler.Filename)
		if err != nil {
			log.Println(err)
		}

		/*
			fmt.Fprintf(w, "%v", handler.Header)
			f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)

			PostFile(handler.Filename, "http://127.0.0.1:8080/uploadfile")
		*/
	}
}

func PostFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func OverlayImage(templateImage string, baseImage string) error {
	imgFile1, err := os.Open(baseImage)
	imgFile2, err := os.Open(templateImage)
	if err != nil {
		fmt.Println(err)
	}

	var img1 image.Image
	if baseImage[len(baseImage)-4:] == ".jpg" {
		var err error
		newImgPath := "assets/images/temp.png"

		img1, err = jpeg.Decode(imgFile1)
		if err != nil {
			return err
		}
		if err := ConvertToPNG(img1); err != nil {
			return err
		}
		imgFile1, err = os.Open(newImgPath)
		img1, _, err = image.Decode(imgFile1)
	} else {
		img1, err = png.Decode(imgFile1)
		if err != nil {
			return err
		}
	}
	img2, _, err := image.Decode(imgFile2)
	if err != nil {
		fmt.Println(err)
	}

	//starting position of the second image (bottom left)
	sp2 := image.Point{img1.Bounds().Dx() - img2.Bounds().Dx(), 0}
	fmt.Printf("%d %d %d %+v\n", img2.Bounds().Dx(), img1.Bounds().Dx(), img2.Bounds().Dx()-img1.Bounds().Dx(), sp2)
	//new rectangle for the second image
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}

	fmt.Printf("%+v %+v\n", sp2, sp2.Add(img2.Bounds().Size()))
	//rectangle for the big image
	r := image.Rectangle{image.Point{0, 0}, r2.Max}

	rgba := image.NewRGBA(r)
	fmt.Printf("%+v %+v\n", image.Point{0, 0}, r2.Max)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Over)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Over)

	out, err := os.Create("assets/images/output.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)

	/*
		var opt jpeg.Options
		opt.Quality = 80

		jpeg.Encode(out, rgba, &opt)
	*/

	return nil
}

func ConvertToPNG(img image.Image) error {
	out, err := os.Create("assets/images/temp.png")
	if err != nil {
		return err
	}

	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
