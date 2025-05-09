# Spotify Auto Genius Opener

A Windows application that automatically opens Genius lyrics pages for currently playing Spotify tracks.

## Features

- 🔍 Monitors currently playing track in Spotify
- 🎵 Automatically detects track changes
- 📖 Opens the corresponding Genius lyrics page
- 🔄 Falls back to Genius search if exact match not found
- ⚡ Real-time tracking with minimal resource usage

## How It Works

1. The application continuously monitors Spotify's window title
2. When a new track starts playing, it:
   - Extracts the track information from the window title
   - Searches for the track on Genius
   - Opens the exact lyrics page if found
   - Falls back to Genius search results if no exact match is found

## Requirements

- Windows operating system
- Spotify desktop application
- Go 1.16 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/spotify-auto-genius-opener.git
cd spotify-auto-genius-opener
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build
```

## Usage

1. Start Spotify and play some music
2. Run the application:
```bash
./spotify-auto-genius-opener
```

The application will automatically:
- Monitor your Spotify playback
- Open Genius pages for new tracks
- Show search results if exact match isn't found

## License

MIT License - feel free to use and modify as you wish!
