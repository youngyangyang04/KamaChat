const { defineConfig } = require('@vue/cli-service')
// const fs = require('fs')
// const path = require('path')
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  // 本地开发环境 - 使用 HTTP
  devServer: {
    host: '0.0.0.0',
    // 端口设为常见的开发用端口，避免与服务器默认的 HTTPS 端口冲突
    port: 8080,
  }
  // HTTPS 配置（生产环境）
  // devServer: {
  //   host: '0.0.0.0',
  //   https: {
  //     // Win10本地部署
  //     // cert: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2.pem')),
  //     // key: fs.readFileSync(path.join(__dirname, 'src/assets/cert/127.0.0.1+2-key.pem')),
  //     // Ubuntu22.04云服务器部署
  //     // cert: fs.readFileSync(path.join("/etc/ssl/certs/server.crt")),
  //     // key: fs.readFileSync(path.join("/etc/ssl/private/server.key")),
  //   },
  //   port: 443,
  // }
})
