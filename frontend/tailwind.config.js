/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#f0f9ff",
          500: "#0284c7",
          900: "#0c4a6e",
        },
      },
      minHeight: {
        touch: "44px", // WCAG Target Size
      },
      minWidth: {
        touch: "44px",
      },
    },
  },
  plugins: [],
};
