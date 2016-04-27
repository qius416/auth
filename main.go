package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	auth "github.com/qius416/auth/authentication"
	"golang.org/x/crypto/bcrypt"
	db "gopkg.in/dancannon/gorethink.v2"
)

// Session for access rethinkdb
var Session *db.Session

func init() {
	var err error
	for Session == nil {
		time.Sleep(time.Second * 5)
		Session, err = db.Connect(db.ConnectOpts{
			Address: "db:28015",
			MaxIdle: 1,
			MaxOpen: 2,
		})
		// test purpose for clustered db
		// session, err = db.Connect(db.ConnectOpts{
		// 	Addresses:     []string{"db1:28015", "db2:28015"},
		// 	DiscoverHosts: true,
		// })
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	// ignore error for create
	db.DBCreate("bourbaki").Run(Session)
	Session.Use("bourbaki")
	db.TableCreate("user", db.TableCreateOpts{PrimaryKey: "email"}).Run(Session)
}

func main() {
	router := httprouter.New()
	router.POST("/signup", signup)
	router.POST("/login", login)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", router))
}

func signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var u auth.User
	err := decoder.Decode(&u)
	if err != nil {
		log.Fatalln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}
		u.Password = string(hash)
		_, dberr := db.Table("user").Insert(u, db.InsertOpts{Conflict: "error"}).RunWrite(Session)
		if dberr != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "user %s exists.", u.Email)
			return
		}
		fmt.Fprintf(w, "user created.")
	}
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var u auth.User
	err := decoder.Decode(&u)
	if err != nil {
		log.Fatalln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	res, dberr := db.Table("user").Get(u.Email).Run(Session)
	defer res.Close()
	if dberr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, dberr.Error())
		return
	}

	if res.IsNil() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	var myuser auth.User
	err = res.One(&myuser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	hasherr := bcrypt.CompareHashAndPassword([]byte(myuser.Password), []byte(u.Password))

	if hasherr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "access rejected.")
		return
	}
	token, tokeErr := auth.MakeToken(myuser.Name, myuser.Role)
	if tokeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, tokeErr.Error())
	} else {
		p := auth.Auth{Token: token}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}
	return
}
