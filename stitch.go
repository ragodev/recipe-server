package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	zero = image.Point{0, 0}
)

func stitch(images []image.Image) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for i, simg := range images {
		draw.Draw(img, simg.Bounds().Add(image.Point{(i % 3) * 133, (i / 3) * 100}), simg, zero, draw.Src)
	}
	return img
}

func loadImages(fileNames []string) []image.Image {
	var images []image.Image
	for _, s := range fileNames {
		f, _ := os.OpenFile(s, os.O_RDONLY, 0644)
		img, _ := jpeg.Decode(f)
		images = append(images, img)
	}
	return images
}

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileNames(ingredients []string) []string {
	var ingredientImages []string
	for _, ingredient := range ingredients {
		ingredientFolder := strings.Join(strings.Split(strings.TrimSpace(ingredient), " "), "-")
		if !exists(path.Join("resized", ingredientFolder)) {
			continue
		}
		fileList := []string{}
		err := filepath.Walk(path.Join("resized", ingredientFolder), func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, ".jpg") || strings.Contains(path, ".JPG") {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if len(fileList) > 0 {
			ingredientImages = append(ingredientImages, fileList[rand.Intn(len(fileList))])
		}
	}
	dest := make([]string, len(ingredientImages))
	perm := rand.Perm(len(ingredientImages))
	for i, v := range perm {
		dest[v] = ingredientImages[i]
	}
	return dest
}

func makeFile(ingredients []string, title string) {
	// resizeEverything()
	fmt.Println(ingredients)
	rand.Seed(time.Now().Unix())
	fileNames := getFileNames(ingredients)
	fmt.Println(fileNames)
	images := loadImages(fileNames)
	img := stitch(images)
	b := bytes.NewBuffer(nil)
	jpeg.Encode(b, img, nil)
	ioutil.WriteFile(path.Join("./images/", title+".jpg"), b.Bytes(), 0644)
}
