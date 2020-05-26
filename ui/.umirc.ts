import { defineConfig } from 'umi';
import { ResponseError } from 'umi-request';

export default defineConfig({
  nodeModulesTransform: {
    type: 'none',
  },

  dva: {}, //使dva生效
  antd: {},

  //proxy 请求代理,如果把这段注释,就走mock形式
  "proxy":{
    "/api/":{ //设置了需要代理的请求头，比如这里定义了 /api ，当你访问如 /api/abc 这样子的请求，就会触发代理
      "target": "https://pvp.qq.com",  //设置代理的目标，即真实的服务器地址
      "changeOrigin": true,           //设置是否跨域请求资源             
      "pathRewrite": { "^/api" : "" },  //表示是否重写请求地址，比如这里的配置，就是把 /api 替换成空字符           
    },
  }




  //使用约定路由,注释掉routes
  //routes: [
  //  { path: '/', component: '@/pages/index' },
  //],
});
