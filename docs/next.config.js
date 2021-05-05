const withNextra = require('nextra')('nextra-theme-docs', './theme.config.js')
module.exports = { basePath: '/docs', ...withNextra() }
