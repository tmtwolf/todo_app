package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}

	/*
		t, err := template.ParseFiles("app/views/templates/top.html")
		if err != nil {
			log.Fatalln(err)
		}
		t.Execute(w, "hello")
	*/
}

func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		todos, _ := user.GetTodosByUser()
		user.Todos = todos
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

func todoSave(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else if r.Method == "POST" {
		// フォームの値の取得、パースフォーム
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// 保存先のユーザーの取得
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Panicln(err)
		}

		// _Todoの作成
		err = user.CreateTodo(r.PostFormValue("content"))
		if err != nil {
			log.Panicln(err)
		}
		http.Redirect(w, r, "/index", 302)
	}
}

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodoById(id)
		if err != nil {
			log.Println(err)
		}
		generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else if r.Method == "POST" {
		// フォームの値の取得、パースフォーム
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// アップデート先Todoの取得
		todo, err := models.GetTodoById(id)
		if err != nil {
			log.Panicln(err)
		}

		// _Todoの更新
		todo.Content = r.PostFormValue("content")
		err = todo.UpdateTodo()
		if err != nil {
			log.Panicln(err)
		}
		http.Redirect(w, r, "/index", 302)
	}
}

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// フォームの値の取得、パースフォーム
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// 削除先Todoの取得
		todo, err := models.GetTodoById(id)
		if err != nil {
			log.Println(err)
		}

		// _Todoの削除
		err = todo.DeleteTodo()
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/index", 302)
	}
}
