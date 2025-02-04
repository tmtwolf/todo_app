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

// 以下Todo
func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}
	todos, err := user.GetTodosByUser()
	user.Todos = todos
	generateHTML(w, user, "layout", "private_navbar", "index")

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
		log.Println("セッション認証失敗")
		http.Redirect(w, r, "/login", 302)
	} else {
		// フォームの値の取得、パースフォーム
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		// アップデート先Todoの取得
		todo, err := models.GetTodoById(id)
		if err != nil {
			log.Println(err)
		}

		// _Todoの更新
		todo.Content = r.PostFormValue("content")
		err = todo.UpdateTodo()
		if err != nil {
			log.Println(err)
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

// 以下ユーザー
func user(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	// テンプレート用のデータ型を作成
	var data models.HttpData
	data.User = user

	// 各処理へ
	if r.Method == "GET" {
		generateHTML(w, data, "layout", "private_navbar", "user")

	} else if r.Method == "POST" {
		// POST内容の取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			data.Msg = "不正なリクエスト"
			generateHTML(w, data, "layout", "private_navbar", "user")
		}

		// ユーザーオブジェクトの更新
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")

		// テンプレート用のデータも更新
		data.User = user

		// データベースのアップデート
		if err = user.UpdateUser(); err == nil {
			data.Msg = "更新に成功しました"
			generateHTML(w, data, "layout", "private_navbar", "user")
		} else {
			log.Println(err)
			data.Msg = "更新に失敗しました"
			generateHTML(w, data, "layout", "private_navbar", "user")
		}
	}
}

func userPass(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	// テンプレート用のデータ型を作成
	var data models.HttpData
	data.User = user

	// 各処理へ
	if r.Method == "GET" {
		generateHTML(w, data, "layout", "private_navbar", "user_pass")

	} else if r.Method == "POST" {
		// POST内容の取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			data.Msg = "不正なリクエスト"
		} else {
			u, err := models.GetUserByEmail(sess.Email)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login", 302)
			}

			oldPass := r.PostFormValue("old_password")
			newPass := r.PostFormValue("new_password")

			if u.Password == models.Encrypt(oldPass) && oldPass != newPass {
				user.Password = models.Encrypt(newPass)
				if err = user.UpdateUserPass(); err == nil {
					data.Msg = "パスワードを更新しました。"
				} else {
					log.Println(err)
					data.Msg = "パスワードの更新に失敗しました。"
				}
			} else {
				data.Msg = "パスワードの更新に失敗しました。"
			}
		}
		generateHTML(w, data, "layout", "private_navbar", "user_pass")
	}
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	// テンプレート用のデータ型を作成
	var data models.HttpData
	data.User = user

	// 各処理へ
	if r.Method == "GET" {
		generateHTML(w, data, "layout", "private_navbar", "user_delete")

	} else if r.Method == "POST" {
		// POST内容の取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			data.Msg = "不正なリクエスト"
		} else {
			u, err := models.GetUserByEmail(sess.Email)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login", 302)
			}

			if u.Password == models.Encrypt(r.PostFormValue("password")) {
				if err = user.DeleteUser(); err == nil {
					generateHTML(w, data, "layout", "public_navbar", "user_deleted")
					return
				} else {
					data.Msg = "ユーザーの削除に失敗しました。"
				}
			} else {
				data.Msg = "ユーザーの削除に失敗しました。"
			}
			generateHTML(w, data, "layout", "private_navbar", "user_delete")
		}
	}
}
