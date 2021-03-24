
import React from 'react';
import { Settings as LayoutSettings, PageLoading } from '@ant-design/pro-layout';

import { history, RequestConfig, RunTimeLayoutConfig } from 'umi';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';

import { queryCurrent } from './services/user';
import MyLayout from '../config/layout';
import { msgError } from '@/utils/notify'
import { errorConfig, errorHandler } from "@/net/request";
import { eosResponseInterceptor } from '@/net/resInterceptor';


// const isDev = process.env.NODE_ENV === 'development';

/** 获取用户信息比较慢的时候会展示一个 loading */
export const initialStateConfig = {
    loading: <PageLoading />,
};


// 定义一个类型
export interface InitialStateType {
    settings?: LayoutSettings;  //默认的布局的属性
    currentUser?: API.CurrentUser;//用户的属性
    fetchUserInfo?: () => Promise<API.CurrentUser | undefined>; //获取用户信息的request

};
/*
* 该方法返回的数据最后会被默认注入到一个 namespace 为 @@initialState  的 model 中。可以通过 useModel
* const { initialState, loading, refresh, setInitialState } = useModel('@@initialState');
* 该配置是一个 async 的 function。会在整个应用最开始执行，返回值会作为全局共享的数据。Layout 插件、Access 插件以及用户都可以通过 useModel('@@initialState') 直接获取到这份数据。
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
            settings: MyLayout, //这个是布局的配置

        };
    }
    return {
        fetchUserInfo,
        settings: MyLayout,
    };
}


// 运行时配置
export const layout: RunTimeLayoutConfig = ({ initialState }) => {
    return {
        rightContentRender: () => <RightContent />, //右上角 UI 完全的自定义
        disableContentMargin: false,
        footerRender: () => <Footer />, //自定义页脚
        onPageChange: () => {
            const { location } = history;
            // 如果没有登录，重定向到 login
            if (!initialState?.currentUser && location.pathname !== '/user/login') {
                history.push('/user/login');
            }

        },
        menuHeaderRender: undefined,
        // 自定义 403 页面
        // unAccessible: <div>unAccessible</div>,
        ...initialState?.settings,

    };
};



// 网络请求的运行时配置（在浏览器里运行）
export const request: RequestConfig = {
    //prefix: '/api', // 所有的请求的前缀,相当于ip+port部分
    //errorConfig: errorConfig,
    errorHandler: errorHandler,
    credentials: 'include', // 默认请求是否带上cookie
    timeout: 5000,
    middlewares: [],
    requestInterceptors: [],
    responseInterceptors: [eosResponseInterceptor],
};


//如果 model里面的effects 中抛异常没有被捕获，会执行 onError，然后才是组件的 dispatch 返回的 Promise 处理。
//如果在 onError 中调用 err. preventDefault() 则后续 dispatch 的 catch 不会执行
export const dva = {
    config: {
        onError(e: Error) {
            console.log("dva全局错误处理:", e.message)
            msgError(e.message);
            e.preventDefault();
        },
    },
};
