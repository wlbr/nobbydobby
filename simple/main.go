package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Guest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

type felixHandler struct {
	mux            *http.ServeMux
	guests         map[string]*Guest
	guestsfilename string
	m              sync.RWMutex
}

func NewFelixHandler() *felixHandler {
	return &felixHandler{mux: http.NewServeMux(), guests: make(map[string]*Guest), guestsfilename: "guests.db"}
}

func (h *felixHandler) ReadGuests(fname string) {
	h.m.Lock()
	defer h.m.Unlock()
	f, err := os.Open(fname)
	if err != nil {
		log.Printf("Error during read Guests from file: %s", err)
		return
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

func (h *felixHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got request from %s: %s %s", r.RemoteAddr, r.Method, r.URL)
	h.mux.ServeHTTP(w, r)
}

func (h *felixHandler) handleReadGuests(w http.ResponseWriter, r *http.Request) {
	h.m.RLock()
	defer h.m.RUnlock()

	w.Header().Add("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	// Write CSV header (optional but good practice)
	if err := csvWriter.Write([]string{"Firstname", "Lastname", "Email"}); err != nil {
		log.Printf("Error writing CSV header: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, g := range h.guests {
		record := []string{g.Firstname, g.Lastname, g.Email}
		if err := csvWriter.Write(record); err != nil {
			log.Printf("Error writing guest CSV record: %s", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	h := NewFelixHandler()
	h.ReadGuests(h.guestsfilename)

	h.mux.HandleFunc("GET /all", h.handleReadGuests)
	h.mux.HandleFunc("POST /register", h.handleRegisterGuest)

	http.ListenAndServe(":8080", h)
}
