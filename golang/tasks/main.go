package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func main() {
	godotenv.Load(".env")
	// open "test.jpg"
	file, err := os.Open(os.Getenv("IMAGE_FOLDER") + "i1.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// resize to using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(250, 200, img, resize.Lanczos3)

	out, err := os.Create(os.Getenv("IMAGE_FOLDER_RESIZE") + "i1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
