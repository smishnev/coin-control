module.exports = {
  darkMode: 'class',
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        secondary: 'var(--color-secondary)',
        background: 'var(--color-background)',
        foreground: 'var(--color-foreground)',
        menu: 'var(--color-menu)',
        'menu-active': 'var(--color-menu-active)',
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
      keyframes: {
    'fade-in': { '0%': { opacity: 0 }, '100%': { opacity: 1 } }
      },
      animation: {
        'fade-in': 'fade-in 0.15s ease'
      }
    },
  },
  plugins: [],
};