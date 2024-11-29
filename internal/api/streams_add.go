package api

import (
	"encoding/json"
	"fmt"
	"github.com/AlexxIT/go2rtc/internal/streams"
	"net/http"
	"os"
	"sync"
)

var mu sync.Mutex
var streamsFile = "streams.yaml"

type Stream struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func InitStreamsAdd() {
	http.HandleFunc("/api/streams/add", addStreamHandler)
}

func addStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var stream Stream
	err := json.NewDecoder(r.Body).Decode(&stream)
	if err != nil {
		http.Error(w, "Erro ao analisar JSON", http.StatusBadRequest)
		return
	}

	if stream.Name == "" || stream.URL == "" {
		http.Error(w, "Nome e URL são obrigatórios", http.StatusBadRequest)
		return
	}

	// Adicionar a stream ao arquivo YAML
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(streamsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo de configuração", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s:\n  - url: %s\n", stream.Name, stream.URL))
	if err != nil {
		http.Error(w, "Erro ao escrever no arquivo de configuração", http.StatusInternalServerError)
		return
	}

	// Adicionar a stream ao cache de streams do programa
	streams.Add(stream.Name, stream.URL)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Stream adicionada com sucesso"))
}
