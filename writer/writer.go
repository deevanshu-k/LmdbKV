package writer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deevanshu-k/lmdbkv/store"
)

var log_key = "WRITER-SERVER"

func StartHttpWriterServer(s *store.Store, address string) error {
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		key, value := body["key"], body["value"]
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}
		if err := s.Set([]byte(key), []byte(value)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Printf("[%s] Listening on %s", log_key, address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("[%s] Error from listner: %v", log_key, err)
	}

	return nil
}
