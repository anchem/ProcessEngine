// LoginHandler
package webclient

import (
	pe "ProcessEngine/src/processengine"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	_ = iota
	USER_AUTH
	USER_DENY
)
const (
	COOKIE_USER_ID = "userId"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		username := r.FormValue("username")
		password := r.FormValue("password")
		userId, result := checkUser(username, password)
		switch result {
		case USER_AUTH:
			// login success
			expiration := time.Now()
			expiration = expiration.AddDate(1, 0, 0)
			cookie := http.Cookie{Name: COOKIE_USER_ID, Value: userId, Expires: expiration}
			http.SetCookie(w, &cookie)
			cur, err := os.Getwd()
			if err != nil {
				log.Panicln(err)
			}
			path := cur[:strings.LastIndex(cur, "\\")]
			t, err := template.ParseFiles(path + "/webpages/html/user_pannel.html")
			if err != nil {
				log.Panicln(err)
			}
			t.Execute(w, nil)
		case USER_DENY:
			w.Write([]byte("user denied"))
		default:
		}
	} else {
		cur, err := os.Getwd()
		if err != nil {
			log.Panicln(err)
		}
		path := cur[:strings.LastIndex(cur, "\\")]
		t, err := template.ParseFiles(path + "/webpages/html/index.html")
		if err != nil {
			log.Panicln(err)
		}
		t.Execute(w, nil)
	}

}
func checkUser(uname string, pwd string) (string, byte) {
	if strings.EqualFold(uname, "") || strings.EqualFold(pwd, "") {
		return "", USER_DENY
	}
	result, userId := pe.GetProcessEngine().CheckUser(uname, pwd)
	if result == 0 {
		return userId, USER_AUTH
	} else {
		return "", USER_DENY
	}
}
