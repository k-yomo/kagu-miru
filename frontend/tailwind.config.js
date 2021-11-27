const colors = require('tailwindcss/colors');

module.exports = {
  mode: 'jit',
  purge: ['./src/**/*.{ts,tsx}'],
  darkMode: 'class', // or 'media' or 'class'
  theme: {
    colors,
    extend: {
      colors: {
        primary: colors.pink,
        'text-primary': colors.gray['800'],
        'text-primary-dark': colors.gray['100'],
        'text-secondary': colors.gray['500'],
        'text-secondary-dark': colors.gray['300'],
        rakuten: '#BF0000',
        'yahoo-shopping': '#FF0132',
      },
      fontFamily: {
        sans: [
          '游ゴシック体',
          'YuGothic',
          '游ゴシック',
          'Yu Gothic',
          'sans-serif',
        ],
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [require('@tailwindcss/line-clamp')],
};
