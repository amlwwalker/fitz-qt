package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/amlwwalker/pdf-editor/utils"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

var ABSOLUTE_PATH_URL string

type QmlBridge struct {
	core.QObject
	business BusinessInterface

	//pdf editor functions (receiving)
	_ func(path string)                                             `slot:"openFile"`            //from qml
	_ func(path string)                                             `slot:"setWorkingDirectory"` //from qml
	_ func(path, fileName, errorType string, imageFile *gui.QImage) `slot:"saveEditedFile"`      //from qml

	//pdf editor functions (sending)
	_ func(p string)    `signal:"loadImage"`
	_ func(data string) `signal:"sendMessage"`
}

//setup functions to communicate between front end and back end

func (q *QmlBridge) ConfigureBridge(config Config) {
	//1. configure the hotloader
	q.business = BusinessInterface{}
	q.business.configureInterface()

	q.ConnectOpenFile(func(path string) {
		fmt.Println("opening " + path)
		q.SendMessage("requested to open " + path)
		path, imageFiles := q.openFileForProcessing(path)

		q.business.iModel.ClearFiles()
		for _, v := range imageFiles {
			var f = NewFile(nil)
			f.SetFilePath(path)
			f.SetFileName(v)
			q.business.iModel.AddFile(f)
		}
		q.SendMessage("created image files for " + path)

	})
	q.ConnectSetWorkingDirectory(func(path string) {
		// ABSOLUTE_PATH_URL =  strings.Replace(path, "file://", "", -1) + "/"
		fmt.Println("chosen path is " + path)

		fmt.Println("chosen path is " + path)
		if runtime.GOOS == "windows" {
			ABSOLUTE_PATH_URL = filepath.Join(strings.TrimPrefix(path, "file:///"), "")
		} else {
			ABSOLUTE_PATH_URL = filepath.Join(strings.TrimPrefix(path, "file://"), "")
		}

		fmt.Println("working directory set to " + ABSOLUTE_PATH_URL)
		q.SendMessage("working directory set to " + ABSOLUTE_PATH_URL)
		//get the files in the specified directory
		if files, err := getDownloadedFiles(); err != nil {
			//couldn't retrieve errors
			fmt.Println("error retrieving previously downloaded files", err)
			q.SendMessage("error retrieving previously downloaded files" + err.Error())
		} else {
			q.business.fModel.ClearFiles()
			for _, v := range files {
				var f = NewFile(nil)
				f.SetFilePath(v.FileName)
				q.SendMessage("adding file " + v.FileName)
				f.SetFileSize(strconv.Itoa(v.FileSize))
				q.business.fModel.AddFile(f)
			}
		}
	})
	q.ConnectSaveEditedFile(func(path, fileName, errorType string, imageFile *gui.QImage) {
		pathToSave := filepath.Join(path, fileName+". error "+errorType)
		fmt.Println("path to save: " + pathToSave)

		buff := core.NewQBuffer(q)
		buff.Open(core.QIODevice__ReadWrite)
		ok := imageFile.Save2(buff, "PNG", -1)
		fmt.Println("Save2", ok)
		data := buff.Data().ConstData()
		fmt.Println("len2:", len(data))
		buff.Close()

		//this wont work while set for POSIX
		path = strings.Replace(path, "original", "edit", -1) //the path of the file
		dirName := strings.Split(path, "edit")[0]            //the path of the directory
		fmt.Println("path is " + path)
		fmt.Println("dirName directory is " + dirName)
		if errorType == "" {
			errorType = "NONE"
		}
		fileName = strings.Replace(fileName, ".orig.", "."+errorType+".", -1)
		if err := ioutil.WriteFile(filepath.Join(path, fileName), []byte(data), 0644); err != nil {
			fmt.Println("error writing to file " + err.Error())
		}
		fmt.Println("WriteFile", filepath.Join(path, fileName))
	})
}

func getDownloadedFiles() ([]utils.File, error) {
	//scan the file system based on the file download location
	//get file name and file size
	//if a user clicks, we are going to open the file if we can
	var files []utils.File
	fileList, err := ioutil.ReadDir(ABSOLUTE_PATH_URL)

	if err != nil {
		return files, err
	}
	//just for debugging
	for _, f := range fileList {
		fmt.Println("name " + f.Name())
		var tmp utils.File
		tmp.FileName = f.Name()
		if strings.Contains(tmp.FileName, ".pdf") {
			tmp.FileSize = int(f.Size())
			files = append(files, tmp)
		}
	}
	return files, nil
}

func CreateDirIfNotExist(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}

// func (q *QmlBridge) openFileForProcessing(filePath string) (string, []string) {
// 	//most don''t need the pdf ext
// 	origDirPath := ABSOLUTE_PATH_URL
// 	subject := strings.Replace(filePath, ".pdf", "", -1)
// 	q.SendMessage("requested to open subject " + subject)
// 	var imageFiles []string

// 	fmt.Println("and.... getting ready to read the pdf " + filepath.Join(origDirPath, subject) + ".pdf")
// 	// return "", imageFiles
// 	// TODO: BUG. CRASHES ON WINDOWS
// 	doc, err := fitz.New(filepath.Join(origDirPath, subject) + ".pdf")
// 	if err != nil {
// 		q.SendMessage("fitz error " + err.Error())
// 	}
// 	q.SendMessage("and.... getting ready to read the directory " + origDirPath + "/" + subject)

// 	defer doc.Close()
// 	q.SendMessage("OK, reading directory " + origDirPath + "/" + subject)
// 	dirName, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject))

// 	originalDir, _ := CreateDirIfNotExist(filepath.Join(origDirPath, subject, "original"))
// 	// Extract pages as images
// 	q.SendMessage("OK, extracting information. Pause ")

// 	for n := 0; n < doc.NumPage(); n++ {
// 		img, err := doc.Image(n)
// 		if err != nil {
// 			q.SendMessage("CreateDirIfNotExist error" + err.Error())
// 		}
// 		q.SendMessage("OK, creating png at " + filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", n)) + "Pause ")

// 		f, err := os.Create(filepath.Join(originalDir, fmt.Sprintf(subject+".orig.%03d.png", n)))
// 		if err != nil {
// 			q.SendMessage("create image error " + err.Error())
// 		}

// 		if err = png.Encode(f, img); err != nil {
// 			q.SendMessage("image encode error " + err.Error())
// 		}

// 		f.Close()
// 		imageFiles = append(imageFiles, fmt.Sprintf(subject+".orig.%03d.png", n))
// 	}
// 	CreateDirIfNotExist(filepath.Join(origDirPath, subject, "edit"))

// 	// Extract pages as text
// 	for n := 0; n < doc.NumPage(); n++ {
// 		text, err := doc.Text(n)
// 		if err != nil {
// 			q.SendMessage("read pdf for text error " + err.Error())
// 		}

// 		f, err := os.Create(filepath.Join(dirName, fmt.Sprintf(subject+".%03d.txt", n)))
// 		if err != nil {
// 			q.SendMessage("create text error " + err.Error())
// 		}

// 		_, err = f.WriteString(text)
// 		if err != nil {
// 			q.SendMessage("write image error " + err.Error())
// 		}

// 		f.Close()
// 	}
// 	//only returns the first image in the pdf
// 	//i.e pdf's should be just 1 page long
// 	return originalDir, imageFiles
// }
