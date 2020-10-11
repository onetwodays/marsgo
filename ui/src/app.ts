/*
http://localhost:3001
在page下，定义_mock.js也可以使用mock功能。如./src/pages/index/_mock.js
 /web201605/js/herolist.json这个会变成/api/web201605/js/herolist.json
 mock/freeheros.json 这个会变成/apimock/freeheros.json
 由于/api/这个的前缀会走代理，而/apimock走不成代理，就走mock数据

在一个中后台中很多页面并不需要跨页面的信息流，也不需要把信息放入 model 中，
所以我们提供了 useRequest hooks，
useRequest 提供了一些快捷的操作和状态，可以大大的节省我们的代码

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

//import { createLogger } from 'redux-logger';
import { ResponseError } from 'umi-request';
import { notification } from 'antd';
import { RequestConfig } from 'umi';



const codeMessage: any = {
    200: '服务器成功返回请求的数据。',
    201: '新建或修改数据成功。',
    202: '一个请求已经进入后台排队（异步任务）。',
    204: '删除数据成功。',
    400: '发出的请求有错误，服务器没有进行新建或修改数据的操作。',
    401: '用户没有权限（令牌、用户名、密码错误）。',
    403: '用户得到授权，但是访问是被禁止的。',
    404: '发出的请求针对的是不存在的记录，服务器没有进行操作。',
    405: '请求方法不被允许。',
    406: '请求的格式不可得。',
    410: '请求的资源被永久删除，且不会再得到的。',
    422: '当创建一个对象时，发生一个验证错误。',
    500: '服务器发生错误，请检查服务器。',
    502: '网关错误。',
    503: '服务不可用，服务器暂时过载或维护。',
    504: '网关超时。',
};

const notification_error_style: any = {
    width: 600,
};

/**
 * 异常处理程序
 */
const errorHandler = (error: ResponseError) => {
    if (error.name === "BizError") {
        notification.error({
            message: `请求错误 ${error.data}`,
            //description: error.data.msg,
        });
        return error.data.code;
    }
    const { response } = error;
    if (response && response.status) {
        const errorText = codeMessage[response.status] || response.statusText;
        const { status, url } = response;

        notification.error({
            message: `请求错误 ${status}: ${url}`,
            description: errorText,
            style: notification_error_style,

        });
        return;
    }

    if (!response) {
        notification.error({
            description: '您的网络发生异常，无法连接服务器',
            message: '网络异常',
            style: notification_error_style,
        });
        return;
    }
    throw error;
};

/*
umi支持的错误格式
interface ErrorInfoStructure {
  success: boolean; // if request is success
  data?: any; // response data
  errorCode?: string; // code for errorType
  errorMessage?: string; // message display to user 
  showType?: number; // error display type： 0 silent; 1 message.warn; 2 message.error; 4 notification; 9 page
  traceId?: string; // Convenient for back-end Troubleshooting: unique request ID
  host?: string; // onvenient for backend Troubleshooting: host of current access server
}
自定义后端接口规范不满足时的适配

当 success 返回是 false 的情况我们会按照 showType 和 errorMessage 
来做统一的错误提示，同时抛出一个异常，异常的格式如下
interface RequestError extends Error {
  data?: any; // 这里是后端返回的原始数据
  info?: ErrorInfoStructure;
}
可以通过 Error.name 是否是 BizError 来判断是否是因为 success 为 false 抛出的错误

*/
const errorConfig: any = {
    adaptor: (res: any) => {
        return {
            success: res.code == 200, //success是false时
            data: res.data,
            errorCode: res.code,
            errorMessage: res.msg, //显示的错误信息
        };
    },
    //showType: 9, //当 showType 为 9 时，默认会跳转到 /exception 页面，你可以通过该配置来修改该路径。

};

/*a1 -> b1 -> response -> b2 -> a2*/
const middlewareA = async (ctx: any, next: Function) => {
    console.log('A before');
    await next();
    console.log('A after');

}

const middlewareB = async (ctx: any, next: Function) => {
    console.log('B before');
    await next();
    console.log('B after');

}


/**
 * 配置request请求时的默认参数
 * 使用app.ts配置RequestConfig就 不能使用extend来配置,不然 errorConfig.adaptor 不会起作用
 * 
 * 该配置返回一个对象。除了 errorConfig 和 middlewares 以外其它配置都是直接透传 umi-request 的全局配置。
 * umi-request 提供中间件机制，之前是通过 request.use(middleware) 的方式引入，现在可以通过 request.middlewares 进行配置。
 * requestInterceptors 该配置接收一个数组，数组的每一项为一个 request 拦截器。等同于 umi-request 的 request.interceptors.request.use()。具体见 umi-request 的拦截器文档
 * responseInterceptors该配置接收一个数组，数组的每一项为一个 response 拦截器。等同于 umi-request 的 request.interceptors.response.use()。具体见 umi-request 的拦截器文档
 */
export const request: RequestConfig = {
    //prefix: 'https://pvp.qq.com',这样会有跨域的问题
    prefix: '/api', //所有的请求的前缀,相当于ip+port部分
    errorHandler: errorHandler,//默认错误处理

    credentials: 'include', //默认请求是否带上cookie
    timeout: 5000,

    errorConfig: {}, //自定义接口规范
    middlewares: [middlewareA, middlewareB],

    requestInterceptors: [],
    responseInterceptors: [],



};






const plugins = [];
// 非生产环境添加 logger
if (process.env.NODE_ENV !== "production") {
    plugins.push(
        require("dva-logger")({
            collapsed: true
        })
    );
}


export const dva = {
    config: {
        //onAction: createLogger(),
        onError(e: Error) {
            e.preventDefault();
            notification.error(e.message);
        },
    },
    plugins,
};



