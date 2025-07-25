package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Utterance represents the structure of an utterance in the transcript.
// It includes the text, speaker, start time, and end time.
type Utterance struct {
	Text    string  `json:"text"`
	Speaker string  `json:"speaker"`
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
}

// CleanUtterance is a simplified version of Utterance for the final output.
// It only includes the text, start time, and end time.
type CleanUtterance struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// Global map to store transcriptions keyed by connection ID.
// This is used to retrieve transcriptions later.
var (
	transcriptions = make(map[string][]CleanUtterance)
	mu             sync.Mutex
)

// getUtterancesFromTranscript fetches the utterances from a completed transcript using the AssemblyAI API.
// It requires the API key and the transcript ID to make the request.
// It returns a slice of Utterance or an error if the request fails.
func getUtterancesFromTranscript(apiKey, transcriptID string) ([]Utterance, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Utterances []Utterance `json:"utterances"`
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Printf("Raw response: %s\n", string(bodyBytes))
		return nil, err
	}

	return data.Utterances, nil
}

// waitUntilCompleted polls the AssemblyAI API until the transcription is completed.
// It takes a client and a transcript ID as parameters.
// It returns the completed transcript or an error if the polling fails.
func waitUntilCompleted(client *assemblyai.Client, transcriptID string) (assemblyai.Transcript, error) {
	for {
		tr, err := client.Transcripts.Get(context.Background(), transcriptID)
		if err != nil {
			return tr, err
		}

		log.Println("Transcript polling status:", tr.Status)

		switch tr.Status {
		case assemblyai.TranscriptStatusCompleted:
			return tr, nil
		case assemblyai.TranscriptStatusError:
			return tr, fmt.Errorf("transcription failed: %s", *tr.Error)
		}

		time.Sleep(3 * time.Second)
	}
}

// upgrader is used to upgrade HTTP connections to WebSocket connections.
// It allows all origins for simplicity, but this should be restricted in production.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleWS handles incoming WebSocket connections.
// It reads binary audio data from the WebSocket, saves it to a temporary file,
// and sends it to AssemblyAI for transcription.
func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	connectionID := uuid.New().String()
	log.Println("New connection:", connectionID)

	mt, data, err := conn.ReadMessage()
	if err != nil || mt != websocket.BinaryMessage {
		log.Println("Failed to read binary audio:", err)
		return
	}

	tmpfile, err := os.CreateTemp("", "*.wav")
	if err != nil {
		log.Println("Temp file creation failed:", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		log.Println("Failed to write to temp file:", err)
		return
	}
	tmpfile.Close()

	audioFile, err := os.Open(tmpfile.Name())
	if err != nil {
		log.Println("Open audio file failed:", err)
		return
	}
	defer audioFile.Close()

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		log.Println("API key not found in environment")
		return
	}
	client := assemblyai.NewClient(apiKey)

	ctx := context.Background()
	params := &assemblyai.TranscriptOptionalParams{
		FormatText:    assemblyai.Bool(true),
		Punctuate:     assemblyai.Bool(true),
		SpeakerLabels: assemblyai.Bool(true),
	}

	transcript, err := client.Transcripts.TranscribeFromReader(ctx, audioFile, params)
	if err != nil {
		log.Println("Transcription failed:", err)
		return
	}

	completedTranscript, err := waitUntilCompleted(client, *transcript.ID)
	if err != nil {
		log.Println("Polling failed:", err)
		return
	}

	utterances, err := getUtterancesFromTranscript(apiKey, *completedTranscript.ID)
	if err != nil {
		log.Println("Failed to get utterances:", err)
		return
	}

	cleaned := make([]CleanUtterance, len(utterances))
	for i, u := range utterances {
		cleaned[i] = CleanUtterance{
			Text:  u.Text,
			Start: u.Start / 1000.0,
			End:   u.End / 1000.0,
		}
	}

	mu.Lock()
	transcriptions[connectionID] = cleaned
	mu.Unlock()

	conn.WriteJSON(map[string]string{"connection_id": connectionID})
}

// handleGetTranscription retrieves the transcription for a given connection ID.
// It responds with the transcription data in JSON format.
// If the transcription is not found, it returns a 404 error.
func handleGetTranscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	mu.Lock()
	data, ok := transcriptions[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "Transcription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWS)
	router.HandleFunc("/transcription/{id}", handleGetTranscription).Methods("GET")

	port := ":8080"
	fmt.Println("Server running on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
