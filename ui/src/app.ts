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

import { ResponseError } from 'umi-request';

export const request = {
    //prefix: 'https://pvp.qq.com',这样会有跨域的问题
    prefix:'/api', //所有的请求都走umirc.ts里面的代理/api
    errorHandler:(error:ResponseError)=>{
        console.log(error);
    },
    

};
