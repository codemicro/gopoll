module.exports = {
  purge: [
    './templates/*.qtpl',
    './templates/parts/*.qtpl',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      spacing: {
        128: '32rem',
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
