/** @type {import('tailwindcss').Config} */

const defaultTheme = require("tailwindcss/defaultTheme");
// const fontInter = require("tailwindcss-font-inter")(defaultTheme.textSizes);
module.exports = {
  content: ["./views/**/*.{html,js}", "./src/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [fontInter],
};
