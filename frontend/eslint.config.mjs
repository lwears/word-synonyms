import globals from 'globals'
import pluginJs from '@eslint/js'
import tseslint from 'typescript-eslint'
import eslintPluginUnicorn from 'eslint-plugin-unicorn'
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended'
import tsParser from '@typescript-eslint/parser'
import reactRefresh from 'eslint-plugin-react-refresh'

export default [
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  eslintPluginUnicorn.configs['flat/recommended'],
  eslintPluginPrettierRecommended,
  {
    ignores: [
      '**/dist',
      '**/eslintrc_old.cjs',
      '**/postcss.config.js',
      '**/vite-env.d.ts',
    ],
  },
  {
    files: ['**/*.{js,mjs,cjs,ts,jsx,tsx}'],
  },
  {
    languageOptions: { globals: globals.browser, parser: tsParser },
    plugins: {
      'react-refresh': reactRefresh,
    },
    rules: {
      'react-refresh/only-export-components': [
        'warn',
        {
          allowConstantExport: true,
        },
      ],
      '@typescript-eslint/no-unused-vars': 'warn',
      '@typescript-eslint/consistent-type-definitions': ['error', 'interface'],
      '@typescript-eslint/consistent-type-imports': 'warn',
      '@typescript-eslint/no-namespace': 'off',
      'unicorn/filename-case': 'off',
      'unicorn/prevent-abbreviations': 'off',
      'unicorn/numeric-separators-style': 'off',
      'unicorn/no-null': 'off',
    },
  },
]
