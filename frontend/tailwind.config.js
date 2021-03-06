const colors = require('tailwindcss/colors');

module.exports = {
  mode: 'jit',
  content: ['./src/**/*.{ts,tsx}'],
  darkMode: 'class',
  theme: {
    colors,
    extend: {
      colors: {
        primary: colors.rose,
        'text-primary': colors.gray['800'],
        'text-primary-dark': colors.gray['300'],
        'text-secondary': colors.gray['500'],
        'text-secondary-dark': colors.gray['400'],
        rakuten: '#BF0000',
        'yahoo-shopping': '#FF0132',
        'paypay-mall': '#977a20',
        brown: {
          50: '#fdf8f6',
          100: '#f2e8e5',
          200: '#eaddd7',
          300: '#e0cec7',
          400: '#d2bab0',
          500: '#bfa094',
          600: '#a18072',
          700: '#977669',
          800: '#846358',
          900: '#43302b',
        },
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
  plugins: [
    require('@tailwindcss/line-clamp'),
    require('@tailwindcss/forms'),
    function ({ addUtilities }) {
      const newUtilities = {
        '.text-shadow': {
          textShadow: '3px 3px 3px #808080',
        },
      };

      addUtilities(newUtilities);
    },
  ],
};
