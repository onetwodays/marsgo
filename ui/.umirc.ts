import { defineConfig } from 'umi'; //写配置时也有提示，可以通过 umi 的 defineConfig 方法定义配置
import { ResponseError } from 'umi-request';

export default defineConfig({
  nodeModulesTransform: {
    type: 'none',
  },

  dva: {}, //使dva生效
  antd: {},
  

  //proxy 请求代理,如果把这段注释,就走mock形式 代理只是将请求服务做了中转，设置proxy不会修改请求地址。
  
  "proxy":{
    "/api/":{                              //1.设置了需要代理的请求头，比如这里定义了 /api ，当你访问如 /api/abc 这样子的请求，就会触发代理
      "target": "https://pvp.qq.com",      //2.设置代理的目标，即真实的服务器地址
      "changeOrigin": true,                //3.设置是否跨域请求资源             
      "pathRewrite": { "^/api/" : "" },    //4.表示是否重写请求地址，比如这里的配置，就是把 /api 替换成空字符           
    },
    "/apieos":{                                  //1.设置了需要代理的请求头，比如这里定义了 /api ，当你访问如 /api/abc 这样子的请求，就会触发代理
      "target": "http://127.0.0.1:18888",      //2.设置代理的目标，即真实的服务器地址
      "changeOrigin": true,                    //3.设置是否跨域请求资源             
      "pathRewrite": { "^/apieos" : "" },       //4.表示是否重写请求地址，比如这里的配置，就是把 /api 替换成空字符           
    },
  },
  




  //使用约定路由,注释掉routes
  //routes: [
  //  { path: '/', component: '@/pages/index' },
  //],
});

//配置文件，包含 umi 内置功能和插件的配置 Umi 在 .umirc.ts 或 config/config.ts 中配置项目和插件，支持 es6