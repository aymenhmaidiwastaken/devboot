#!/usr/bin/env node
// Records a compelling devboot demo → demo.cast → demo.gif
// Story: show config → diff → apply → status
// Usage: node demo/record.js

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const ROOT = path.resolve(__dirname, "..");
const BIN = path.join(ROOT, "devboot.exe");
const AGG = path.join(
  process.env.HOME || process.env.USERPROFILE,
  ".local", "bin", "agg.exe"
);
const CAST = path.join(__dirname, "demo.cast");
const GIF = path.join(__dirname, "demo.gif");

function run(args) {
  try {
    return execSync(`"${BIN}" ${args}`, { cwd: ROOT, encoding: "utf-8", timeout: 30000 });
  } catch (e) { return e.stdout || ""; }
}

function esc(s) {
  let r = "";
  for (let i = 0; i < s.length; i++) {
    const c = s.charCodeAt(i);
    if (c === 0x5c) r += "\\\\";
    else if (c === 0x22) r += '\\"';
    else if (c === 0x09) r += "\\t";
    else if (c === 0x0a) r += "\\n";
    else if (c === 0x0d) continue;
    else if (c === 0x1b) r += "\\u001b";
    else if (c < 0x20) r += "\\u" + c.toString(16).padStart(4, "0");
    else r += s[i];
  }
  return r;
}

// ── Capture real outputs ──
console.log("Capturing...");

const configYaml = fs.readFileSync(path.join(__dirname, "devboot.yaml"), "utf-8");
const outDiff = run("diff demo/devboot.yaml");
const outApply = run("apply demo/devboot.yaml --only git --no-tui");
const outStatus = run("status demo/devboot.yaml");

// ── Build cast ──
console.log("Building .cast...");

const events = [];
let t = 0.0;

function emit(data) {
  events.push(`[${t.toFixed(4)}, "o", "${data}"]`);
}

function prompt() {
  t += 0.5;
  emit("\\u001b[1;32m❯\\u001b[0m ");
}

function typeCmd(text) {
  for (const ch of text) {
    t += 0.04 + Math.random() * 0.03;
    emit(esc(ch));
  }
}

function enter() {
  t += 0.08;
  emit("\\r\\n");
}

// Output lines slowly — the key fix
// delay = seconds between each line
function outSlow(text, delay) {
  const lines = text.split("\n");
  for (const line of lines) {
    t += delay;
    emit(esc(line) + "\\r\\n");
  }
}

function pause(s) {
  t += s;
  emit("");
}

function comment(text) {
  t += 0.4;
  emit("\\u001b[2;3m# " + esc(text) + "\\u001b[0m\\r\\n");
}

function clearScreen() {
  t += 0.1;
  emit("\\u001b[2J\\u001b[H");
}

// ═══════════════════════════════════════════════
// Scene 1: Show the config file
// ═══════════════════════════════════════════════
prompt();
comment("Your entire dev environment in one file:");
pause(0.6);
prompt();
typeCmd("cat devboot.yaml");
enter();
pause(0.3);
outSlow(configYaml, 0.06);  // 60ms per line — readable scroll
pause(4.0);                  // hold so people can read the end

// ═══════════════════════════════════════════════
// Scene 2: Diff — preview changes
// ═══════════════════════════════════════════════
clearScreen();
prompt();
comment("See what would change before applying:");
pause(0.6);
prompt();
typeCmd("devboot diff");
enter();
pause(0.3);
outSlow(outDiff, 0.08);     // 80ms per line
pause(4.0);

// ═══════════════════════════════════════════════
// Scene 3: Apply — the main event
// ═══════════════════════════════════════════════
clearScreen();
prompt();
comment("One command. Everything configured:");
pause(0.6);
prompt();
typeCmd("devboot apply");
enter();
pause(0.4);
outSlow(outApply, 0.12);    // 120ms per line — feels like real work
pause(3.5);

// ═══════════════════════════════════════════════
// Scene 4: Status — proof it worked
// ═══════════════════════════════════════════════
clearScreen();
prompt();
typeCmd("devboot status");
enter();
pause(0.3);
outSlow(outStatus, 0.08);
pause(4.0);

// Final
prompt();
pause(2.0);

// Write cast file
const header = JSON.stringify({
  version: 2,
  width: 90,
  height: 35,
  timestamp: Math.floor(Date.now() / 1000),
  env: { SHELL: "/bin/bash", TERM: "xterm-256color" },
  title: "DevBoot — Fresh machine to productive in one command",
});

fs.writeFileSync(CAST, header + "\n" + events.join("\n") + "\n");
console.log(`Cast: ${events.length} events, ${t.toFixed(1)}s`);

// Render GIF
console.log("Rendering GIF...");
try {
  execSync(
    `"${AGG}" "${CAST}" "${GIF}" --theme dracula --font-size 16`,
    { cwd: ROOT, stdio: "inherit", timeout: 300000 }
  );
  const size = (fs.statSync(GIF).size / 1024).toFixed(0);
  console.log(`\nDone! demo/demo.gif (${size} KB)`);
} catch (e) {
  console.log("agg error:", e.message);
}
