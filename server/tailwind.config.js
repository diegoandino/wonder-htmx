/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './**/*.html', // Include all HTML files in your project
    './**/*.templ', // Include Templ files if you're using a custom templating engine
    './app/**/*.go', // Consider including Go HTML templates if you're using inline templates
    './static/js/**/*.js', // Include JavaScript files for dynamic class additions
    './handlers/**/*.go',
  ],
  theme: {
    extend: {
      colors: {
        // Example custom colors
        'custom-blue': '#007bff',
        'custom-gray': '#6c757d',
      },
      spacing: {
        // Example custom spacing
        '72': '18rem',
        '84': '21rem',
        '96': '24rem',
      },
      borderRadius: {
        // Example custom border radius
        'xl': '1rem',
      },
      // Include any other theme extensions you might need
    },
    // Consider customizing other theme values or adding new ones as per your design requirements
  },
  plugins: [
    
  ],
};

