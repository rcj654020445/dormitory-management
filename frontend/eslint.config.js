// @ts-check
import eslintPluginVue from 'eslint-plugin-vue'
import vueParser from 'vue-eslint-parser'
import tseslint from 'typescript-eslint'

// Vue SFC <script setup lang="ts"> files
const vueTsConfig = {
  files: ['**/*.vue'],
  processor: eslintPluginVue.processors['.vue'],
  plugins: {
    vue: eslintPluginVue,
  },
  languageOptions: {
    parser: vueParser,
    parserOptions: {
      parser: tseslint.parser,
      ecmaVersion: 'latest',
      sourceType: 'module',
    },
  },
  rules: {
    ...eslintPluginVue.configs['flat/essential'][0].rules,
    'vue/multi-word-component-names': 'off',
    'no-unused-vars': 'off',
    'no-console': 'off',
  },
}

// TypeScript files (.ts / .tsx)
const tsConfig = {
  files: ['**/*.ts', '**/*.tsx', '**/*.mts', '**/*.cts'],
  plugins: {
    '@typescript-eslint': tseslint.plugin,
  },
  languageOptions: {
    parser: tseslint.parser,
    parserOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
    },
  },
  rules: {
    ...tseslint.configs.recommended[0].rules,
    '@typescript-eslint/no-explicit-any': 'warn',
    '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
  },
}

// JS files
const jsConfig = {
  files: ['**/*.js', '**/*.jsx', '**/*.mjs', '**/*.cjs'],
  languageOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    globals: {
      console: 'readonly',
      window: 'readonly',
      document: 'readonly',
      navigator: 'readonly',
      NodeJS: 'readonly',
    },
  },
  rules: {
    'no-unused-vars': 'off',
    'no-console': 'off',
  },
}

export default tseslint.config(
  { ignores: ['dist/**', 'node_modules/**', '*.min.js'] },
  vueTsConfig,
  tsConfig,
  jsConfig,
)
