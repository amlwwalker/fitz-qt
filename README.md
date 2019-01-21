# Demo of using fitz (pdf converter)

This is a demo of using the fitz pdf converter library with therecipe/qt library.

Currently works on Mac, but cannot get it to work on Windows. Fitz causes a segmentation fault


The application offers an interface to be able to load all the pages of a pdf as PNG files and paint over the top of them, then resave them out

## Instructions

* Compile inside the `qt` directory with `qtdeploy test desktop`
* Once it compiles and runs you will see a window open with a bar like

![https://imgur.com/a/eqDfjCj]()

* To use the app

    1. Button most left of toolbar: Set Working Directory (Choose a directory that has a PDF in (helpful if nothing else in there at this stage))
    2. Button that looks like a PDF: Choose a PDF from the ones it found in that directory (this converts the pages to images so can take a moment)
    3. On the right of the tool bar is a button that that allows you to choose a page to edit (the PNG)
      *. The image will load in the middle
    4. Click on the paintbrush to enable painting. Now click once, someone on the image and move the mouse around. Click again to stop painting
    5. Click the floppy disk. It will automatically save the files inside of the folder with the same name as the pdf
