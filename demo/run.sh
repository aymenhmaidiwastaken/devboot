#!/bin/bash
# ============================================================================
#  DevBoot Demo — Run this script to see all features in action
#  Usage: bash demo/run.sh
# ============================================================================

set -e

BINARY="go run ."
CYAN='\033[36m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

pause() {
    echo ""
    echo -e "${DIM}────────────────────────────────────────────────────────${RESET}"
    echo ""
    sleep 1
}

header() {
    echo ""
    echo -e "${BOLD}${CYAN}  ▸ $1${RESET}"
    echo ""
    sleep 0.5
}

# ── Intro ──────────────────────────────────────────────────────────────────

clear
echo ""
echo -e "${BOLD}${CYAN}"
echo '    ____             ____              __  '
echo '   / __ \___  _   __/ __ )____  ____  / /_ '
echo '  / / / / _ \| | / / __  / __ \/ __ \/ __/ '
echo ' / /_/ /  __/| |/ / /_/ / /_/ / /_/ / /_   '
echo '/_____/\___/ |___/_____/\____/\____/\__/   '
echo ""
echo -e "  Fresh machine to productive in one command.${RESET}"
echo ""
sleep 2

# ── 1. Version ─────────────────────────────────────────────────────────────

header "1. devboot version"
$BINARY version
pause

# ── 2. Doctor — diagnose your environment ──────────────────────────────────

header "2. devboot doctor — diagnose your environment"
$BINARY doctor
pause

# ── 3. Profile system — curated tool sets ──────────────────────────────────

header "3. devboot profile list — curated profiles"
$BINARY profile list
pause

header "3b. devboot profile show terminal — inspect a profile"
$BINARY profile show terminal
pause

header "3c. devboot profile search docker — search profiles"
$BINARY profile search docker
pause

# ── 4. Export — reverse-engineer current setup ─────────────────────────────

header "4. devboot export — scan and export your current environment"
$BINARY export
pause

# ── 5. Init — generate config (plain mode for demo) ───────────────────────

header "5. devboot init --plain — generate starter config"
# Use demo config instead
echo -e "${DIM}  (using demo/devboot.yaml)${RESET}"
echo ""
head -20 demo/devboot.yaml
echo -e "${DIM}  ... (truncated)${RESET}"
pause

# ── 6. Status — show installed vs configured ──────────────────────────────

header "6. devboot status — dashboard view"
$BINARY status demo/devboot.yaml
pause

# ── 7. Diff — preview what would change ───────────────────────────────────

header "7. devboot diff — preview changes before applying"
$BINARY diff demo/devboot.yaml
pause

# ── 8. Apply — dry run with --only git ─────────────────────────────────────

header "8. devboot apply --only git — apply just the git section"
$BINARY apply demo/devboot.yaml --only git --no-tui
pause

# ── 9. History — see what devboot did ──────────────────────────────────────

header "9. devboot history — audit trail"
$BINARY history
pause

# ── 10. Profile export — generate yaml from profile ───────────────────────

header "10. devboot profile export rust — export profile as YAML"
$BINARY profile export rust
pause

# ── Done ───────────────────────────────────────────────────────────────────

echo ""
echo -e "${BOLD}${CYAN}"
echo "  ┌─────────────────────────────────────────────────────┐"
echo "  │                                                     │"
echo "  │   Demo complete! Here's what you just saw:          │"
echo "  │                                                     │"
echo "  │   ✓ doctor    — environment diagnostics             │"
echo "  │   ✓ profile   — curated tool sets (7 profiles)      │"
echo "  │   ✓ export    — reverse-engineer current machine    │"
echo "  │   ✓ init      — interactive config wizard           │"
echo "  │   ✓ status    — installed vs configured dashboard   │"
echo "  │   ✓ diff      — preview changes before applying     │"
echo "  │   ✓ apply     — install everything from config      │"
echo "  │   ✓ history   — full audit trail                    │"
echo "  │   ✓ add       — interactive tool picker             │"
echo "  │   ✓ uninstall — rollback what devboot installed     │"
echo "  │                                                     │"
echo "  │   Get started:                                      │"
echo "  │     devboot init        # interactive wizard        │"
echo "  │     devboot apply       # set up everything         │"
echo "  │                                                     │"
echo "  └─────────────────────────────────────────────────────┘"
echo -e "${RESET}"
echo ""
