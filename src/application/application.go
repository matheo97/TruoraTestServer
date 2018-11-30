package application

import (
	// native packages
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	// local packages
	"recipes"
	// GitHub packages
	"github.com/gorilla/mux"
	// Standard SQL Override
	_ "github.com/lib/pq"
)

// App represents the application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) getRecipesEndpoint(w http.ResponseWriter, req *http.Request) {

  	recipes, err := recipes.GetRecipes(a.DB)
  	if err != nil {
  		respondWithError(w, http.StatusInternalServerError, err.Error())
  		return
  	}
  	respondWithJSON(w, http.StatusOK, recipes)
}

func (a *App) createRecipeEndpoint(w http.ResponseWriter, req *http.Request) {
  	var r recipes.Recipe
  	decoder := json.NewDecoder(req.Body)
  	if err := decoder.Decode(&r); err != nil {
  		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
  		return
  	}
  	defer req.Body.Close()
  	if err := r.CreateRecipe(a.DB); err != nil {
  		respondWithError(w, http.StatusInternalServerError, err.Error())
  		return
  	}
  	respondWithJSON(w, http.StatusCreated, r)
}

func (a *App) modifyRecipeEndpoint(w http.ResponseWriter, req *http.Request) {
  	params := mux.Vars(req)
  	id, err := strconv.Atoi(params["id"])
  	if err != nil {
  		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
  		return
  	}
  	var r recipes.Recipe
  	decoder := json.NewDecoder(req.Body)
  	if err := decoder.Decode(&r); err != nil {
  		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
  		return
  	}
  	defer req.Body.Close()
  	r.ID = id
  	if _, err := r.UpdateRecipe(a.DB); err != nil {
  		respondWithError(w, http.StatusInternalServerError, err.Error())
  		return
  	}
  	respondWithJSON(w, http.StatusOK, r)
}

func (a *App) deleteRecipeEndpoint(w http.ResponseWriter, req *http.Request) {
  	params := mux.Vars(req)
  	id, err := strconv.Atoi(params["id"])
  	if err != nil {
  		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
  		return
  	}
  	r := recipes.Recipe{ID: id}
  	if _, err := r.DeleteRecipe(a.DB); err != nil {
  		respondWithError(w, http.StatusInternalServerError, err.Error())
  		return
  	}
  	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) searchRecipesEndpoint(w http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
  	name, ok := params["name"]
  	if ok != true {
  		respondWithError(w, http.StatusBadRequest, "Invalid recipe name")
  		return
  	}

  	recipes, err := recipes.GetRecipesByName(a.DB, name)
  	if err != nil {
  		respondWithError(w, http.StatusInternalServerError, err.Error())
  		return
  	}
  	respondWithJSON(w, http.StatusOK, recipes)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

// Initialize sets up the database connection, router, and routes for the app
func (a *App) Initialize(user, dbname string) {

		connectionString := fmt.Sprintf("postgresql://%s@localhost:26257/%s?sslmode=disable", user, dbname)

  	var err error
		fmt.Println(connectionString)
  	a.DB, err = sql.Open("postgres", connectionString)
  	if err != nil {
  		log.Fatal(err)
  	}

  	a.Router = mux.NewRouter()

		//Allow CORS
		a.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  		w.Header().Set("Access-Control-Allow-Origin", "*")
  		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
  		w.WriteHeader(http.StatusNoContent)
  		return
  	})
		a.Router.StrictSlash(true)

  	v1 := a.Router.PathPrefix("/v1").Subrouter()

  	v1.HandleFunc("/recipes", a.getRecipesEndpoint).Methods("GET")
  	v1.HandleFunc("/recipes", a.createRecipeEndpoint).Methods("POST")
  	v1.HandleFunc("/editRecipe/{id:[0-9]+}", a.modifyRecipeEndpoint).Methods("POST")
  	v1.HandleFunc("/recipes/{id:[0-9]+}", a.deleteRecipeEndpoint).Methods("DELETE")
  	v1.HandleFunc("/recipes/{name:[a-zA-Z\\s]+}", a.searchRecipesEndpoint).Methods("GET")
}

// Run starts the app and serves on the specified port
func (a *App) Run(port string) {
	log.Print("Now serving recipes ...")
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}
