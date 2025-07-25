import websocket
import json
import argparse
import os
import sys
import requests

def format_txt_output(transcripts):
    lines = []
    for u in transcripts:
        start = round(u["start"], 2)
        end = round(u["end"], 2)
        text = u["text"]
        lines.append(f"{start} - {end}: {text}")
    return "\n".join(lines)

def main():
    parser = argparse.ArgumentParser(description="Send WAV file to WebSocket server and retrieve transcription.")
    parser.add_argument("filepath", help="Path to the .wav audio file")
    parser.add_argument("--url", default="ws://localhost:8080/ws", help="WebSocket server URL")
    parser.add_argument("--api", default="http://localhost:8080/transcription", help="HTTP API to get transcription by ID")
    parser.add_argument("--output", choices=["uuid", "json", "txt"], default="json", help="Output format")
    args = parser.parse_args()

    if not os.path.exists(args.filepath):
        print(f"Error: File '{args.filepath}' not found.")
        sys.exit(1)

    try:
        print(f"[WS] Connecting to {args.url} ...")
        ws = websocket.create_connection(args.url)

        with open(args.filepath, "rb") as f:
            audio_data = f.read()
            ws.send(audio_data, opcode=websocket.ABNF.OPCODE_BINARY)

        response = ws.recv()
        ws.close()

        ws_data = json.loads(response)
        connection_id = ws_data.get("connection_id")
        if not connection_id:
            print("Error: No connection_id returned.")
            return

        if args.output == "uuid":
            print(json.dumps({"connection_id": connection_id}, indent=2))
            return

        api_url = f"{args.api}/{connection_id}"
        print(f"[HTTP] Getting full transcript from {api_url}")
        resp = requests.get(api_url)
        if resp.status_code != 200:
            print("Error retrieving transcription:", resp.text)
            return

        utterances = resp.json()

        if args.output == "json":
            print("WS result:")
            print(json.dumps(ws_data, indent=2))

            print("\nHTTP result:")
            print(json.dumps(utterances, indent=2))

        elif args.output == "txt":
            text_output = format_txt_output(utterances)
            txt_filename = f"{connection_id}.txt"
            with open(txt_filename, "w", encoding="utf-8") as f:
                f.write(text_output)
            print(f"Transcript saved to {txt_filename}")

    except Exception as e:
        print("Error:", e)
        sys.exit(1)

if __name__ == "__main__":
    main()
