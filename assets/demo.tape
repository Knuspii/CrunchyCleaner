# Output settings
Output crunchycleaner-demo.gif
Set FontSize 15
Set Width 1200
Set Height 600
Set TypingSpeed 100ms
Set WindowBar Colorful

Hide
  Type "export PATH=$PATH:$(pwd)/../bin && clear"
  Enter
Show

# --- Start the TUI ---
Type "crunchycleaner -d"
Sleep 500ms
Enter
Sleep 2s # Wait for the TUI to render

# --- Navigate the Menu ---
# Simulate the user browsing through the options
Down @ 500ms
Sleep 500ms
Down @ 500ms
Sleep 500ms
Up @ 800ms
Sleep 1s
Type " "
Down @ 500ms
Type " "
Sleep 2s
Type "a"
Sleep 2s
Type "c"
Sleep 5s

# --- Exit ---
Enter @ 500ms
Sleep 3s
