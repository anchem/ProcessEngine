// UserControllHandler
package webclient

import (
	pe "ProcessEngine/src/processengine"
	"log"
	"net/http"
	"strings"
)

func UserControllHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cmd := r.FormValue("cmd")
	userId, _ := r.Cookie(COOKIE_USER_ID)
	log.Println(userId.Value)
	switch {
	//	case strings.EqualFold(strings.TrimSpace(cmd), "uploadFile"):
	//		fName := r.FormValue("fileName")
	//		fName = strings.TrimSpace(fName)
	//		if !strings.EqualFold(fName, "") {
	//			log.Println(fName)
	//			userFile[fName] = userId.Value
	//		}
	case strings.EqualFold(strings.TrimSpace(cmd), "query"):
		result := queryProcessDataByUser(userId.Value)
		w.Write([]byte(result))
	case strings.EqualFold(strings.TrimSpace(cmd), "delete"):
		procKey := r.FormValue("procKey")
		if strings.EqualFold(strings.TrimSpace(procKey), "") {
			// can not find key
			result := "{\"statusCode\":1,\"reason\":\"process key can not be null\"}"
			w.Write([]byte(result))
		} else {
			log.Println("yes!")
			result := "{\"statusCode\":0}"
			w.Write([]byte(result))
		}
	}
}
func handleError(err error, str string, w *http.ResponseWriter) {

}
func queryProcessDataByUser(userId string) string {
	resultStr, flag := pe.GetProcessEngine().QueryProcessDefByUser(userId)
	var result string
	if flag == pe.QUERY_SUCCESS {
		result = "{\"statusCode\":0,\"list\":" + resultStr + "}"
	} else {
		result = "{\"statusCode\":1,\"reason\":\"" + resultStr + "\"}"
	}
	return result
}
