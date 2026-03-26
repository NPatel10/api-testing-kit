import js from "@eslint/js";
import tsPlugin from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";
import { defineConfig } from "eslint/config";
import globals from "globals";
import svelte from "eslint-plugin-svelte";
import svelteParser from "svelte-eslint-parser";

const sharedGlobals = {
	...globals.browser,
	...globals.node,
};

const svelteRuneGlobals = {
	$derived: "readonly",
	$effect: "readonly",
	$state: "readonly",
};

export default defineConfig([
	{
		ignores: [".svelte-kit/**", "build/**", "dist/**", "coverage/**", "node_modules/**"],
	},
	{
		linterOptions: {
			reportUnusedDisableDirectives: "off",
		},
	},
	js.configs.recommended,
	...svelte.configs["flat/recommended"],
	{
		files: ["**/*.{js,mjs,cjs,ts,tsx}"],
		languageOptions: {
			globals: sharedGlobals,
		},
	},
	{
		files: ["**/*.{ts,tsx}"],
		languageOptions: {
			parser: tsParser,
			parserOptions: {
				ecmaVersion: "latest",
				sourceType: "module",
			},
		},
		plugins: {
			"@typescript-eslint": tsPlugin,
		},
		rules: {
			"no-unused-vars": "off",
			"@typescript-eslint/no-unused-vars": [
				"error",
				{
					argsIgnorePattern: "^_",
					varsIgnorePattern: "^_",
				},
			],
		},
	},
	{
		files: ["**/*.svelte.{js,ts}"],
		languageOptions: {
			globals: {
				...sharedGlobals,
				...svelteRuneGlobals,
			},
			parser: tsParser,
			parserOptions: {
				ecmaVersion: "latest",
				sourceType: "module",
			},
		},
		plugins: {
			"@typescript-eslint": tsPlugin,
		},
		rules: {
			"no-unused-vars": "off",
			"@typescript-eslint/no-unused-vars": [
				"error",
				{
					argsIgnorePattern: "^_",
					varsIgnorePattern: "^_",
				},
			],
		},
	},
	{
		files: ["**/*.svelte"],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				parser: tsParser,
				extraFileExtensions: [".svelte"],
				ecmaVersion: "latest",
				sourceType: "module",
			},
			globals: {
				...sharedGlobals,
				...svelteRuneGlobals,
			},
		},
		plugins: {
			"@typescript-eslint": tsPlugin,
		},
		rules: {
			"no-unused-vars": "off",
			"@typescript-eslint/no-unused-vars": [
				"error",
				{
					argsIgnorePattern: "^_",
					varsIgnorePattern: "^_",
				},
			],
		},
	},
]);
