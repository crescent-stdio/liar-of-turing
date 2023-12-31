import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {
        "liar-blue": "#3b82f6",
        "liar-blue-light": "#93c5fd",
        "liar-blue-lighter": "#d1e5fe",
        "liar-blue-dark": "#1e40af",
        "liar-blue-darker": "#18338c",
      },
    },
  },
  plugins: [],
};
export default config;
