const { defineConfig } = require('@vue/cli-service')
const fs = require('fs')
const path = require('path')

function parseBoolean(value, defaultValue = false) {
  if (value === undefined || value === null || value === '') {
    return defaultValue
  }

  switch (String(value).trim().toLowerCase()) {
    case '1':
    case 'true':
    case 'yes':
    case 'on':
      return true
    case '0':
    case 'false':
    case 'no':
    case 'off':
      return false
    default:
      throw new Error(`Invalid boolean value: ${value}`)
  }
}

function resolveFilePath(filePath) {
  return path.isAbsolute(filePath) ? filePath : path.resolve(__dirname, filePath)
}

const devServer = {
  host: process.env.DEV_SERVER_HOST || '0.0.0.0',
  port: Number.parseInt(process.env.DEV_SERVER_PORT || '8080', 10),
}

if (parseBoolean(process.env.DEV_SERVER_HTTPS, false)) {
  const certPath = process.env.DEV_SERVER_HTTPS_CERT
  const keyPath = process.env.DEV_SERVER_HTTPS_KEY

  if (!certPath || !keyPath) {
    throw new Error(
      'DEV_SERVER_HTTPS=true requires both DEV_SERVER_HTTPS_CERT and DEV_SERVER_HTTPS_KEY'
    )
  }

  devServer.https = {
    cert: fs.readFileSync(resolveFilePath(certPath)),
    key: fs.readFileSync(resolveFilePath(keyPath)),
  }
}

module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  chainWebpack: (config) => {
    config.plugin('html').tap((args) => {
      args[0].title = 'GoChat'
      return args
    })
  },
  devServer,
})
