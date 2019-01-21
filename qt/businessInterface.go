package main

// controller "github.com/amlwwalker/wingit/packages/controller"

//this handles interfacing with any business logic occuring elsewhere
type BusinessInterface struct {
	iModel *FileModel //list of files for a user
	fModel *FileModel //list of files for a user
}

//handles the interface between the backend architecture
//and the bridge
func (b *BusinessInterface) configureInterface() {
	b.fModel = NewFileModel(nil)
	b.iModel = NewFileModel(nil)
}
