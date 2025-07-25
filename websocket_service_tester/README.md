# ASR Service Tester

A WebSocket service tester for Automatic Speech Recognition (ASR) services that simulates real-time audio streaming and captures transcription results.

## Description

This script converts audio files to PCM format and streams them in real-time to a WebSocket-based ASR service, simulating live audio input. It's designed to test the performance and accuracy of speech recognition services by sending audio data with proper timing intervals and collecting transcription results.

## Features

- ✅ Real-time audio streaming simulation
- ✅ Audio format conversion (any format → PCM S16LE, 16kHz, mono)
- ✅ WebSocket connection management
- ✅ Transcription result logging
- ✅ Configurable connection parameters
- ✅ Performance metrics tracking

## Requirements

- **Node.js**: Version 22 or higher
- **ffmpeg**: For audio format conversion
- **TypeScript**: For compilation (included as dev dependency)

### Installing Requirements

#### Node.js 22
```bash
# Using nvm (recommended)
nvm install 22
nvm use 22

# Or download from https://nodejs.org/
```

#### FFmpeg
```bash
# macOS (using Homebrew)
brew install ffmpeg

# Ubuntu/Debian
sudo apt update
sudo apt install ffmpeg

# Windows (using Chocolatey)
choco install ffmpeg

# Or download from https://ffmpeg.org/download.html
```

## Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd asr_service_tester
   ```

2. **Install dependencies**:
   ```bash
   yarn install
   # or
   npm install
   ```

3. **Build the project**:
   ```bash
   yarn build
   # or
   npm run build
   ```

## Usage

### Basic Usage

```bash
# Using npm scripts (recommended)
yarn build && yarn start --wsUrl=<websocket_url> --audioPath=<path_to_audio_file>

# Or directly
node dist/src/index.js --wsUrl=<websocket_url> --audioPath=<path_to_audio_file>
```

### Parameters

- `--wsUrl`: **Required** - WebSocket URL of the ASR service
- `--audioPath`: **Required** - Path to the audio file to be processed
- `--connectionId`: **Optional** - Custom connection ID (auto-generated if not provided)

### Examples

#### Example 1: Basic ASR Testing
```bash
node dist/src/index.js \
  --wsUrl=ws://localhost:8080/asr \
  --audioPath=/Users/salim/repos/asr_service_tester/8m_audio.wav
```

#### Example 2: With Custom Connection ID
```bash
node dist/src/index.js \
  --wsUrl=wss://your-asr-service.com/transcribe \
  --audioPath=./8m_audio.wav \
  --connectionId=test-session-001
```

## Output

The script will:
1. Connect to the specified WebSocket URL
2. Convert the audio file to PCM format (16kHz, mono, 16-bit)
3. Stream audio data in real-time chunks
4. Display transcription results in the console
5. Filter and process only messages containing transcription text
6. Save results to `temp/<session-id>.txt` with timestamps

### Sample Output
```
websocket is open
Sent chunk of size: 1024 0.032 1024
Sent chunk of size: 1024 0.064 2048
{
  "start": 0.50,
  "end": 1.20,
  "text": "Hello, this is a test"
}
{
  "start": 1.20,
  "end": 2.80,
  "text": "of the speech recognition service"
}
Conversion completed
websocket is closed
```

## Project Structure

```
asr_service_tester/
├── src/
│   └── index.ts          # Main script
├── dist/                 # Compiled JavaScript output
├── temp/                 # Transcription output directory (auto-created)
├── 8m_audio.wav         # Sample audio file
├── example_result.txt   # Example result file
├── package.json         # Dependencies and scripts
├── tsconfig.json        # TypeScript configuration
├── yarn.lock           # Dependency lock file
└── README.md           # This file
```

## Troubleshooting

### Common Issues

1. **FFmpeg not found**:
   ```
   Error: FFmpeg is not installed. Please install it first.
   ```
   Solution: Install FFmpeg using the instructions above.

2. **Audio file not found**:
   ```
   Error: Input file does not exist
   ```
   Solution: Check the file path and ensure the audio file exists.

3. **WebSocket connection failed**:
   ```
   websocket is error
   ```
   Solution: Verify the WebSocket URL is correct and the service is running.

4. **Permission errors**:
   Solution: Ensure the script has write permissions to create the `temp/` directory.

5. **Output directory not found**:
   ```
   Error: ENOENT: no such file or directory, open '../temp/...'
   ```
   Solution: The script automatically creates the `temp/` directory. Ensure the parent directory exists and has write permissions.

## Development

### Building from Source
```bash
# Install dependencies
yarn install

# Build the project
yarn build

# Run the script
yarn start --wsUrl=<url> --audioPath=<path>

# For continuous development (watches for changes)
yarn watch
```

### Available Scripts
The project includes these npm scripts:
```json
{
  "scripts": {
    "build": "tsc",                    // Compile TypeScript to JavaScript
    "start": "node dist/src/index.js", // Run the compiled script
    "dev": "tsc && node dist/src/index.js", // Build and run in one command
    "watch": "tsc --watch",            // Watch for changes and recompile
    "clean": "rm -rf dist",            // Clean build output
  }
}
```

## License

**PROPRIETARY SOFTWARE - ALL RIGHTS RESERVED**

Copyright © 2024 Meeting.AI. All rights reserved.

This software and associated documentation files (the "Software") are the exclusive property of Meeting.AI. No part of this Software may be copied, modified, distributed, or used in any form or by any means without the express written permission of Meeting.AI.

**RESTRICTIONS:**
- ❌ **NO COPYING** - You may not copy, reproduce, or duplicate any part of this Software
- ❌ **NO MODIFICATION** - You may not modify, adapt, alter, translate, or create derivative works
- ❌ **NO DISTRIBUTION** - You may not distribute, sublicense, rent, lease, or transfer this Software
- ❌ **NO REVERSE ENGINEERING** - You may not reverse engineer, decompile, or disassemble this Software

**VIOLATION NOTICE:** Unauthorized use, copying, modification, or distribution of this Software is strictly prohibited and may result in legal action.

For licensing inquiries, contact: legal@meeting.ai

## Author

Meeting.AI