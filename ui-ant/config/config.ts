// https://umijs.org/config/
// 构建时配置
import { defineConfig } from 'umi'; // 写配置时也有提示，可以通过 umi 的 defineConfig 方法定义配置
import defaultSettings from './layout';
import proxy from './proxy';
import routes from "./routes";

const { REACT_APP_ENV } = process.env;

// 这是一个函数,聚合umi的插件的配置和其他配置，
export default defineConfig({
    hash: true,
    antd: {},
    dva: {
        hmr: true,
    },
    //构建时layout的配置
    layout: defaultSettings,
    locale: {

        // default zh-CN
        default: 'zh-CN',
        antd: true,
        // default true, when it is true, will use `navigator.language` overwrite default
        baseNavigator: true,

    },
    dynamicImport: {
        loading: '@/components/PageLoading/index',
    },
    targets: {
        ie: 11,
    },
    // umi routes: https://umijs.org/docs/routing
    routes: routes,
    // Theme for antd: https://ant.design/docs/react/customize-theme-cn
    theme: {
        // ...darkTheme,
        'primary-color': defaultSettings.primaryColor,
    },
    // @ts-ignore
    title: false,
    ignoreMomentLocale: true,
    proxy: proxy[REACT_APP_ENV || 'dev'], //代理
    manifest: {
        basePath: '/',
    },
    request: {
        // dataField 对应接口统一格式中的数据字段，比如接口如果统一的规范是 { success: boolean, data: any} ，那么就不需要配置，这样你通过 useRequest 消费的时候会生成一个默认的 formatResult，直接返回 data 中的数据，方便使用。如果你的后端接口不符合这个规范，可以自行配置 dataField 。配置为 '' （空字符串）的时候不做处理。
        dataField: '',
    }



});
