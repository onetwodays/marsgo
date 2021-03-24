// https://github.com/umijs/umi-request/blob/master/README_zh-CN.md#%E4%B8%AD%E9%97%B4%E4%BB%B6

import { notification } from 'antd';

export const eosResponseInterceptor = async (res: Response, req) => {
    //克隆响应对象做解析处理
    if (res.url.includes('/v1/chain')) {
        let temp = await res.clone().json();
        console.log("响应拦截器输出:", temp); //temp 是服务端返回的http body对象
        //console.log("响应拦截器输出res:", res); //temp 是request库封装的响应体对象
        //console.log("响应拦截器输出req:", req); //temp 是request库封装的请求体对象

        //eos的错误信息拦截，code对应http status，message是status的文字解释，error是具体错误
        const { code, message, error } = temp;
        if (code && code != 200) {
            notification.error({
                message: `EOS请求错误:${res.url} ${code}:${message}`,
                description: JSON.stringify(temp, null, 4),
            });
            //const err = new Error(`{error}`);
            //throw err;
        }
    }

    return res;
}
