<p align="center">
  <a href="https://github.com/knuspii/crunchycleaner/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/knuspii/crunchycleaner/go.yml?style=for-the-badge&logo=githubactions&logoColor=white&label=Build&labelColor=111111&color=00c853" alt="Build">
  </a>
  <a href="https://github.com/knuspii/crunchycleaner/releases">
    <img src="https://img.shields.io/github/v/release/knuspii/crunchycleaner?style=for-the-badge&logo=github&logoColor=white&label=Release&labelColor=111111&color=ff6d00" alt="Release">
  </a>
  <a href="https://github.com/knuspii/crunchycleaner/releases">
    <img src="https://img.shields.io/github/downloads/knuspii/crunchycleaner/total?style=for-the-badge&logo=github&logoColor=white&label=Downloads&labelColor=111111&color=8e24aa" alt="Downloads">
  </a>
  <a href="https://github.com/knuspii/crunchycleaner">
    <img src="https://img.shields.io/github/stars/knuspii/crunchycleaner?style=for-the-badge&logo=github&logoColor=white&label=Stars&labelColor=111111&color=fbc02d" alt="Stars">
  </a>
</p>

<div align="center">
<h1>CrunchyCleaner <img src="assets/logo.png" width="100" height="100" alt="Logo"></h1>
<img src="assets/crunchycleaner-demo.gif" width="1000" height="1000" alt="Preview">
</div>

### 🧹 A lightweight, software [cache](https://wikipedia.org/wiki/Cache_(computing)) cleanup tool for Windows & Linux.
> Made by: Knuspii, (M)
- **Cross-Platform:** Works on both **Windows** and **Linux**
- **Lightweight:** Single binary, no dependencies (just download and run it)
- **TUI (Text-UI):** Simple, minimalist interface, no confusing menus

---

## Supported Software List:
| Category | Software / Path |
| :--- | :--- |
| **System:** | System Logs, System Temp, Thumbnail, Shell History, Package Managers, Shader Cache |
| **Browsers:** | Firefox, Google Chrome, Microsoft Edge, Brave, Opera, Zen |
| **Apps:** | Steam, Epic Games(Heroic), Discord, Telegram, Spotify |
| **Development:** | Visual Studio Code, JetBrains IDE, Slack, Go, Pip, NPM, Yarn, Cargo, Electron, NuGet, Gradle, Playwright, Python uv, Node.js |

**Flatpak is supported** \
**Snap is NOT supported**

---

## 📥 How to Install & Download
### Easiest way:
[Simply download the latest release from here ⬇ and run the binary.](https://github.com/Knuspii/CrunchyCleaner/releases/latest) \
[![DownloadCrunchyCleaner](https://img.shields.io/badge/⬇%20Download%20CrunchyCleaner-Latest%20Release-brightgreen?style=for-the-badge)](https://github.com/Knuspii/CrunchyCleaner/releases/latest)


### Alternative: Install via Terminal
If you prefer installing CrunchyCleaner through your terminal, use one
of the following commands.

**Linux (using sudo):**
``` bash
curl -L https://github.com/Knuspii/CrunchyCleaner/releases/latest/download/crunchycleaner -o cc && sudo install -m 755 cc /usr/local/bin/crunchycleaner && rm cc
```

**Windows (using PowerShell as Administrator):**
``` powershell
$d="$env:ProgramFiles\CrunchyCleaner"; md $d -F; iwr "https://github.com/Knuspii/CrunchyCleaner/releases/latest/download/crunchycleaner.exe" -OutFile "$d\crunchycleaner.exe"; $p=[Environment]::GetEnvironmentVariable('Path',2); if($p -notlike "$d"){[Environment]::SetEnvironmentVariable('Path',"$p;$d",2)}
```

**Go (using `go install`):**
``` bash
go install github.com/knuspii/crunchycleaner/v2@latest
```

After installation, restart your terminal and verify it by running:
``` bash
crunchycleaner -v
```

---

## Options:
```
  -a	Automate cleaning (select all and start immediately)
  -d	Dry-run mode without deleting files (for testing)
  -s	Show all hotkeys
  -v	Show version information
```

---

> [!WARNING]
> You use this tool at your own risk!

> [!NOTE]
> AI was used for this project in some parts. BUT THIS IS NOT VIBECODED!

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
\
[![Download CrunchyCleaner](https://a.fsdn.com/con/app/sf-download-button)](https://sourceforge.net/projects/crunchycleaner/files/latest/download)
[![Total Downloads](https://img.shields.io/sourceforge/dt/crunchycleaner.svg)](https://sourceforge.net/projects/crunchycleaner/files/latest/download)

---

### 🎖️ Featured on
 - [Awesome TUIs](https://github.com/rothgar/awesome-tuis)
 - [Awesome Go](https://github.com/avelino/awesome-go)
 - [Awesome Windows](https://github.com/0PandaDEV/awesome-windows)
 - [Linux Links](https://www.linuxlinks.com/crunchycleaner-lightweight-cache-cleanup-tool)

### External Dependencies
This project uses the following external dependencies:
- **[github.com/eiannone/keyboard](https://github.com/eiannone/keyboard)** – used for cross-platform keyboard input (MIT License)
- **[golang.org/x/term](https://pkg.go.dev/golang.org/x/term)** – used for the menu list (BSD 3-Clause License)
- **[github.com/shirou/gopsuti](https://github.com/shirou/gopsutil)** – used for cross-platform system and hardware metrics (BSD 3-Clause License)
