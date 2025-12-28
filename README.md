<h1>THIS PROGRAM IS CURRENTLY BEING REWRITTEN COMPLETLY</h1>

[![Golang](https://img.shields.io/static/v1?label=Made%20with&message=Go&logo=go&color=007ACC)](https://go.dev/)
[![Go version](https://img.shields.io/github/go-mod/go-version/knuspii/crunchycleaner)](https://github.com/knuspii/crunchycleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/Knuspii/crunchycleaner)](https://goreportcard.com/report/github.com/Knuspii/crunchycleaner)
[![Build](https://github.com/Knuspii/crunchycleaner/actions/workflows/go.yml/badge.svg)](https://github.com/Knuspii/crunchycleaner/actions/workflows/go.yml)
[![GitHub Issues](https://img.shields.io/github/issues/knuspii/crunchycleaner)](https://github.com/knuspii/crunchycleaner/issues)
[![GitHub Stars](https://img.shields.io/github/stars/knuspii/crunchycleaner?style=social)](https://github.com/knuspii/crunchycleaner/stargazers)

<h1>CrunchyCleaner 🧹</h1>

![Preview](preview.png)

### ✨ A lightweight, cross-platform system cleanup tool
CrunchyCleaner is made to be simple, easy and *very crunchy indeed!*\
It helps you clear out caches, temp files, logs, and more — without confusing menus and 100+ options.

## 📥 [[Download here]](https://github.com/Knuspii/crunchycleaner/releases) <- Click here to download CrunchyCleaner!

## 🔑 Key features:

- 💻 **Cross-Platform**: Works on both **Windows** and **Linux**
- ⚡ **Lightweight**: Single binary, no dependencies (just download and run it)
- 🎨 **TUI (Text-UI)**: Simple, minimalist interface, no confusing menus
- 🧹 **Multiple Modes**:
  - Safe Clean (harmless cache cleanup)
  - Full Clean (deep cleanup of system junk)
  - User Clean (profile-specific cleanup)

## ⚙️ Start options:
```
Usage:
  crunchycleaner [option]

Options:
  No option   Run with TUI (Text-UI)
  -t          Run with TUI (Text-UI)
  -s          Run Safe-Cleanup
  -sy         Run Safe-Cleanup (non-interactive for scripts)
  -f          Run Full-Cleanup
  -fy         Run Full-Cleanup (non-interactive for scripts)
  -u <user>   Run User-Cleanup
  -uy <user>  Run User-Cleanup (non-interactive for scripts)
  -v          Show version
  -h          Show this help page
```

> [!WARNING]
> You use this tool at your own risk!

## External Dependencies
This project uses the following external dependencies:
- **[github.com/eiannone/keyboard](https://github.com/eiannone/keyboard)** – used for cross-platform keyboard input (MIT License)
