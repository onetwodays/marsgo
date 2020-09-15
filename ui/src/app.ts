//http://localhost:3001
//在page下，定义_mock.js也可以使用mock功能。如./src/pages/index/_mock.js
/*
├── package.json
├── .umirc.ts
├── .env
├── dist
├── mock
├── public
└── src
    ├── .umi
    ├── layouts/index.tsx
    ├── pages
        ├── index.less
        └── index.tsx
    └── app.ts
*/
//网络请求的运行时配置
import { createLogger }    from 'redux-logger';
import { ResponseError }   from 'umi-request';
import { message }         from 'antd';
import {RequestConfig}     from 'umi';


export const dva = {
    config: {
      onAction: createLogger(),
      onError(e: Error) {
        message.error(e.message, 3);
      },
    },
};

//所有请求的 prefix
export const request:RequestConfig= {
    //prefix: 'https://pvp.qq.com',这样会有跨域的问题
    prefix:'/api', //所有的请求的前缀,相当于ip+port部分
    errorHandler:(error:ResponseError)=>{
        //console.log(error);
        message.error(error.request);
        message.error(error);
        
    },
    timeout:3000,
    errorConfig:{},
    middlewares:[ 
        async function middlewareA(ctx, next) {
            //console.log('A before');
            await next();
            //console.log('A after');
       },
       async function middlewareB(ctx, next) {
            //console.log('B before');
            await next();
            //console.log('B after');
       },
    ],
    requestInterceptors:[],
    responseInterceptors:[],

    

};
// /web201605/js/herolist.json这个会变成/api/web201605/js/herolist.json
// mock/freeheros.json 这个会变成/apimock/freeheros.json
// 由于/api/这个的前缀会走代理，而/apimock走不成代理，就走mock数据