package ppserver

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/profile", HandlePostProfile)

}

// profile contains the metadata and the raw content of a pprof profile.
type Profile struct {
	Name    string
	Content string
	Time    time.Time
}

func HandlePostProfile(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	p := r.FormValue("p")
	if p == "" {
		http.NotFound(w, r)
		return
	}
	var profile Profile
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&profile)
	if err != nil {
		log.Printf("handle post profile error parsing JSON for %v: %v", p, err)
		http.Error(w, "invalid profile content", http.StatusBadRequest)
		return
	}

	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.String()))
	/*	body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("handleProfile read body error: %v", err)
			return
		}
	*/
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Profile", nil), &profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("decoded profile: %+q", profile)
	return
}
