// ##################################################################
// CrunchyCleaner
// Here are all paths to the specific caches.
// ##################################################################

package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func getPrograms() []Program {
	if runtime.GOOS == "windows" {
		// Windows
		home, _ := os.UserHomeDir()
		appData := os.Getenv("APPDATA")
		localAppData := os.Getenv("LOCALAPPDATA")
		programFilesX86 := os.Getenv("ProgramFiles(x86)")
		programFiles := os.Getenv("ProgramFiles")
		programData := os.Getenv("ProgramData")
		winDir := os.Getenv("WINDIR")
		return []Program{
			{"System Logs (Admin)", []string{
				filepath.Join(winDir, "Panther"),
				filepath.Join(winDir, "Logs"),
			}, false},
			{"Font Cache (Admin)", []string{filepath.Join(winDir, "ServiceProfiles/LocalService/AppData/Local/FontCache")}, false},
			{"System Temp Folders (Admin)", []string{filepath.Join(winDir, "Temp")}, false},
			{"Update Logs (Admin)", []string{filepath.Join(winDir, "SoftwareDistribution/Download")}, false},
			{"User Temp Folder", []string{filepath.Join(localAppData, "Temp")}, false},
			{"Thumbnail Cache", []string{filepath.Join(localAppData, "Microsoft/Windows/Explorer")}, false},
			{"Powershell History", []string{filepath.Join(appData, "Microsoft/Windows/PowerShell/PSReadLine")}, false},
			{"Package Manager Caches", []string{
				filepath.Join(home, "scoop/apps/scoop/current/cache"),
				filepath.Join(programData, "chocolatey/cache"),
			}, false},
			{"Firefox Cache", []string{
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/cache2"),
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/jumpListCache"),
				filepath.Join(appData, "Mozilla/Firefox/Profiles/*/shader-cache"),
			}, false},
			{"Chrome Cache", []string{
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Code Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/*/Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Media Cache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(localAppData, "Microsoft/Edge/User Data/Default/Cache"),
				filepath.Join(localAppData, "Microsoft/Edge/User Data/*/Cache"),
				filepath.Join(localAppData, "Microsoft/Edge/User Data/Default/Media Cache"),
			}, false},
			{"Brave Cache", []string{
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/Default/Cache"),
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/*/Cache"),
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/Default/Media Cache"),
			}, false},
			{"Opera Cache", []string{
				filepath.Join(localAppData, "Opera Software/Opera Stable/Cache"),
				filepath.Join(localAppData, "Opera Software/Opera Stable/Code Cache"),
			}, false},
			{"Zen Browser Cache", []string{
				filepath.Join(localAppData, "zen/Profiles/*/cache2"),
				filepath.Join(localAppData, "zen/Profiles/*/jumpListCache"),
				filepath.Join(appData, "zen/Profiles/*/shader-cache"),
			}, false},
			{"Vivaldi Cache", []string{
				filepath.Join(localAppData, "Vivaldi/User Data/Default/Cache"),
				filepath.Join(localAppData, "Vivaldi/User Data/*/Cache"),
				filepath.Join(localAppData, "Vivaldi/User Data/Default/Media Cache"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(localAppData, "Thunderbird/Profiles/*/cache2"),
			}, false},
			{"Steam Cache", []string{
				filepath.Join(programFilesX86, "Steam/appcache"),
				filepath.Join(programFiles, "Steam/appcache"),
				filepath.Join(localAppData, "Steam/htmlcache"),
			}, false},
			{"Epic Games Cache", []string{filepath.Join(localAppData, "EpicGamesLauncher/Saved/webcache")}, false},
			{"Discord Cache", []string{
				filepath.Join(appData, "discord/Cache"),
				filepath.Join(appData, "discord/Code Cache"),
				filepath.Join(appData, "discord/GPUCache"),
			}, false},
			{"Telegram Cache", []string{filepath.Join(appData, "Telegram Desktop/tdata/user_data/cache")}, false},
			{"Spotify Cache", []string{filepath.Join(localAppData, "Spotify/Storage")}, false},
			{"VS Code Cache", []string{
				filepath.Join(appData, "Code/Cache"),
				filepath.Join(appData, "Code/CachedData"),
				filepath.Join(appData, "Code/CachedExtensionVSIXs"),
				filepath.Join(appData, "Code/User/workspaceStorage"),
				filepath.Join(appData, "Code/GPUCache"),
			}, false},
			{"JetBrains IDE Cache", []string{filepath.Join(localAppData, "JetBrains/*/system/caches")}, false},
			{"Slack Cache", []string{
				filepath.Join(appData, "Slack/Cache"),
				filepath.Join(appData, "Slack/Code Cache"),
				filepath.Join(appData, "Slack/GPUCache"),
			}, false},
			{"Shader Cache", []string{
				filepath.Join(localAppData, "D3DSCache"),
				filepath.Join(localAppData, "NVIDIA/GLCache"),
				filepath.Join(localAppData, "NVIDIA/PerDriverVersion/DXCache"),
			}, false},
			{"Golang Cache", []string{filepath.Join(localAppData, "go-build")}, false},
			{"Pip Cache", []string{filepath.Join(localAppData, "pip/cache")}, false},
			{"NPM Cache", []string{filepath.Join(appData, "npm-cache/_cacache")}, false},
			{"Yarn Cache", []string{
				filepath.Join(localAppData, "Yarn/Cache"),
				filepath.Join(appData, "Yarn/Cache"),
			}, false},
			{"Cargo Cache", []string{
				filepath.Join(home, ".cargo/registry/cache"),
				filepath.Join(home, ".cargo/git/db"),
			}, false},
			{"Electron Cache", []string{
				filepath.Join(localAppData, "electron/Cache"),
				filepath.Join(localAppData, "electron/Code Cache"),
				filepath.Join(localAppData, "electron/GPUCache"),
			}, false},
			{"NuGet Cache", []string{filepath.Join(home, ".nuget/packages")}, false},
			{"Gradle Cache", []string{filepath.Join(home, ".gradle/caches")}, false},
			{"Playwright Cache", []string{filepath.Join(localAppData, "ms-playwright")}, false},
			{"Python uv Cache", []string{filepath.Join(localAppData, "uv/cache")}, false},
			{"Node.js Cache", []string{filepath.Join(localAppData, "node-gyp")}, false},
			{"Bun Cache", []string{filepath.Join(home, ".bun/install/cache")}, false},
			{"Windows WebCache", []string{filepath.Join(localAppData, "Microsoft/Windows/WebCache")}, false},
			{"Android Studio", []string{
				filepath.Join(localAppData, "Google", "AndroidStudio*", "caches"),
			}, false},
			{"Claude Cache", []string{
				filepath.Join(appData, "Claude/Cache"),
				filepath.Join(appData, "Claude/Code Cache"),
				filepath.Join(appData, "Claude/GPUCache"),
			}, false},
			{"User Cache", []string{filepath.Join(home, ".cache")}, false},
		}
	} else {
		// Linux
		home, _ := os.UserHomeDir()
		cache := ".cache/"
		flatpak := ".var/app/"
		return []Program{
			{"System Logs (Root)", []string{"/var/log/"}, false},
			{"System Temp Folders (Root)", []string{"/tmp"}, false},
			{"Package Manager Caches (Root)", []string{
				"/var/cache/apt/archives/*",
				"/var/cache/apt/pkgcache.bin",
				"/var/cache/apt/srcpkgcache.bin",
				"/var/cache/pacman/pkg/*",
				"/var/cache/dnf/*",
				"/var/cache/snapd/*",
				"/var/cache/apk/*",
			}, false},
			{"Thumbnail Cache", []string{filepath.Join(home, cache, "thumbnails")}, false},
			{"Shell History", []string{filepath.Join(home, "*_history")}, false},
			{"Firefox Cache", []string{
				filepath.Join(home, cache, "mozilla/firefox/*/cache2"),
				filepath.Join(home, flatpak, "org.mozilla.firefox/cache/mozilla/firefox/*/cache2"),
			}, false},
			{"Chromium Cache", []string{
				filepath.Join(home, cache, "chromium/*/Cache"),
				filepath.Join(home, cache, "chromium/*/Code Cache"),
				filepath.Join(home, flatpak, "com.google.Chrome/cache/chromium/*/Cache"),
				filepath.Join(home, flatpak, "com.google.Chrome/cache/chromium/*/Code Cache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(home, cache, "microsoft-edge/*/Cache"),
				filepath.Join(home, cache, "microsoft-edge/*/Code Cache"),
				filepath.Join(home, flatpak, "com.microsoft.Edge/cache/microsoft-edge/*/Cache"),
				filepath.Join(home, flatpak, "com.microsoft.Edge/cache/microsoft-edge/*/Code Cache"),
			}, false},
			{"Brave Cache", []string{
				filepath.Join(home, cache, "BraveSoftware/Brave-Browser/*/Cache"),
				filepath.Join(home, cache, "BraveSoftware/Brave-Browser/*/Code Cache"),
				filepath.Join(home, flatpak, "com.brave.Browser/cache/Brave-Browser/*/Cache"),
				filepath.Join(home, flatpak, "com.brave.Browser/cache/Brave-Browser/*/Code Cache"),
			}, false},
			{"Opera Cache", []string{
				filepath.Join(home, cache, "opera/Cache"),
				filepath.Join(home, ".config/opera/Cache"),
				filepath.Join(home, flatpak, "com.opera.Opera/cache/opera/Cache"),
				filepath.Join(home, flatpak, "com.opera.Opera/config/opera/Cache"),
			}, false},
			{"Zen Browser Cache", []string{
				filepath.Join(home, cache, "zen/*/cache2"),
				filepath.Join(home, flatpak, "app.zen_browser.zen/cache/zen/*/cache2"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(home, cache, "thunderbird/*/cache2"),
				filepath.Join(home, flatpak, "org.mozilla.Thunderbird/cache/mozilla/Thunderbird/*/cache2"),
			}, false},
			{"Steam Cache", []string{
				filepath.Join(home, ".steam/steam/appcache"),
				filepath.Join(home, ".local/share/Steam/appcache"),
				filepath.Join(home, ".local/share/Steam/config/htmlcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/steam/steam/appcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/.local/share/Steam/appcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/.local/share/Steam/config/htmlcache"),
			}, false},
			{"Epic Games (Heroic/Lutris) Cache", []string{
				filepath.Join(home, ".config/heroic/WebCache"),
				filepath.Join(home, ".local/share/lutris/runtime"),
				filepath.Join(home, flatpak, "com.heroicgameslauncher.hgl/config/heroic/WebCache"),
				filepath.Join(home, flatpak, "com.heroicgameslauncher.hgl/.local/share/lutris/runtime"),
			}, false},
			{"Discord Cache", []string{
				filepath.Join(home, ".config/discord/Cache"),
				filepath.Join(home, ".config/discord/Code Cache"),
				filepath.Join(home, ".config/discord/GPUCache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/Cache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/Code Cache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/GPUCache"),
			}, false},
			{"Telegram Cache", []string{filepath.Join(
				home, ".local/share/TelegramDesktop/tdata/user_data/cache"),
				filepath.Join(home, flatpak, "org.telegram.desktop/data/TelegramDesktop/tdata/user_data/cache"),
			}, false},
			{"Spotify Cache", []string{
				filepath.Join(home, cache, "spotify"),
				filepath.Join(home, flatpak, "com.spotify.Client/cache/spotify"),
			}, false},
			{"VS Code Cache", []string{
				filepath.Join(home, ".config/Code/Cache"),
				filepath.Join(home, ".config/Code/Code Cache"),
				filepath.Join(home, ".config/Code/CachedData"),
				filepath.Join(home, ".config/Code/GPUCache"),
				filepath.Join(home, ".config/Code/User/workspaceStorage"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/Cache"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/Code Cache"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/CachedData"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/GPUCache"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/User/workspaceStorage"),
			}, false},
			{"JetBrains IDE Cache", []string{filepath.Join(home, cache, "JetBrains/*/caches")}, false},
			{"Slack Cache", []string{
				filepath.Join(home, ".config/Slack/Cache"),
				filepath.Join(home, ".config/Slack/Code Cache"),
				filepath.Join(home, ".config/Slack/GPUCache"),
				filepath.Join(home, flatpak, "com.slack.Slack/config/Slack/Cache"),
				filepath.Join(home, flatpak, "com.slack.Slack/config/Slack/Code Cache"),
				filepath.Join(home, flatpak, "com.slack.Slack/config/Slack/GPUCache"),
			}, false},
			{"Nix Cache", []string{
				filepath.Join(home, ".cache/nix")}, false},
			{"Shader Cache", []string{
				filepath.Join(home, cache, "mesa_shader_cache"),
				filepath.Join(home, cache, "nvidia/GLCache"),
			}, false},
			{"Golang Cache", []string{filepath.Join(home, cache, "go-build")}, false},
			{"Pip Cache", []string{filepath.Join(home, cache, "pip")}, false},
			{"NPM Cache", []string{filepath.Join(home, ".npm/_cacache")}, false},
			{"Yarn Cache", []string{filepath.Join(home, cache, "yarn")}, false},
			{"Cargo Cache", []string{filepath.Join(home, ".cargo/registry/cache")}, false},
			{"Electron Cache", []string{filepath.Join(home, cache, "electron")}, false},
			{"NuGet Cache", []string{filepath.Join(home, ".nuget/packages")}, false},
			{"Gradle Cache", []string{filepath.Join(home, ".gradle/caches")}, false},
			{"Playwright Cache", []string{filepath.Join(home, cache, "ms-playwright")}, false},
			{"Python uv Cache", []string{filepath.Join(home, cache, "uv")}, false},
			{"Node.js Cache", []string{filepath.Join(home, cache, "node-gyp")}, false},
		}
	}
}
