<p align="center">
  <img src="https://img.shields.io/badge/maintenance-actively--developed-brightgreen.svg" />
  <a href="https://goreportcard.com/report/github.com/knuspii/crunchycleaner"><img src="https://goreportcard.com/badge/github.com/knuspii/crunchycleaner" alt="Go Report Card" /></a>
  <a href="https://github.com/knuspii/crunchycleaner/actions/workflows/go.yml"><img src="https://github.com/knuspii/crunchycleaner/actions/workflows/go.yml/badge.svg" alt="Build" /></a>
  <a href="https://github.com/knuspii/crunchycleaner/stargazers"><img src="https://img.shields.io/github/stars/knuspii/crunchycleaner?style=social" alt="GitHub Stars" /></a>
  <br>
  <img src="https://img.shields.io/github/downloads/knuspii/crunchycleaner/total?color=green?" />
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue?logo=linux&logoColor=white" alt="Platform" />
</p>

<div align="center">
<h1>CrunchyCleaner <img src="assets/logo.png" width="64" height="64" alt="Logo"></h1>
<img src="assets/crunchycleaner-demo.gif" width="1000" height="1000" alt="Preview">
</div>

### 🧹 A lightweight, software [cache](https://wikipedia.org/wiki/Cache_(computing)) cleanup tool for Windows & Linux.
> Made by: Knuspii, (M)
- **Cross-Platform:** Works on both **Windows** and **Linux**
- **Lightweight:** Single binary, no dependencies (just download and run it)
- **TUI (Text-UI):** Simple, minimalist interface, no confusing menus

---

### Supported Software List:
| Category | Software / Path |
| :--- | :--- |
| **System** | System Logs, System Temp Folders, Thumbnail Cache, Shader Cache |
| **Browsers** | Firefox, Google Chrome, Microsoft Edge, Brave, Opera |
| **Development** | Visual Studio Code, Slack, JetBrains IDE, Go, Pip, NPM, Yarn, Cargo, NuGet, Gradle |
| **Server** | coming soon* |
| **Gaming** | Steam, Epic Games(Heroic) |
| **Apps** | Discord, Spotify, Thunderbird, Telegram |

**Flatpak is supported** \
**Snap is not supported**

## 📥 How to Install & Download ![Download](https://img.shields.io/github/downloads/knuspii/crunchycleaner/total?color=green)
Paste this into your terminal and restart it afterwards. \
For Linux (using sudo):
```
curl -L https://github.com/Knuspii/CrunchyCleaner/releases/latest/download/crunchycleaner -o cc && sudo install -m 755 cc /usr/local/bin/crunchycleaner && rm cc
```
For Windows (using Powershell as Admin):
```
iwr https://github.com/Knuspii/CrunchyCleaner/releases/latest/download/crunchycleaner.exe -OutFile cc.exe; mv -Force cc.exe $env:SystemRoot\System32\crunchycleaner.exe
```
For Go (using go install):
```
go install github.com/knuspii/crunchycleaner/v2@latest
```
Or download binary:
[[Download here]](https://github.com/knuspii/crunchycleaner/releases)

### 🎖️ Featured on
[![Awesome TUIs](https://img.shields.io/badge/Awesome-TUIs-orange?style=flat-square)](https://github.com/rothgar/awesome-tuis)
[![Awesome Go](https://img.shields.io/badge/Awesome-Go-blue?style=flat-square)](https://github.com/avelino/awesome-go)
[![Awesome Windows](https://img.shields.io/badge/Awesome-Windows-blueviiet?style=flat-square)](https://github.com/0PandaDEV/awesome-windows)

---

## Options:
```
  -a    Automate cleaning (select all and start immediately)
  -d    Simulation mode without deleting files (for testing)
  -t    Skip terminal resizing and environment initialization
  -v    Display version information
```

If you use a terminal emulator (kitty, konsole) you might wann use the `-t` flag.

---

> [!WARNING]
> You use this tool at your own risk!

> [!NOTE]
> AI was used for this project in some parts.

## Q&A
**Q:** Will this break my system? \
**A:** No.\
\
**Q:** What does it actually delete? \
**A:** Only cache files and temp files.\
\
**Q:** Why another cleaner? \
**A:** Because this one is easy, simple, small and lightweight.

## Other Downloads
CrunchyCleaner is also available on SourceForge \
[![Download CrunchyCleaner](https://a.fsdn.com/con/app/sf-download-button)](https://sourceforge.net/projects/crunchycleaner/files/latest/download)
[![Download CrunchyCleaner](https://img.shields.io/sourceforge/dt/crunchycleaner.svg)](https://sourceforge.net/projects/crunchycleaner/files/latest/download)

### External Dependencies
This project uses the following external dependencies:
- **[github.com/eiannone/keyboard](https://github.com/eiannone/keyboard)** – used for cross-platform keyboard input (MIT License)
- **[github.com/shirou/gopsutil](https://github.com/shirou/gopsutil)** – used for cross-platform system and hardware metrics (BSD 3-Clause License)
