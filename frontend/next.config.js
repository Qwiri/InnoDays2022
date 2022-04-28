/** @type {import('next').NextConfig} */

const path = require('path')
const withSass = require('@zeit/next-sass');
module.exports = withSass({
  cssModules: true
})
module.exports = {
  sassOptions: {
  includePaths: [path.join('./', 'styles')],
  },
}
const nextConfig = {
  reactStrictMode: true,
  images: {
    domains: ['upload.wikimedia.org',
              'avatars.dicebear.com'],
  },
}

module.exports = nextConfig
