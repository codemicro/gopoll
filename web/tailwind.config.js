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
      colors: {
        theme: {
          DEFAULT: '#233D4D',
          highlight: {
            DEFAULT: '#FFAE03',
            lighter: "#ffb740",
            darker: "#E59C02"
          },
        },
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
