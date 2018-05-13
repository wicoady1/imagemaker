package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wicoady1/gowatermark"
	"github.com/wicoady1/imagemaker/util"

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

		//-----

		file2, handler2, err2 := r.FormFile("overlayfile")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		defer file2.Close()

		destFile2, err2 := os.Create("assets/images/" + handler2.Filename)
		defer destFile2.Close()

		io.Copy(destFile2, file2)
		destFile2.Sync()

		err = OverlayImage("assets/images/"+handler2.Filename, "assets/images/"+handler.Filename)
		if err != nil {
			log.Println(err)
		}

		err = util.RenderPage(w, "imageresult", map[string]string{
			"ImageResult": "/assets/images/output.png",
		})
		if err != nil {
			log.Println(err)
		}

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
	imageType := gowatermark.ImagePNG
	if baseImage[len(baseImage)-4:] == ".jpg" || baseImage[len(baseImage)-5:] == ".jpeg" {
		imageType = gowatermark.ImageJPEG
	}
	mainImage, err := gowatermark.New(baseImage, imageType)
	if err != nil {
		log.Println(err)
		return err
	}

	imageType2 := gowatermark.ImagePNG
	if templateImage[len(templateImage)-4:] == ".jpg" || templateImage[len(templateImage)-5:] == ".jpeg" {
		imageType2 = gowatermark.ImageJPEG
	}
	err = mainImage.AddOverheadImage(templateImage, imageType2)
	if err != nil {
		log.Println(err)
		return err
	}

	err = mainImage.OutputImageToFile("assets/images/", "output", gowatermark.ImagePNG)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
