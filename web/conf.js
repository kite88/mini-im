/**
 * 后端接口地址配置
 *
 * 默认与页面同源，本地开发 / 线上部署都无需修改。
 * 若前后端分离部署，直接改成后端地址即可，例如：
 *   const apiDomain = 'http://127.0.0.1:2580'
 *   const wsDomain  = 'ws://127.0.0.1:2580'
 */
const apiDomain = window.location.origin
const wsDomain = (window.location.protocol === 'https:' ? 'wss://' : 'ws://') + window.location.host
