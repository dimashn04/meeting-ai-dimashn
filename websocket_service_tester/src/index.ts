import { Readable } from 'stream'
import WebSocket from 'ws'
import * as fs from 'fs'

import { spawn } from 'child_process'
import path from 'path'
import { ulid } from 'ulid'

class AudioConverter {
  constructor() {
    // Check if ffmpeg is installed
    this.checkFFmpeg()
  }

  checkFFmpeg() {
    try {
      spawn('ffmpeg', ['-version'])
    } catch (error) {
      throw new Error('FFmpeg is not installed. Please install it first.')
    }
  }

  convertToPCMStream(inputFile: string): Promise<Readable> {
    return new Promise((resolve, reject) => {
      // Check if input file exists
      if (!fs.existsSync(inputFile)) {
        reject(new Error('Input file does not exist'))
        return
      }

      // FFmpeg command to convert to PCM S16LE
      const ffmpeg = spawn('ffmpeg', [
        '-i', inputFile,        // Input file
        '-f', 's16le',         // Force format to signed 16-bit little-endian
        '-acodec', 'pcm_s16le', // Audio codec
        '-ac', '1',            // Mono channel
        '-ar', '16000',        // Sample rate 16kHz
        'pipe:1'               // Output to stdout
      ])

      // Handle errors
      ffmpeg.stderr.on('data', (data) => {
        console.error(`FFmpeg stderr: ${data}`)
      })

      ffmpeg.on('error', (error) => {
        reject(new Error(`FFmpeg error: ${error.message}`))
      })

      ffmpeg.on('close', (code) => {
        if (code !== 0) {
          reject(new Error(`FFmpeg process exited with code ${code}`))
        }
      })

      // Return the readable stream
      resolve(ffmpeg.stdout)
    })
  }
}

// Example usage
async function main() {
  const fileId = ulid()

  // find url from args
  const strWsUrl = process.argv.find((arg) => arg.startsWith('--wsUrl='))?.split('=')[1]
  if (!strWsUrl) {
    throw new Error('wsUrl is required')
  }

  const audioPath = process.argv.find((arg) => arg.startsWith('--audioPath='))?.split('=')[1]
  if (!audioPath) {
    throw new Error('audioPath is required')
  }

  let connectionId = process.argv.find((arg) => arg.startsWith('--connectionId='))?.split('=')[1]
  if (!connectionId) {
    connectionId = fileId
  }

  const wsUrl = new URL(strWsUrl)
  wsUrl.searchParams.set('connection_id', connectionId)

  const outputPath = path.join(__dirname, '/../temp', `${fileId}.txt`)

  // Ensure temp directory exists
  const tempDir = path.dirname(outputPath)
  if (!fs.existsSync(tempDir)) {
    fs.mkdirSync(tempDir, { recursive: true })
  }

  // tracker
  let receivedBytes = 0
  let totalDurationMs = 0
  let totalSentBytes = 0

  const ws = new WebSocket(wsUrl.toString())

  await new Promise<void>(res => {
    ws.onopen = () => {
      console.log('websocket is open')

      res()
    }

    ws.onerror = (error) => {
      console.error('websocket is error', error)
    }
  })

  ws.onclose = (_e) => {
    console.log('websocket is closed')

    process.exit(0)
  }

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data as string)
    console.log(data)


    if (!data.text) {
      console.log('Unknown message', data)
      return
    }

    fs.appendFileSync(outputPath, `${data.start.toFixed(2)} - ${(data.end).toFixed(2)}: ${data.text}\n`)
  }

  // init convert
  const converter = new AudioConverter()
  try {
    const pcmStream = await converter.convertToPCMStream(audioPath)

    // Example: Process the stream directly
    pcmStream.on('data', (chunk) => {
      receivedBytes += chunk.length

      // Process PCM data chunks here
      setTimeout((totalDurationMs: number) => {
        if (ws.readyState !== WebSocket.OPEN) {
          console.log('WebSocket is not open. Cannot send data.')
          return
        }

        ws.send(chunk)

        totalSentBytes += chunk.length
        console.log('Sent chunk of size:', chunk.length, totalDurationMs / 1000, totalSentBytes)
      }, Math.floor(totalDurationMs), totalDurationMs)
      
      totalDurationMs += chunk.length / (16000 * 2) * 1000
    })

    pcmStream.on('end', () => {
      console.log('Conversion completed')
    })
  } catch (error) {
    console.error('Error:', error instanceof Error ? error.message : String(error))
  }
}

(async () => {
  await main()
})()