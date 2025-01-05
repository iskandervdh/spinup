/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#d8fe05',
        'primary-dark': '#b0cf01',
        background: '#575757',

        info: '#5dade2',
        'info-dark': '#3297da',
        success: '#88d567',
        warning: '#ffd700',
        error: '#c94b55',
      },
      fontFamily: {
        azeret: ['Azeret'],
      },
    },
  },
  plugins: [],
};
