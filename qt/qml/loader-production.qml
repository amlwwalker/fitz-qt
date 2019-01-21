import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.4
import QtQuick.Controls.Material 2.0
import QtQuick.Controls.Universal 2.0
import Qt.labs.settings 1.0
import QtQuick.Dialogs 1.0
Item {
    id: window
    width: 900
    height: 700
    visible: true
    property int xpos
    property int ypos
    property bool painting: false
    property int brushSize: 35
    property string brushColor: "red"
    property bool mousePressed: false
    property string currentlySelectedFilePath: ""
    property string currentlySelectedFileName: ""
         ToolBar {
            id: toolbar
            Material.foreground: "white"
            Material.background: Material.BlueGrey
             z: 100
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.top: parent.top
            RowLayout {
                spacing: 20
                anchors.fill: parent
	            ToolButton {
	                contentItem: Image {
	                    fillMode: Image.Pad
	                    horizontalAlignment: Image.AlignHCenter
	                    verticalAlignment: Image.AlignVCenter
	                    source: "images/FA/black/png/22/wrench.png"
	                }
	                onClicked: {
                        setWorkingDirectory.open()
	                }
	            }                
	            ToolButton {
	                contentItem: Image {
	                    fillMode: Image.Pad
	                    horizontalAlignment: Image.AlignHCenter
	                    verticalAlignment: Image.AlignVCenter
	                    source: "images/FA/black/png/22/file-pdf-o.png"
	                }
	                onClicked: {
                        drawer.open()
	                }
	            }
                ToolButton {
                    contentItem: Image {
                        id: paintButton
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/black/png/22/paint-brush.png"
                    }
                    onClicked: {
                        painting = !painting
                        if (painting) {
                            paintButton.source = "images/FA/white/png/22/paint-brush.png"
                        } else {
                            paintButton.source = "images/FA/black/png/22/paint-brush.png"
                        }
                    }
                }
                ToolButton {
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/black/png/22/plus.png"
                    }
                    onClicked: {
                    brushSize += 5
                    }
                }
                ToolButton {
                    contentItem: Image {
                            fillMode: Image.Pad
                            horizontalAlignment: Image.AlignHCenter
                            verticalAlignment: Image.AlignVCenter
                            source: "images/FA/black/png/22/minus.png"
                    }
                    onClicked: {
                    brushSize -= 5
                    }
                }
                Label {
                    id: brushSizeLabel
                    text: brushSize
                    font.pixelSize: 20
                    elide: Label.ElideRight
                    horizontalAlignment: Qt.AlignHCenter
                    verticalAlignment: Qt.AlignVCenter
                }
                ToolButton {
                    id: clearCanvas
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/black/png/22/eraser.png"
                    }
                    onClicked: {
                        painting = false
                        paintButton.source = "images/FA/black/png/22/paint-brush.png"
                        myCanvas.clear()
                    }
                }
                ToolButton {
                    id: settingsViewer
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/black/png/22/save.png"
                    }
                    onClicked: {
                        painting = false
                        paintButton.source = "images/FA/black/png/22/paint-brush.png"
                    if (currentlySelectedFilePath == "" || currentlySelectedFileName == "") {
                        return false
                    }
                    source.grabToImage(function(result){
                        QmlBridge.saveEditedFile(currentlySelectedFilePath, currentlySelectedFileName, errorType.text, result.image)
                    
                })
                    }
                }
                Label {
                    id: titleLabel
                    text: "Fitz Tester"
                    font.pixelSize: 20
                    elide: Label.ElideRight
                    horizontalAlignment: Qt.AlignHCenter
                    verticalAlignment: Qt.AlignVCenter
                    Layout.fillWidth: true
                }
	            ToolButton {
	                contentItem: Image {
	                    fillMode: Image.Pad
	                    horizontalAlignment: Image.AlignHCenter
	                    verticalAlignment: Image.AlignVCenter
	                    source: "images/FA/black/png/22/image.png"
	                }
	                onClicked: {

                        drawerRight.open()
	                }
	            }
            }
        }
        ToolBar {
            id: subToolBar
            Material.foreground: "white"
            Material.background: Material.Red
             z: 100
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.top: toolbar.bottom
            RowLayout {
                spacing: 20
                anchors.fill: parent
	            ToolButton {
	                contentItem: TextField {
                        id: errorType
                        verticalAlignment: TextInput.AlignVCenter
                        placeholderText: "insert error type for this edit"
                        color: "black"
                        background: Rectangle {
                            radius: 2
                            implicitWidth: 300
                            implicitHeight: 24
                            border.color: "#333"
                            border.width: 1
                        }
                    }
	            }
            }
        }
        FileDialog {
            id: setWorkingDirectory
            selectFolder: true
            onAccepted: {
                QmlBridge.setWorkingDirectory(setWorkingDirectory.folder)
            }
        }
    //menu
    Drawer {
        id: drawer
        width: 200
        height: window.height
        edge: Qt.RightEdge
        ListView {
            id: listView
            currentIndex: -1
            anchors.fill: parent

            delegate: ItemDelegate {
                width: parent.width
                text: model.filePath
                highlighted: ListView.isCurrentItem
                onClicked: {
                    if (listView.currentIndex != index) {
                        listView.currentIndex = index
                        titleLabel.text = model.filePath
                        QmlBridge.openFile(model.filePath)
                    }
                    drawer.close()
                }
            }
            model: FilesModel

            ScrollIndicator.vertical: ScrollIndicator { }
        }
    }
    Drawer {
        id: drawerRight
        width: Math.min(window.width, window.height) / 3 * 2
        height: window.height
        edge: Qt.RightEdge
        ListView {
            id: listViewRight
            currentIndex: -1
            anchors.fill: parent

            delegate: ItemDelegate {
                width: parent.width
                text: model.fileName
                highlighted: ListView.isCurrentItem
                onClicked: {
                    if (listViewRight.currentIndex != index) {
                        listViewRight.currentIndex = index
                        titleLabel.text = model.fileName
                        currentlySelectedFilePath = model.filePath
                        currentlySelectedFileName = model.fileName
                        img.source = "file:///" + model.filePath + "/" + model.fileName
                        myCanvas.requestPaint()
                        flickableCanvas.contentWidth = img.width
                        flickableCanvas.contentHeight = img.height
                    }
                    drawerRight.close()
                }
            }
            model: ImageFilesModel

            ScrollIndicator.vertical: ScrollIndicator { }
        }
    }
    ScrollView {
            anchors.left: parent.left
            anchors.top: subToolBar.bottom
            anchors.right: parent.right
            clip: true
            ScrollBar.horizontal.policy: ScrollBar.AlwaysOn
            ScrollBar.vertical.policy: ScrollBar.AlwaysOn

            Flickable {
                id: flickableCanvas
                anchors.fill: parent
                
                Item {
                    id: source
                    width: Math.max(parent.width,img.width)
                    height: Math.max(parent.height,img.height)
                Image
                {
                    id: img
                }
                Canvas {
                    id: myCanvas
                    anchors.left: img.left
                    anchors.right: img.right
                    anchors.top: img.top
                    anchors.bottom: img.bottom
                    onPaint: {
                        var ctx = getContext('2d')
                        ctx.fillStyle = brushColor
                        xpos = mousearea.mouseX;
                        ypos = mousearea.mouseY;
                        if (painting)
                        ctx.fillRect(xpos-7, ypos-7, brushSize, brushSize)
                    }
                    function clear() {
                        var ctx = getContext("2d");
                        ctx.reset();
                        myCanvas.requestPaint();
                    }              
                    MouseArea{
                        id:mousearea
                        anchors.fill: parent
                        hoverEnabled: mousePressed
                        onClicked: {
                            xpos = mouseX
                            ypos = mouseY
                            mousePressed = !mousePressed
                            if (painting)
                                myCanvas.requestPaint()
                        }
                        onPositionChanged: {
                            xpos = mouseX
                            ypos = mouseY
                            if (mousePressed) {
                                myCanvas.requestPaint();
                            }
                        }
                    }
                }
            }
        }
        }
    Connections {
        target: QmlBridge
        onLoadImage: {
            img.source = "file://" + p
            myCanvas.requestPaint()
            flickableCanvas.contentWidth = img.width
            flickableCanvas.contentHeight = img.height
        }
    }

    
}
