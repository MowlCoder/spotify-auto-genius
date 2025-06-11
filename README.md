# Spotify Auto Genius Opener

An application that automatically opens Genius lyrics pages for currently playing Spotify tracks.

## Features

- 🔍 Monitors currently playing track in Spotify
- 🎵 Automatically detects track changes
- 📖 Opens the corresponding Genius lyrics page
- 🔄 Falls back to Genius search if exact match not found
- ⚡ Real-time tracking with minimal resource usage
- 💻 Cross-platform support (Windows & Linux)

## How It Works

**Windows:**
1. The application continuously monitors Spotify's window title
2. When a new track starts playing, it:
   - Extracts the track information from the window title
   - Searches for the track on Genius
   - Opens the exact lyrics page if found
   - Falls back to Genius search results if no exact match is found

**Linux:**
1. The application uses D-Bus to get track's metadata from Spotify
2. When a new track starts playing, it:
   - Gets track's artist and title
   - Searches for the track on Genius
   - Opens the exact lyrics page if found
   - Falls back to Genius search results if no exact match is found

## Requirements

- Windows or Linux operating system
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
go build -o spotify-auto-genius
```

## Usage

1. Start Spotify and play some music
2. Run the application:
```bash
./spotify-auto-genius
```

The application will automatically:
- Monitor your Spotify playback
- Open Genius pages for new tracks
- Show search results if exact match isn't found

## License

MIT License - feel free to use and modify as you wish!
