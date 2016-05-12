// WebApp
package webclient

import (
	pe "ProcessEngine/src/processengine"
	"file"
	"fmt"
	"net/http"
)

func ServeHttp() {
	//	http.Handle("/css/", http.FileServer(http.Dir("./webpages/css")))
	//	http.Handle("/js/", http.FileServer(http.Dir("./webpages/js")))
	http.HandleFunc("/", LoginHandler)
	http.HandleFunc("/pannel", UserControllHandler)
	go receiveFile()
	http.ListenAndServe(":8088", nil)
}
func receiveFile() {
	fileRcver := file.NewFileReceiver()
	go handleFile(fileRcver.CPLch)
	fileRcver.Receive()
}

func handleFile(nameCh chan *file.UserFile) {
	var uf *file.UserFile
	for {
		uf = <-nameCh
		fmt.Println("--complete--", uf.UserId, ":", uf.FileName)
		result := pe.GetProcessEngine().DeployWithUser(uf.FileName, pe.DEPLOY_FILE_XML, pe.DEPLOY_MODE_DEFAULT, uf.UserId, 0)
		var str string
		switch result {
		case pe.DEPLOY_SUCCESS:
			str = "depoly successfully"
		case pe.DEPLOY_FAILED_ERROR:
			str = "depoly failed"
		case pe.DEPLOY_FAILED_REPITITION:
			str = "depoly repitition"
		}
		fmt.Println(str)
	}
}
