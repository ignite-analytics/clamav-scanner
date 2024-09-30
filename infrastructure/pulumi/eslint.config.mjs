import eslint from '@eslint/js'
// import globals from 'globals'
import tseslint from 'typescript-eslint'

export default tseslint.config({
	extends: [eslint.configs.recommended, ...tseslint.configs.recommended],
	files: ['**/*.ts'],
	ignores: ['dist', 'node_modules'],
	languageOptions: {
		ecmaVersion: 'latest'
	},
	plugins: {
		// Add any plugins here
	},
	rules: {
		// Add any rules here
	}
})
