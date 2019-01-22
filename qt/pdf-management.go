package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func (q *QmlBridge) openFileForProcessing(filePath string) (string, []string) {
	//most don''t need the pdf ext
	origDirPath := ABSOLUTE_PATH_URL
	subject := strings.Replace(filePath, ".pdf", "", -1)
	q.SendMessage("requested to open subject " + subject + " THEN PAUSING")
	var imageFiles []string
	fmt.Println("and.... getting ready to read the pdf " + filepath.Join(origDirPath, subject) + ".pdf")
	pdfPath := filepath.Join(origDirPath, subject) + ".pdf"

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.SetResolution(300, 300)
	mw.SetCompressionQuality(100)
	mw.SetImageFormat("png")
	if err := mw.ReadImage(pdfPath); err != nil {
		fmt.Println("Error reading pages ", err)
	}
	pages := mw.GetNumberImages()

	q.SendMessage("OK reading pdf " + pdfPath)
	// dirName, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject))
	originalDir, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject, "original"))
	fmt.Println("pages ", pages)
	for n := uint(0); n < pages; n++ {
		//get the first pdf page to convert to an image
		q.SendMessage("OK, creating png at " + filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", n)))

		mw.SetIteratorIndex(int(n)) // This being the page offset
		fmt.Println("index ", mw.GetIteratorIndex())
		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_OFF)
		mw.WriteImage(filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", mw.GetIteratorIndex())))

		imageFiles = append(imageFiles, fmt.Sprintf(subject+".orig.%03d.png", n))
	}

	//now setup the edit directory and the csv
	CreateDirIfNotExist(filepath.Join(origDirPath, subject, "edit"))

	return originalDir, imageFiles
}
