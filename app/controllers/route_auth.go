package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		if _, err = models.GetUserByEmail(email); err == nil {
			http.Redirect(w, r, "/login", 302)
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: password,
		}

		if err := user.CreateUser(); err != nil {
			log.Fatalln(err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}

	if user.Password == models.Encrypt(r.PostFormValue("password")) {
		sess, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    sess.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/top", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		sess := models.Session{
			UUID: cookie.Value,
		}

		if err = sess.DeleteSessionByUUID(); err != nil {
			log.Println("セッションの削除に失敗")
		}

		http.Redirect(w, r, "/login", 302)
	}
}
