import type { Config } from 'tailwindcss'

const config: Config = {
	content: [
		'./pages/**/*.{js,ts,jsx,tsx,mdx}',
		'./components/**/*.{js,ts,jsx,tsx,mdx}',
		'./app/**/*.{js,ts,jsx,tsx,mdx}',
	],
	theme: {
		extend: {
			backgroundImage: {
			'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
			'gradient-conic':
				'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
			},
			colors: {
				'palette': {
					100: "#E4EEE7",
					200: '#060B07',
					300: '#9BD7AE',
					400: '#51FC84',
					500: '#3ED46C',
					600: '#46464A'
				},
			},
			fontFamily: {
				Roboto: ["Roboto Mono", "monospace"],
				Display: ["Red Hat Display", "sans-serif"],
			},
		},
	},
	plugins: [],
}

export default config
