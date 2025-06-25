# Backend API Documentation

The Terminal DevTool backend is a Go service that processes media files using FFmpeg. This document describes the available API endpoints and their usage.

## Endpoints

### Health Check
```
GET /health
```

Returns information about the server status and FFmpeg availability.

Example response:
```json
{
  "status": "OK",
  "version": "0.1.0",
  "ffmpeg_available": true,
  "ffmpeg_version": "ffmpeg version 4.4.1",
  "components": {
    "libx264": "available",
    "libvpx": "available",
    "libopus": "available"
  }
}
```

### Process Media
```
POST /api/process
```

Process a media file with specified options.

Request body:
```json
{
  "input": "input.mp4",
  "output": "output.webm",
  "resolution": "1280x720",
  "bitrate": "800k",
  "format": "webm",
  "codec": "libvpx-vp9",
  "frame_rate": "30",
  "crf": "23",
  "preset": "medium"
}
```

Parameters:
- `input` (required): Path to the input file
- `output` (optional): Path to the output file. If not provided, a default name will be generated
- `resolution` (optional): Output resolution (e.g., "1280x720")
- `bitrate` (optional): Video bitrate (e.g., "800k")
- `format` (optional): Output format (e.g., "mp4", "webm", "gif")
- `codec` (optional): Video codec (e.g., "libx264", "libvpx-vp9")
- `frame_rate` (optional): Frame rate (e.g., "30")
- `crf` (optional): Constant Rate Factor for quality-based compression (lower is better, e.g., "23")
- `preset` (optional): Encoding preset (e.g., "ultrafast", "fast", "medium", "slow")

Response:
```json
{
  "output": "output.webm"
}
```

### Compare Media
```
POST /api/compare
```

Compare two media files to analyze differences.

Request body:
```json
{
  "original": "original.mp4",
  "processed": "processed.mp4"
}
```

Parameters:
- `original` (required): Path to the original file
- `processed` (required): Path to the processed file

Response:
```json
{
  "original": {
    "filename": "original.mp4",
    "format": "mp4",
    "duration": "60.0s",
    "resolution": "1920x1080",
    "bitrate": "5000000",
    "size": 37500000,
    "codec": "h264",
    "frame_rate": "30 fps"
  },
  "processed": {
    "filename": "processed.mp4",
    "format": "mp4",
    "duration": "60.0s",
    "resolution": "1280x720",
    "bitrate": "2000000",
    "size": 15000000,
    "codec": "h264",
    "frame_rate": "30 fps"
  },
  "size_diff_percent": 60.0,
  "resolution_changed": true,
  "bitrate_reduction_percent": 60.0,
  "format_changed": false,
  "codec_changed": false
}
```

### Compress Media
```
POST /api/compress
```

Compress a video file using specified bitrate.

Request body:
```json
{
  "input": "input.mp4",
  "output": "compressed.mp4",
  "bitrate": "800k"
}
```

Parameters:
- `input` (required): Path to the input file
- `output` (optional): Path to the output file. If not provided, a default name will be generated (original_filename_compressed.ext)
- `bitrate` (required): Target bitrate (must end with 'k' or 'M', e.g., "800k", "2M")

Response:
```json
{
  "output": "compressed.mp4",
  "status": "success",
  "message": "Video compressed successfully"
}
```

### Get Media Info
```
GET /api/info?path=file.mp4
```

Get detailed information about a media file.

Parameters:
- `path` (required): Path to the media file

Response:
```json
{
  "filename": "file.mp4",
  "format": "mp4",
  "duration": "120.5s",
  "resolution": "1920x1080",
  "bitrate": "5000000",
  "size": 75000000,
  "codec": "h264",
  "frame_rate": "30 fps"
}
```

## Error Handling

API errors are returned with appropriate HTTP status codes and error messages:

```json
{
  "error": "Error message"
}
```

Common status codes:
- `400 Bad Request`: Invalid input parameters
- `404 Not Found`: The requested resource was not found
- `500 Internal Server Error`: Server-side error processing the request
