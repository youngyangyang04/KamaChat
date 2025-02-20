const { defineConfig } = require('@vue/cli-service')
const fs = require('fs')
const path = require('path')
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  devServer: {
    host: '0.0.0.0',
    // 此处开启 https,并加载本地证书（否则浏览器左上角会提示不安全）
    https: {
      cert: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2.pem')),
      key: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2-key.pem')),
    },
    // 注意： https的端口必须是443
    port: 443,
  }
})
