# DevBoot Demo — VHS Tape
# Prerequisites: vhs, ffmpeg, ttyd (all installed)
# Run: vhs demo/demo.tape
# Output: demo/demo.gif

Output demo/demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1000
Set Height 600
Set Padding 20
Set Theme { "name": "Catppuccin Mocha", "black": "#45475A", "red": "#F38BA8", "green": "#A6E3A1", "yellow": "#F9E2AF", "blue": "#89B4FA", "magenta": "#F5C2E7", "cyan": "#94E2D5", "white": "#BAC2DE", "brightBlack": "#585B70", "brightRed": "#F38BA8", "brightGreen": "#A6E3A1", "brightYellow": "#F9E2AF", "brightBlue": "#89B4FA", "brightMagenta": "#F5C2E7", "brightCyan": "#94E2D5", "brightWhite": "#A6ADC8", "background": "#1E1E2E", "foreground": "#CDD6F4", "selection": "#585B70", "cursor": "#F5E0DC" }
Set TypingSpeed 40ms
Set WindowBar Colorful

Sleep 500ms

# ── version ──
Type "devboot version"
Enter
Sleep 2s

# ── profiles ──
Type "devboot profile list"
Enter
Sleep 3s

# ── profile detail ──
Type "devboot profile show terminal"
Enter
Sleep 3s

# ── status dashboard ──
Type "devboot status demo/devboot.yaml"
Enter
Sleep 4s

# ── diff preview ──
Type "devboot diff demo/devboot.yaml"
Enter
Sleep 4s

# ── apply git ──
Type "devboot apply demo/devboot.yaml --only git --no-tui"
Enter
Sleep 3s

# ── export ──
Type "devboot profile export rust"
Enter
Sleep 3s

Sleep 1s
