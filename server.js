/**
 * Root process launcher for pxxl.app Node.js environment fallback
 * Spawns the native Curexal Go backend executable.
 */
import { spawn } from "child_process";
import path from "path";
import fs from "fs";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Locate binary
const binCandidates = [
  path.join(__dirname, "apps", "api", "bin", "curexal-backend"),
  path.join(__dirname, "apps", "api", "bin", "curexal-backend.exe"),
  path.join(__dirname, "bin", "curexal-backend"),
  path.join(__dirname, "curexal-backend"),
];

let binPath = binCandidates.find((p) => fs.existsSync(p));

if (!binPath) {
  console.log("==> Building Curexal Go backend on startup...");
  const { execSync } = await import("child_process");
  execSync("cd apps/api && go build -ldflags=\"-w -s\" -o bin/curexal-backend ./cmd/CUREXAL", {
    stdio: "inherit",
  });
  binPath = binCandidates.find((p) => fs.existsSync(p));
}

if (!binPath) {
  console.error("Could not find or build curexal-backend binary.");
  process.exit(1);
}

console.log(`==> Launching Curexal Backend from: ${binPath}`);
const child = spawn(binPath, [], {
  cwd: path.join(__dirname, "apps", "api"),
  env: process.env,
  stdio: "inherit",
});

child.on("exit", (code, signal) => {
  console.log(`Curexal backend process exited with code ${code}, signal ${signal}`);
  process.exit(code || 0);
});
