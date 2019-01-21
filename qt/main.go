package main

import (
	"os"

	//for compiled files
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
)

type Config struct {
}

func main() {

	core.QCoreApplication_SetOrganizationName("OrganizationName")                 //needed to fix an QML Settings issue on windows
	quick.QQuickWindow_SetSceneGraphBackend(quick.QSGRendererInterface__Software) //needed to get the application working on VMs when using the windows docker images

	//0. set any required env vars for qt
	// os.Setenv("QT_QUICK_CONTROLS_STYLE", "material") //set style to material
	os.Setenv("QML_DISABLE_DISK_CACHE", "true") //disable caching files

	var config Config
	//3. Create a bridge to the frontend
	var qmlBridge = NewQmlBridge(nil)
	qmlBridge.ConfigureBridge(config)
	// turn on high definition scaling
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	// quick.QQuickWindow_SetSceneGraphBackend(quick.QSGRendererInterface__Software)
	//4. Configure the qml binding and create an application
	widgets.NewQApplication(len(os.Args), os.Args)

	//create a view
	var view = quick.NewQQuickView(nil)
	view.SetTitle("Fitz Test")

	//configure the view to know about the bridge
	//this needs to happen before anything happens on another thread
	//else the thread might beat the context property to setup

	view.RootContext().SetContextProperty("QmlBridge", qmlBridge)
	view.RootContext().SetContextProperty("FilesModel", qmlBridge.business.fModel)
	view.RootContext().SetContextProperty("ImageFilesModel", qmlBridge.business.iModel)

	view.SetSource(core.NewQUrl3("qrc:/qml/loader-production.qml", 0))

	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.Show()
	widgets.QApplication_Exec()

}
