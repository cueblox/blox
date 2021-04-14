const withPWA = require('next-pwa');

module.exports = withPWA({
  images: {
    domains: ['res.cloudinary.com'],
  },
  future: { webpack5: true },
  pwa: {
    dest: 'public',
    disable: process.env.NODE_ENV === 'development',
  }
})
