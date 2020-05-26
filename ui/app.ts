//umi 运行时配置,

import {message} from 'antd'
import {ResponseError} from 'umi-request'
export const requset ={
    //prefix:'https://pvp.qq.com', //配置所有请求的prefix
    prefix:'/api', //请求代理的条件,当uri是/api开头时,请求发往.umirc.ts里面配置的代理里面
    errorHandler:(error:ResponseError)=>{
        message.error(error.message);
    },
}