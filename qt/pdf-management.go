package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func (q *QmlBridge) openFileForProcessing(filePath string) (string, []string) {
	//most don''t need the pdf ext
	origDirPath := ABSOLUTE_PATH_URL
	subject := strings.Replace(filePath, ".pdf", "", -1)
	go q.SendMessage("requested to open subject " + subject + " THEN PAUSING")
	var imageFiles []string
	time.Sleep(4 * time.Millisecond)
	fmt.Println("and.... getting ready to read the pdf " + filepath.Join(origDirPath, subject) + ".pdf")
	// return "", imageFiles
	pdfPath := filepath.Join(origDirPath, subject) + ".pdf"

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	// pw := imagick.NewPixelWand()
	defer mw.Destroy()
	mw.SetResolution(300, 300)
	mw.SetCompressionQuality(100)
	mw.SetImageFormat("png")
	if err := mw.ReadImage(pdfPath); err != nil {
		fmt.Println("Error reading pages ", err)
	}
	pages := mw.GetNumberImages()

	go q.SendMessage("OK, reading pdf " + pdfPath)
	dirName, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject))
	originalDir, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject, "original"))
	fmt.Println("pages ", pages)
	for n := uint(0); n < pages; n++ {
		//get the first pdf page to convert to an image
		go q.SendMessage("OK, creating png at " + filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", n)))

		mw.SetIteratorIndex(int(n)) // This being the page offset
		fmt.Println("index ", mw.GetIteratorIndex())
		// s, _ := strconv.ParseUint(mw.GetIteratorIndex(), 10, 32)
		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_OFF)
		mw.WriteImage(filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", mw.GetIteratorIndex())))

		// png := //the resulting image from the conversion of this page of the pdf
		// //write it out now as a png
		// f, err := os.Create(filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", n)))
		// if err != nil {
		// 	q.SendMessage("create image error " + err.Error())
		// }

		// if err = png.Encode(f, img); err != nil {
		// 	q.SendMessage("image encode error " + err.Error())
		// }

		// f.Close()
		imageFiles = append(imageFiles, fmt.Sprintf(subject+".orig.%03d.png", n))
	}

	//now setup the edit directory and the csv
	CreateDirIfNotExist(filepath.Join(origDirPath, subject, "edit"))

	//create a csv file for the directory
	//scoped code.
	{
		f, err := os.Create(filepath.Join(dirName, "drawing_data.csv"))
		if err != nil {
			q.SendMessage("create csv error " + err.Error())
		}
		defer f.Close()

		if _, err = f.WriteString("drawing_name,error_type,path\r\n"); err != nil {
			q.SendMessage("write error " + err.Error())
		}
		f.Close()
	}
	return originalDir, imageFiles
}

func pdfConverter() {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	// pw := imagick.NewPixelWand()
	defer mw.Destroy()
	mw.SetResolution(300, 300)
	mw.SetCompressionQuality(100)
	mw.SetImageFormat("png")
	mw.ReadImage("ML-MT-PH-00005.000_REV1A.pdf")
	pages := mw.GetNumberImages()

	for i := uint(0); i < pages; i++ {
		mw.SetIteratorIndex(int(i)) // This being the page offset
		fmt.Println("index ", mw.GetIteratorIndex())
		// s, _ := strconv.ParseUint(mw.GetIteratorIndex(), 10, 32)
		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_OFF)
		mw.WriteImage("test" + fmt.Sprint(mw.GetIteratorIndex()) + ".png")
	}

}

// background := imagick.NewPixelWand()
// background.SetAlpha(0)
// background.SetColor("red")
// mw.SetBackgroundColor(background)
