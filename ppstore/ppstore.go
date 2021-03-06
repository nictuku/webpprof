package ppstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	// Install the QL SQL driver.
	_ "github.com/cznic/ql/driver"
	"github.com/nictuku/mothership/login"
	"github.com/nictuku/webpprof/ppcommon"
)

var (
	db   *sql.DB
	once sync.Once
)

func dbInit() {
	once.Do(func() {
		var err error
		db, err = sql.Open("ql", "ql.db")
		if err != nil {
			log.Fatalf("db opening error: %v", err)
		}
	})
}

// HandlePostProfile receives a pprof profile and stores it.
func HandlePostProfile(w http.ResponseWriter, r *http.Request) {
	dbInit()
	p := r.FormValue("p")
	if p == "" {
		http.NotFound(w, r)
		return
	}
	var profile ppcommon.Profile
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&profile); err != nil {
		log.Printf("handle post profile error parsing JSON for %v: %v", p, err)
		http.Error(w, "invalid profile content", http.StatusBadRequest)
		return
	}
	if err := saveProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleReadProfile shows a requested profile without authentication..
func HandleReadProfile(w http.ResponseWriter, r *http.Request) {
	readProfile(w, nil)
}

// HandleAuthReadProfile shows a requested profile after authentication.
func HandleAuthReadProfile(w http.ResponseWriter, r *http.Request) {
	passport, err := login.CurrentPassport(r)
	if err != nil {
		log.Printf("Redirecting to ghlogin: %v. Path: %q. Referrer: %q", err, r.URL.Path, r.Referer())
		http.SetCookie(w, &http.Cookie{Name: "ref", Value: r.URL.Path})
		http.Redirect(w, r, "/ghlogin", http.StatusFound)
		return
	}
	log.Println("login from user", passport.Email)
	readProfile(w, passport)
}

func readProfile(w io.Writer, passport *login.Passport) error {
	dbInit()
	log.Println("reading")
	rows, err := db.Query(`SELECT content, t FROM profiles WHERE name == "heap" && user == $1 ORDER BY t DESC LIMIT 1;`, passport.Email)
	if err != nil {
		fmt.Println(err)
		log.Fatalln("QUERY error:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var content string
		var t time.Time
		if err := rows.Scan(&content, &t); err != nil {
			log.Fatal(err)
		}
		raw, err := ppcommon.RawProfile(content)
		if err != nil {
			log.Println("RawProfile:", err)
		} else {
			content = raw
		}
		fmt.Fprintf(w, "%s\n", content)
		return nil
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Schema:
// CREATE TABLE profiles (user string, name string, content blob, t time);
func saveProfile(p *ppcommon.Profile) (err error) {
	// TODO: Move to Go's sql driver way of doing things.
	tx, err := db.Begin()
	if err != nil {
		return
	}
	result, err := tx.Exec(`
		INSERT INTO profiles VALUES ($1, $2, $3, now());`,
		"yves.junqueira@gmail.com", p.Name, p.Content)

	if err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	aff, err := result.RowsAffected()
	if err != nil {
		return
	}

	fmt.Printf("LastInsertId %d, RowsAffected %d\n", id, aff)

	log.Printf("decoded profile: %+q", p)
	return
}
