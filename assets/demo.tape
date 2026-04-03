# Output settings
Output crunchycleaner-demo.gif
Set FontSize 17
Set Width 1200
Set Height 600
Set TypingSpeed 100ms
Set Padding 10
Set WindowBar Colorful

Hide
  Type "export PATH=$PATH:$(pwd)/../bin && clear"
  Enter
Show

# --- Start the TUI ---
Type "crunchycleaner -d # This is a test run!"
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
Enter
Down @ 500ms
Enter
Down @ 500ms
Sleep 1s
Type "a"
Sleep 2s
Type "c"
Sleep 6s
Enter

