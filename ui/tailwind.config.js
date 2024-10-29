/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/react-tailwindcss-datepicker/dist/index.esm.{js,ts}",
  ],
  theme: {
    extend: {
      colors: {
        'geopost-black': '#414042',
        'geopost-midgray': '#CAC4BE',
        'geopost-red': '#DC0032',
        'geopost-blue': '#4F46E5',
        'geopost-darkblue': '#342CDC',
        'geopost-lightgray': '#D9D9D9',
        'geopost-lightergray': '#E6E7E8',
        'geopost-darkgray': '#4D4646',
        'geopost-darkred': '#A90034',
        'geopost-green': '#78A55A',
      },
      boxShadow: {
        'organization-card': '0px 2px 5px 0px rgba(0, 0, 0, 0.25)',
        'project-card': '0px 4px 4px rgba(0, 0, 0, 0.25)',
      },
    },
  },
  plugins: [],
}

