package main

import (
	"fmt"
	"log"
	"net/http"
	"self/imagemaker/util"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func ResultImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := util.RenderPage(w, "imageresult", map[string]string{
		"ImageResult": "/assets/images/output.png",
	})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/uploadfile", UploadFile)
	router.POST("/uploadfile", UploadFile)
	router.GET("/resultimage", ResultImage)
	router.POST("/resultimage", ResultImage)
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	log.Println("Serving on 8080")

	log.Fatal(http.ListenAndServe(":8080", router))
	//log.Fatal(http.ListenAndServe(":8081", http.FileServer(http.Dir("./assets/images/default.png"))))
}
