/*
运行时配置文件，可以在这里扩展运行时的能力，比如修改路由、修改 render 方法等。
运行时配置和配置的区别是他跑在浏览器端，基于此，我们可以在这里写函数、jsx、import 浏览器端依赖等等，注意不要引入 node 依赖

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

import React from 'react';
import { BasicLayoutProps, Settings as LayoutSettings } from '@ant-design/pro-layout';
import { notification } from 'antd';
import { history, RequestConfig } from 'umi';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import { ResponseError } from 'umi-request';
import { queryCurrent } from './services/user';
import defaultSettings from '../config/defaultSettings';

export interface InitialStateType {
  settings?: LayoutSettings;
  currentUser?: API.CurrentUser;
  fetchUserInfo: () => Promise<API.CurrentUser | undefined>;

};
/* 
* 该方法返回的数据最后会被默认注入到一个 namespace 为 @@initialState  的 model 中。可以通过 useModel
* const { initialState, loading, refresh, setInitialState } = useModel('@@initialState');

*/
export async function getInitialState(): Promise<InitialStateType> {
  // 先定义一个函数
  const fetchUserInfo = async () => {
    try {
      const currentUser = await queryCurrent();
      return currentUser;
    } catch (error) {
      console.log("getInitialState异常");
      history.push('/user/login');
      console.log("getInitialState异常");
    }
    return undefined;
  };
  // 如果是登录页面，不执行
  if (history.location.pathname !== '/user/login') {
    const currentUser = await fetchUserInfo();
    /*
    {
      data: {isLogin: false}
      isLogin: false
      errorCode: "401"
      errorMessage: "嘿嘿请先登录！"
      success: true
    }
    */
    return {
      fetchUserInfo,
      currentUser,
      settings: defaultSettings,
    };
  }
  return {
    fetchUserInfo,
    settings: defaultSettings,
  };
}
// 运行时配置
export const layout = ({
  initialState,
}: {
  initialState: { settings?: LayoutSettings; currentUser?: API.CurrentUser };
}): BasicLayoutProps => {
  return {
    rightContentRender: () => <RightContent />,
    disableContentMargin: false,
    footerRender: () => <Footer />,
    onPageChange: () => {
      console.log("layout11111");
      const { currentUser } = initialState;
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!currentUser && location.pathname !== '/user/login') {
        history.push('/user/login');
      }
      console.log("layout 22222");
    },
    menuHeaderRender: undefined,
    ...initialState?.settings,

  };
};

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

/**
 * 异常处理程序
 * {
        data: {
          isLogin: false,
        },
        errorCode: '401',
        errorMessage: '嘿嘿请先登录！',
        success: true,
   }
 */
const errorHandler = (error: ResponseError) => {
  // 是否是 BizError 来判断是否是因为 success 为 false 抛出的错误
  if (error.name === "BizError") {
    notification.error({
      message: `请求错误 ${error.data}`,
      // description: error.data.msg,
    });
    return error.data.code;
  }
  const { response, data } = error;
  if (response && response.status) {
    const errorText = codeMessage[response.status] || response.statusText; //自己生成的
    const { status, url } = response;

    //请求错误 401: http://localhost:8000/api/currentUser,[object Object]
    //用户没有权限（令牌、用户名、密码错误）


    notification.error({
      message: `请求错误 ${status}: ${url},${data}`,
      description: errorText,
      style: NotificationErrorStyle,

    });

  }

  if (!response) {
    notification.error({
      description: '您的网络发生异常，无法连接服务器',
      message: '网络异常',
      style: NotificationErrorStyle,
    });

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
// 相当于做了一个格式转换
const errorConfig: any = {
  adaptor: (res: any) => {
    return {
      //...res,
      success: res.code === 200, // success是false时
      data: res,
      errorCode: res.errorCode,
      errorMessage: res.errorMessage, // 显示的错误信息
    };
  },
  // showType: 9, //当 showType 为 9 时，默认会跳转到 /exception 页面，你可以通过该配置来修改该路径。

};

// a1 -> b1 -> response -> b2 -> a2 
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
  //prefix: '/api', // 所有的请求的前缀,相当于ip+port部分
  errorHandler: errorHandler,
  credentials: 'include', // 默认请求是否带上cookie
  timeout: 5000,

  //errorConfig: errorConfig, // 自定义接口规范
  middlewares: [middlewareA, middlewareB],

  requestInterceptors: [],
  responseInterceptors: [],



};





/*
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
*/