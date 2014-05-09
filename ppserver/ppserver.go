package ppserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	// Install the QL SQL driver.
	_ "github.com/cznic/ql/driver"
	"github.com/nictuku/webpprof/ppcommon"
)

var (
	db   *sql.DB
	once sync.Once
)

// HandlePostProfile receives a pprof profile and stores it.
func HandlePostProfile(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		var err error
		db, err = sql.Open("ql", "ql.db")
		if err != nil {
			log.Fatalf("db opening error: %v", err)
		}
	})

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

// Schema:
//
// CREATE TABLE profiles (user string, profile blob, t time);

func saveProfile(p *Profile) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	result, err := tx.Exec(`
		INSERT INTO profiles VALUES ($1, $2, now());`,
		"yves.junqueira@gmail.com", profile.Content)

	if err != nil {
		return
	}
	if err := tx.Commit(); err != nil {
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

	log.Printf("decoded profile: %+q", profile)
	return
}
