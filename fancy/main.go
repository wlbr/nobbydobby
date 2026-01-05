package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Guest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func (g *Guest) String() string {
	return fmt.Sprintf(`"%s","%s","%s"`, g.Firstname, g.Lastname, g.Email)
}

type felixHandler struct {
	address        string
	fulladdress    string
	handlers       map[string]func(http.ResponseWriter, *http.Request)
	mux            *http.ServeMux
	guests         map[string]*Guest
	guestsfilename string
	m              sync.RWMutex
}

func NewFelixHandler(address string) *felixHandler {
	h := &felixHandler{address: address, mux: http.NewServeMux(), guests: make(map[string]*Guest), guestsfilename: "guests.db"}
	h.fulladdress = address
	if address[0] == ':' {
		h.fulladdress = "localhost" + address
	}
	h.handlers = make(map[string]func(http.ResponseWriter, *http.Request))
	return h

}

func (h *felixHandler) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.mux.HandleFunc(pattern, handler)
	h.handlers[pattern] = handler
}

func (h *felixHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<h1>Routes:</h1>\n<ul>"))
	for pattern, _ := range h.handlers {
		pats := strings.Split(pattern, " ")
		if len(pats) == 1 {
			w.Write([]byte(`<li>ALL<a href="http://` + h.fulladdress + pats[0] + `"> ` + pats[0] + "</a></li> \n"))
		} else if pats[0] == "GET" {
			w.Write([]byte("<li> " + pats[0] + ` <a href="http://` + h.fulladdress + pats[1] + `">` + pats[1] + "</a></li> \n"))
		} else {
			w.Write([]byte("<li>" + pattern + "</li> \n"))
		}
	}
	w.Write([]byte("</ul>"))
}

func (h *felixHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got request from %s: %s %s", r.RemoteAddr, r.Method, r.URL)
	h.mux.ServeHTTP(w, r)
}

// ------------------------------

func (h *felixHandler) ReadGuests(fname string) error {
	h.m.Lock()
	defer h.m.Unlock()
	f, err := os.Open(fname)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Printf("Error during read Guests from file: %s", err)
		return fmt.Errorf("Error during read Guests from file: %s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		j := scanner.Bytes()
		g := &Guest{}
		if err := json.Unmarshal(j, g); err != nil {
			log.Printf("Error unmarshalling guest from file: %s", err)
		} else {
			h.guests[g.Email] = g
		}
	}
	return nil
}

func (h *felixHandler) AddGuest(g *Guest) bool {
	h.m.Lock()
	defer h.m.Unlock()
	if _, ok := h.guests[g.Email]; ok {
		log.Printf("Email %s already on guestlist, skipped", g.Email)
		return false
	}
	h.guests[g.Email] = g

	// Append the new guest to the file
	f, err := os.OpenFile(h.guestsfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening guests file for append: %s", err)
		return false // Or handle error appropriately
	}
	defer f.Close()

	var j []byte
	if j, err = json.Marshal(g); err != nil {
		log.Printf("Error marshalling new guest for append: %s", err)
		return false
	}
	if _, err = f.Write(j); err != nil {
		log.Printf("Error writing new guest to file: %s", err)
		return false
	}
	if _, err = f.WriteString("\n"); err != nil {
		log.Printf("Error writing newline after new guest to file: %s", err)
		return false
	}

	return true
}

// ------------------------------

func (h *felixHandler) handleReadGuests(w http.ResponseWriter, r *http.Request) {
	h.m.RLock()
	defer h.m.RUnlock()

	w.Header().Add("Content-Type", "text/plain")
	//w.Header().Add("Content-Type", "text/csv")
	for _, g := range h.guests {
		fmt.Fprintln(w, g)
	}
}

func (h *felixHandler) handleRegisterGuest(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	g := &Guest{}
	if err := json.Unmarshal(b, g); err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.AddGuest(g) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusAlreadyReported)
	}
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	h := NewFelixHandler(":8080")
	if err := h.ReadGuests(h.guestsfilename); err != nil {
		log.Fatalf("Error reading guests database.")
	}

	h.HandleFunc("GET /all", h.handleReadGuests)
	h.HandleFunc("POST /register", h.handleRegisterGuest)
	h.HandleFunc("/", h.handleRoot)

	log.Printf("Starting server at http://%s \n", h.fulladdress)
	if err := http.ListenAndServe(h.address, h); err != nil {
		fmt.Println(err)
	}
}
