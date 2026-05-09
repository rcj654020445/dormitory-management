#!/usr/bin/env node
/**
 * Validates that imports respect the Vue3 frontend layer hierarchy.
 * Pure Node.js ESM — no ts-node required.
 *
 * Layer 0: src/types/       — No internal imports allowed
 * Layer 1: src/api/, src/stores/ — May import types only
 * Layer 2: src/views/        — May import types, api, stores
 * Layer 3: src/components/   — May import types, api, stores, views
 *
 * Usage: node scripts/lint-deps.mjs
 */
import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Layer rules: pattern -> { layer, forbiddenImports }
const rules = [
  {
    pattern: "frontend/src/types/",
    layer: 0,
    allowedBases: ["../types/"], // only self-references
    note: "Layer 0: types must have NO internal imports",
  },
  {
    pattern: "frontend/src/api/",
    layer: 1,
    allowedBases: ["../types/", "./"], // can import types and self
    note: "Layer 1: api may import types",
  },
  {
    pattern: "frontend/src/stores/",
    layer: 2,
    allowedBases: ["../types/", "../api/", "./"],
    note: "Layer 2: stores may import types and api",
  },
  {
    pattern: "frontend/src/views/",
    layer: 3,
    allowedBases: ["../types/", "../api/", "../stores/", "./"],
    note: "Layer 3: views may import types, api, stores",
  },
  {
    pattern: "frontend/src/components/",
    layer: 4,
    allowedBases: ["../types/", "../api/", "../stores/", "../views/", "./"],
    note: "Layer 4: components may import anything except external @/",
  },
];

let violations = 0;

function resolveImport(baseFile, importPath) {
  // Resolve relative import to absolute path
  const baseDir = path.dirname(baseFile);
  if (importPath.startsWith("../")) {
    return path.resolve(baseDir, importPath);
  }
  if (importPath.startsWith("./")) {
    return path.resolve(baseDir, importPath);
  }
  return null; // absolute or @ alias
}

function checkFile(filePath) {
  if (
    filePath.includes("node_modules") ||
    filePath.includes(".test.") ||
    filePath.includes(".spec.")
  ) {
    return;
  }

  const ext = path.extname(filePath);
  if (ext !== ".ts" && ext !== ".vue" && ext !== ".tsx") return;

  let content = fs.readFileSync(filePath, "utf8");

  // For .vue files, extract only the <script> content
  if (ext === ".vue") {
    const scriptMatch = content.match(/<script[^>]*>([\s\S]*?)<\/script>/i);
    if (scriptMatch) {
      content = scriptMatch[1];
    } else {
      return; // No script block
    }
  }

  const lines = content.split("\n");

  for (const rule of rules) {
    if (!filePath.includes(rule.pattern)) continue;

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];

      // Match import statements
      const importMatch = line.match(/import\s+.+\s+from\s+['"]([^'"]+)['"]/);
      if (!importMatch) continue;

      const importPath = importMatch[1];

      // Skip external packages (@/, vue, etc.)
      if (importPath.startsWith("@/") || (!importPath.startsWith(".") && !importPath.startsWith("../"))) {
        continue;
      }

      // Skip self-references (importing same file)
      const resolved = resolveImport(filePath, importPath);
      if (resolved === filePath) continue;

      // Check if this import is in the allowed list for this layer
      const importBase = path.dirname(resolved || "");
      const allowed = rule.allowedBases.some((allowedBase) => {
        const allowedPath = path.resolve(__dirname, "..", allowedBase.replace("./", "frontend/src/"));
        return importBase.startsWith(allowedPath.replace(/\/$/, ""));
      });

      if (allowed) continue;

      // Find which layer this import belongs to
      let targetLayer = -1;
      let targetName = "";
      for (const r of rules) {
        if (importBase.includes(path.resolve(__dirname, "..", r.pattern.replace("frontend/src/", "frontend/src/")))) {
          targetLayer = r.layer;
          targetName = r.pattern;
          break;
        }
      }

      if (targetLayer < 0) continue; // external or unclassified

      // Check for reverse dependency (higher layer imports lower layer)
      if (targetLayer < rule.layer) {
        console.error(`\n[LINT-DEPS] ${filePath}:${i + 1}`);
        console.error(`  Layer ${rule.layer} (${rule.pattern}) imports Layer ${targetLayer} (${targetName}) — REVERSE DEPENDENCY`);
        console.error(`  ${line.trim()}`);
        console.error(`  Fix: ${rule.note}`);
        console.error(`  If you need data from ${targetName}, pass it as a prop or parameter.`);
        violations++;
      }
    }
  }
}

function walkDir(dir) {
  if (!fs.existsSync(dir)) return;
  if (dir.includes("node_modules")) return;

  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walkDir(fullPath);
    } else if (
      entry.name.endsWith(".ts") ||
      entry.name.endsWith(".vue") ||
      entry.name.endsWith(".tsx")
    ) {
      checkFile(fullPath);
    }
  }
}

const srcDir = path.join(__dirname, "..", "frontend", "src");
walkDir(srcDir);

if (violations === 0) {
  console.log("✅ All frontend imports follow the layer hierarchy");
  process.exit(0);
} else {
  console.error(`\n❌ Found ${violations} layer violation(s)`);
  process.exit(1);
}
