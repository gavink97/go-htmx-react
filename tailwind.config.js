/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors')
// const defaultTheme = require('tailwindcss/defaultTheme');
// const plugin = require('tailwindcss/plugin')

module.exports = {
  content: [
    './src/**/*.{ts,tsx}',
    './internal/**/*.templ'
  ],
theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        mobile: "2rem",
        tablet: "4rem",
        desktop: "5rem",
      },
    },
    extend: {
      colors: {
        primary: colors.blue,
        secondary: colors.yellow,
        neutral: colors.gray,
      }
    },
  },
  plugins: [
        require("tailwindcss-animate"),
        require('@tailwindcss/forms'),
        require('@tailwindcss/container-queries')
    ],
}
