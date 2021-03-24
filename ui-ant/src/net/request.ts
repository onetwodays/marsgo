import { notification } from 'antd';
import type { ResponseError } from "umi-request";

const codeMessage = {
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

const NotificationErrorStyle: any = {
    width: 600,
};

/*
export interface ResponseError<D = any> extends Error {
  name: string;
  data: D;
  response: Response;
  request: {
    url: string;
    options: RequestOptionsInit;
  };
  type: string;
}

interface Response extends Body {
    readonly headers: Headers;
    readonly ok: boolean;
    readonly redirected: boolean;
    readonly status: number;
    readonly statusText: string;
    readonly trailer: Promise<Headers>;
    readonly type: ResponseType;
    readonly url: string;
    clone(): Response;
}

errorConfig?: {
    errorPage?: string;
    adaptor?: (resData: any, ctx: Context) => ErrorInfoStructure;
  };
*/
// 是对接口做的适配，判断响应的报文是否是错误的，如果是错误的进入errorHandler处理
// 当后端接口不满足该规范的时候你需要通过该配置把后端接口数据转换为该格式，该配置只是用于错误处理，不会影响最终传递给页面的数据格式。

export const errorConfig: any = {

    //adaptor是一个函数

    adaptor: (httpResponse: any) => {
        console.log("进入errorConfig.adaptor处理")
        console.log("errorConfig:", httpResponse);
        if (!httpResponse.code) {
            return {
                success: true,
                ...httpResponse,
                errorCode: 100,
                errorMessage: "haha",

            }

        } else {
            return {
                //...res,
                success: httpResponse.code === 0, // success是false时
                data: httpResponse,
                errorCode: httpResponse.errorCode,
                errorMessage: httpResponse.errorMessage, // 显示的错误信息
            };
        }

    },
    // showType: 9, //当 showType 为 9 时，默认会跳转到 /exception 页面，你可以通过该配置来修改该路径。

};

//   主要处理http状态码的问题
export const errorHandler = (error: ResponseError) => {


    console.log("进入errorHandler处理")
    const { response } = error;
    console.log("error:", error);
    console.log("response:", response);

    if (error.name == 'BizError') {

        notification.error({
            message: `errorConfig.adaptor->success==false,msg:` + error.message,
            style: NotificationErrorStyle,
        });
    } else {

        //请求已发送但服务端返回状态码非2xx的响应
        if (response && response.status) {
            const { status, url } = response;
            const errorText = codeMessage[status] || response.statusText;

            notification.error({
                message: `请求错误:状态码:${status}: ${url},${error.message}`,
                description: errorText,
                style: NotificationErrorStyle,
            });
        }

        //服务端没有返回数据时
        if (!response) {
            notification.error({
                description: '您的网络发生异常，无法连接服务器',
                message: '网络异常:' + error.message,
                style: NotificationErrorStyle,
            });
        }

    }

    throw error;
};


