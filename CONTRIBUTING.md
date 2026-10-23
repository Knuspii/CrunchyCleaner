# 🚀 How Can you Contribute?
## 1. Adding New Programs/Caches
The most common way to contribute is by adding support for more applications.
* Open the cc_caches.go file.
* Locate the getPrograms() function
* Add a new Program struct to the slice.
* **Important:** Ensure that the application is kinda popular, I won't add SuperEditorX5000 lol.
* **Important:** Ensure you provide paths for both Windows and Linux where possible.
* Use filepath.Join and environment variables (like localAppData) and maintain cross-platform compatibility.

## 2. Improving the TUI
I aim for a "Crunchy" retro terminal feel. Improvements to the menu navigation, spinner, or banner are welcome, provided they don't add external dependencies.

## 3. Bug Reports & Feature Requests
If you find a bug or have an idea:
* Check the Issues tab to see if it has already been reported.
* Open a new issue with a clear title and description of the environment (OS, Terminal).

## 📜 Code Style Guidelines
* Go Idioms: Follow standard Go formatting (go fmt).
* Cross-Platform First: Always check if your code breaks on the "other" OS.
* No Heavy Dependencies: I prefer standard library solutions or very lightweight packages (like keyboard).

## 📥 Pull Request Process
* Fork the repository and create your branch from main.
* Test your changes use ```go test .```
* Run a dry-run ```go run . -d``` to ensure paths are detected correctly.
* Commit with descriptive messages (e.g., "Add: Spotify cache support").
* Open a PR and describe what your changes do and why they are necessary.

## ⚖️ License
By contributing, you agree that your contributions will be licensed under the GNU General Public License v3.0.
