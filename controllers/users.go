package controllers

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var user user

	if error = json.Unmarshal(request, &user); error != nil {
		w.Write([]byte("Erro ao converter o usuário para struct"))
		return
	}

	db, error := database.ConnectDatabase()
	if error != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
		return
	}

	defer db.Close()

	statement, error := db.Prepare("insert into usuarios (nome, email) values (?, ?)")

	if error != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statement.Close()

	insert, error := statement.Exec(user.Nome, user.Email)

	if error != nil {
		w.Write([]byte("Erro ao executar o statement"))
		return
	}

	idInserted, error := insert.LastInsertId()

	if error != nil {
		w.Write([]byte("Erro ao inserir o dado"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! ID: %d", idInserted)))

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, error := database.ConnectDatabase()

	if error != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
	}

	defer db.Close()

	data, error := db.Query("select * from usuarios")

	var users []user

	for data.Next() {
		var user user

		if error := data.Scan(&user.ID, &user.Nome, &user.Email); error != nil {
			w.Write([]byte("Erro ao escanear usuários"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)

	if error := json.NewEncoder(w).Encode(users); error != nil {
		w.Write([]byte("Erro ao converter os usuários para JSON"))
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, error := strconv.ParseUint(params["id"], 10, 32)

	if error != nil {
		w.Write([]byte("Erro ao converter id para int"))
		return
	}

	db, error := database.ConnectDatabase()

	if error != nil {
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}

	data, error := db.Query("select * from usuarios where id = ?", ID)
	if error != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}

	var user user

	if data.Next() {
		if error := data.Scan(&user.ID, &user.Nome, &user.Email); error != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(w).Encode(user); error != nil {
		w.Write([]byte("Erro ao converter usuário para json"))
		return
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, error := strconv.ParseUint(params["id"], 10, 32)
	if error != nil {
		w.Write([]byte("Erro ao parsear o id"))
		return
	}

	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição"))
		return
	}

	var user user
	if error := json.Unmarshal(body, &user); error != nil {
		w.Write([]byte("Erro ao converter usuário para struct"))
		return
	}

	db, error := database.ConnectDatabase()
	if error != nil {
		w.Write([]byte("Erro ao conectar com um banco"))
		return
	}

	defer db.Close()

	statement, error := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if error != nil {
		w.Write([]byte("Erro ao gerar o statement"))
		return
	}

	defer statement.Close()

	if _, error := statement.Exec(user.Nome, user.Email, ID); error != nil {
		w.Write([]byte("Erro ao atualizar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, error := strconv.ParseUint(params["id"], 10, 32)
	if error != nil {
		w.Write([]byte("Erro ao parsear id"))
		return
	}

	db, error := database.ConnectDatabase()
	if error != nil {
		w.Write([]byte("Erro ao conectar ao banco"))
		return
	}
	defer db.Close()

	statement, error := db.Prepare("delete from usuarios where id = ?")
	if error != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}
	defer statement.Close()

	if _, error := statement.Exec(ID); error != nil {
		w.Write([]byte("Erro ao deletar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
