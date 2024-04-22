package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Home()))
}

func StartAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/interaction", getInteraction).Methods("GET")
	r.HandleFunc("/search", handleSearch).Methods("GET")
	r.HandleFunc("/drugs", drugsTest).Methods("GET")
	r.Handle("/", &homeHandler{}).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("s")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}
	fmt.Printf("Request Recieved: %s\n", query)
	results, _ := searchFiles(strings.ToLower(query))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func searchFiles(searchTerm string) ([]string, error) {
	var foundFiles []string
	files, err := ioutil.ReadDir("./cache/interactions")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), searchTerm) {
			foundFiles = append(foundFiles, strings.Replace(strings.Replace(file.Name(), ".json", "", -1), "-interactions-", "", -1))
		}
	}

	return foundFiles, nil
}

func getInteraction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("i")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}
	query = fmt.Sprintf("./cache/interactions/-interactions-%s.json", query)
	content, err := ioutil.ReadFile(query)
	if err != nil {
		http.Error(w, "Invalid Interaction", http.StatusBadRequest)
		return
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonData)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(Interactions(jsonData)))
	}

}

func drugsTest(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("/Users/khushchauhan/Desktop/FYP/CODE/scraper/drugstest.txt")
	if err != nil {
		http.Error(w, "Invalid Interaction "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}
