/**
 * 在生产环境 代理是无法生效的，所以这里没有生产环境的配置
 * The agent cannot take effect in the production environment
 * so there is no configuration of the production environment
 * For details, please see
 * https://pro.ant.design/docs/deploy
 */
export default {
    dev: {

        "/api/": {                              //1.设置了需要代理的请求头，比如这里定义了 /api ，当你访问如 /api/abc 这样子的请求，就会触发代理
            "target": "https://pvp.qq.com",      //2.设置代理的目标，即真实的服务器地址
            "changeOrigin": true,                //3.设置是否跨域请求资源             
            "pathRewrite": { "^/api/": "" },    //4.表示是否重写请求地址，比如这里的配置，就是把 /api 替换成空字符           
        },
        "/apieos": {                                  //1.设置了需要代理的请求头，比如这里定义了 /api ，当你访问如 /api/abc 这样子的请求，就会触发代理
            "target": "http://127.0.0.1:18888",      //2.设置代理的目标，即真实的服务器地址
            "changeOrigin": true,                    //3.设置是否跨域请求资源             
            "pathRewrite": { "^/apieos": "" },       //4.表示是否重写请求地址，比如这里的配置，就是把 /api 替换成空字符           
        },
    },
    test: {
        '/api/': {
            target: 'https://preview.pro.ant.design',
            changeOrigin: true,
            pathRewrite: { '^': '' },
        },
    },
    pre: {
        '/api/': {
            target: 'your pre url',
            changeOrigin: true,
            pathRewrite: { '^': '' },
        },
    },
};
