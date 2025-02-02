# SnapResolve 🔍✨

**Instant AI-Powered Screenshot Analysis**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/JasonLovesDoggo/snapresolve/pulls)


## About  
SnapResolve lets you analyze anything on your screen with GPT-4 Vision using a single hotkey. Capture screenshots, get real-time AI insights in a non-intrusive overlay, and streamline your workflow across Windows, macOS, and Linux.

## Key Features  
🔑 **Custom Hotkeys** - Ctrl+Shift+S by default  
💡 **AI Insights** - GPT-4/Gemini Pro Vision integration with streaming responses  
🖥️ **System Tray** - Always available, never in your way  

## Installation (WIP) 🚧  
*Compiled binaries coming soon!*  

**Build from source**:  
```bash
# Requires Go 1.20+ and Node.js 18+
git clone https://github.com/JasonLovesDoggo/snapresolve
cd snapresolve
cd frontend && pnpm install
cd .. && wails build
```

## Development
```bash
# Frontend (Svelete)
cd frontend && pnpm dev

# Backend (Go)
wails dev
```

## Tech Stack
- **Backend**: Go (hotkey/screenshot logic)
- **Frontend**: Svelete + Vite (UI)
- **Packaging**: Wails v2

## License
MIT © [JasonLovesDoggo](https://github.com/JasonLovesDoggo)