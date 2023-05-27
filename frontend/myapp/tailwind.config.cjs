const colors = require('tailwindcss/colors')

const config = {
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    "./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}",
  ],

  theme: {
    colors: {
      'tahiti': {
        light: '#67e8f9',
        DEFAULT: '#06b6d4',
        dark: '#0e7490',
      },
      'redTest': "#c4b5fd",
      'testYellow':{
        light: 'text-yellow-100',
        dark:  'text-yellow-100',
      },
    },
    
  },

  plugins: [
    require('flowbite/plugin')
  ],
  darkMode: 'class',
};

module.exports = config;