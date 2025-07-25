# Real-Time Transcription WebSocket Service  

A WebSocket-based transcription service using [AssemblyAI](https://www.assemblyai.com/), which receives audio in `.wav` format over WebSocket, transcribes it, and returns a `connection_id`. The full utterances can be retrieved using an HTTP GET request.

---

## Project Structure  

```
.
├── client.py                      # Python WebSocket client
├── main.go                        # Go server for WebSocket & HTTP
├── go.mod / go.sum                # Go modules
├── .env                           # Environment config (API keys)
├── requirements.txt               # Python dependencies
├── README.md                      # Project documentation
└── websocket_service_tester/      # Test WAV files
```

---

## Requirements  

### System Dependencies  

- Go 1.20+  
- Python 3.8+  
- Git  

### Go Dependencies (auto-installed via `go mod tidy`)  

- `github.com/AssemblyAI/assemblyai-go-sdk`  
- `github.com/cenkalti/backoff`  
- `github.com/coder/websocket`  
- `github.com/google/go-querystring`  
- `github.com/google/uuid`  
- `github.com/gorilla/mux`  
- `github.com/gorilla/websocket`  
- `github.com/joho/godotenv`  

### Python Dependencies  

Install using:  

```bash
pip install -r requirements.txt  
```

Content:  
```txt
websocket-client
```

---  

## Environment Setup  

Create a `.env` file in the root directory:  

```env
ASSEMBLYAI_API_KEY=your_assemblyai_api_key_here  
```

You can get your API key from [https://app.assemblyai.com](https://app.assemblyai.com)  

---  

## Running the Server  

Start the Go server:  

```bash
go run main.go  
```

Output:  

```
Server running on :8080  
```

---

## Running the Client

Basic usage:  

```bash
python client.py websocket_service_tester/8m_audio.wav  
```

Advanced usage:

```bash
python client.py <audio.wav> [--url ws://localhost:8080/ws] [--output txt|json|id]  
```

### Output Options  

- `--output id` -> Print only connection ID  
- `--output json` -> Print full JSON response (utterances)
- `--output txt` -> Print and save formatted text timestamps

Example result for `--output txt`:  

```
2.84 - 5.86: Hey Satya, I'm here and ready to dive in.  
5.86 - 7.84: Satya Nadella  
...
```

---  

## API Endpoints  

### 1. WebSocket (POST Binary Audio)  

**URL:** `ws://localhost:8080/ws`  

- Sends `.wav` audio binary  
- Returns:  
```json
{
  "connection_id": "your-uuid"  
}
```

---

### 2. HTTP GET Transcription  

**URL:** `http://localhost:8080/transcription/{connection_id}`  

- Example:  
```bash
curl http://localhost:8080/transcription/9d5b56ba-ff0c-413a-bf5c-1bdb3ce908de  
```

- Response:  
```json
[  
  {  
    "text": "Hey Satya, I'm here and ready to dive in.",  
    "start": 2.84,  
    "end": 5.86  
  },  
  ...  
]  
```

---  

## Notes  

- Make sure your `.wav` file is short and mono-channel for faster transcription.  
- The utterances endpoint will only be available after the transcription is **completed**.  
- WebSocket receives only one audio per connection (no streaming).  

---

## Troubleshooting

- If you receive `transcription not found`:  
  - Make sure you use the correct `connection_id`.  
  - Wait a few seconds, the transcription might still be processing.  
- If Go complains about missing packages, run:  
  ```bash
  go mod tidy  
  ```

---  

## License

This project is for educational and demonstration purposes. Please refer to AssemblyAI's terms for usage of their API.  