package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/edneyosf/gloged"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const appName = "Pet Path" // Nome da aplicação
const version = "0.1"            // Versão da aplicação
const port = "8000"              // Porta que será rodada a aplicação
const userDb = "admin"
const passwordDb = "123"
const urlDb = "localhost:3306"
const nameDb = "pet_path"

var database *sql.DB

func main() {
	connectToDb()
	log.DebugMode = true
	router := mux.NewRouter()

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")

	log.I(appName + " v" + version)
	log.S("API iniciada na porta " + port)
	log.E(http.ListenAndServe(":"+port, router).Error())
}

/* REQUEST */

// Para efetuar o login pelo app
func login(w http.ResponseWriter, r *http.Request) {
	const method = "login"
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Em(method, err.Error())
		return
	} else {
		valid := areUserPasswordValid(user)
		response := Response{false, "Usuário ou senha inválidos"}

		log.D(method, "user", fmt.Sprintf("%v", user))

		if valid {
			response.Status = true
			response.Message = "Logado com sucesso"
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Para efetuar o registro pelo app
func register(w http.ResponseWriter, r *http.Request) {
	const method = "register"
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Em(method, err.Error())
		return
	} else {
		exists, _ := userExists(user)
		response := Response{false, "Usuário já registrado"}

		log.D(method, "user", fmt.Sprintf("%v", user))

		if !exists {
			if addUser(user) {
				response.Status = true
				response.Message = "Registrado com sucesso"
			} else {
				response.Message = "Falha ao registrar usuário"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

/* DATABASE */

func connectToDb() {
	var err error = nil

	database, err = sql.Open("mysql", userDb+":"+passwordDb+"@tcp("+urlDb+")/"+nameDb)

	if err != nil {
		panic(err.Error())
	}
}

// Registra um novo usuário
func addUser(user User) bool {
	status := true

	_, err := database.Exec("INSERT INTO user (name_user, password) VALUES ('" + user.This + "', '" + user.Password + "')")

	if err != nil {
		status = false
		panic(err.Error())
	}

	return status
}

// Valida o login de um usuário
func areUserPasswordValid(user User) bool {
	status := false
	id := -1

	results, err := database.Query("SELECT id_user FROM user WHERE name_user = '" + user.This + "' AND password = '" + user.Password + "' LIMIT 1")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}

		if id != -1 {
			status = true
			break
		}
	}

	return status
}

// Verifica a existência de um usuário
func userExists(user User) (bool, int) {
	status := false
	id := -1

	results, err := database.Query("SELECT id_user FROM user WHERE name_user = '" + user.This + "' LIMIT 1")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}

		if id != -1 {
			status = true
			break
		}
	}

	return status, id
}

/* STRUCT */

type User struct {
	This     string `json:"user"`
	Password string `json:"password"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
