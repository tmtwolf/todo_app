package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/template"
	"todo_app/app/models"
	"todo_app/config"
)

func generateHTML(w http.ResponseWriter, data any, filenames ...string) {
	var files []string
	for _, filename := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", filename))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// セッションがあればセッションを返し、なければエラーを返す
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess := models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
		return sess, err
	}
	return
}

// URLの正規表現の型
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// _/todos/edit/1
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}

		qi, err := strconv.Atoi(q[2]) // q[2]、インデックス番号2はIDとして受け取る
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, qi)
	}
}

func StartMainServer() (err error) {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/user", user)
	http.HandleFunc("/user_pass", userPass)
	http.HandleFunc("/user_delete", userDelete)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))

	port := os.Getenv("PORT")
	log.Println("Port :" + port)
	if port == "" {
		port = "4000"
		log.Println("Fix_Port :" + port)
	}
	return http.ListenAndServe(":"+port, nil)
}
