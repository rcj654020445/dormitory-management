#!/usr/bin/env ts-node
/**
 * Validates Vue3 frontend quality rules:
 * - File size limits (max 500 lines for .ts, 300 lines for .vue)
 * - No console.log in non-test files (use a logger)
 * - No any type without comment
 *
 * Usage: npx ts-node scripts/lint-quality.ts
 */
import * as fs from "fs";
import * as path from "path";

const MAX_FILE_LINES_TS = 500;
const MAX_FILE_LINES_VUE = 300;
let violations = 0;

function checkFile(filePath: string): void {
  const isVue = filePath.endsWith(".vue");
  const isTest = filePath.includes(".test.") || filePath.includes(".spec.");
  const isMain = filePath.endsWith("main.ts") || filePath.endsWith("main.js");
  const maxLines = isVue ? MAX_FILE_LINES_VUE : MAX_FILE_LINES_TS;

  // Skip node_modules
  if (filePath.includes("node_modules")) return;

  const content = fs.readFileSync(filePath, "utf8");
  const lines = content.split("\n");

  // Check file size
  if (lines.length > maxLines) {
    console.error(
      `✗ ${filePath} has ${lines.length} lines (max ${maxLines})`
    );
    console.error(
      `  Fix: Split this file into smaller, focused modules.`
    );
    violations++;
  }

  // Check for console.log in non-test, non-main files
  if (!isTest && !isMain) {
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const trimmed = line.trim();

      // Skip comments
      if (trimmed.startsWith("//") || trimmed.startsWith("/*") || trimmed.startsWith("*")) {
        continue;
      }

      // Check for console.log
      if (line.includes("console.log(") || line.includes("console.error(")) {
        console.error(
          `✗ ${filePath}:${i + 1} — Use structured logger instead of console.log`
        );
        console.error(
          `  Fix: Import a logger from src/utils/logger.ts or use Vue's built-in logging.`
        );
        violations++;
      }

      // Check for 'any' type without comment
      if (line.includes(": any") || line.includes("<any>") || line.includes("as any")) {
        const hasComment = i > 0 && (lines[i - 1].includes("//") || lines[i - 1].includes("/*"));
        if (!hasComment) {
          console.error(
            `✗ ${filePath}:${i + 1} — Avoid 'any' type. Use explicit types or add a comment explaining why any is necessary.`
          );
          violations++;
        }
      }
    }
  }
}

function walkDir(dir: string): void {
  if (!fs.existsSync(dir)) return;
  if (dir.includes("node_modules")) return;

  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walkDir(fullPath);
    } else if (
      (entry.name.endsWith(".ts") && !entry.name.endsWith(".d.ts")) ||
      entry.name.endsWith(".vue")
    ) {
      checkFile(fullPath);
    }
  }
}

walkDir("frontend/src");

if (violations === 0) {
  console.log("✓ All quality checks passed");
  process.exit(0);
} else {
  console.error(`\n✗ Found ${violations} quality violation(s)`);
  process.exit(1);
}
