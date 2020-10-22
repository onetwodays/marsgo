// https://umijs.org/config/
// 构建时配置
import { defineConfig } from 'umi'; // 写配置时也有提示，可以通过 umi 的 defineConfig 方法定义配置
import defaultSettings from './defaultSettings';
import proxy from './proxy';

const { REACT_APP_ENV } = process.env;

export default defineConfig({
  hash: true,
  antd: {},
  dva: {
    hmr: true,
  },
  //构建时layout的配置
  layout: {
    name: 'Ant Design Pro',
    locale: true, //通过 layout 配置的 locale 配置开启国际化，开启后路由里配置的菜单名会被当作菜单名国际化的 key，插件会去 locales 文件中查找 menu.[key]对应的文案，默认值为改 key
    ...defaultSettings,
  },
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
  routes: [
    {
      path: '/user',
      layout: false,
      routes: [
        {
          name: 'login',
          path: '/user/login',
          component: './user/login',
        },
      ],
    },
    {
      path: '/welcome',
      name: 'welcome',
      icon: 'smile',
      component: './Welcome',

    },
    {
      path: '/admin',
      name: 'admin',
      icon: 'crown',
      access: 'canAdmin',
      component: './Admin',
      routes: [
        {
          path: '/admin/sub-page',
          name: 'sub-page',
          icon: 'smile',
          component: './Welcome',
        },
      ],
    },
    {
      path: '/eos',
      name: 'eos',
      icon: 'crown',
      //component: './eos/chain',
      routes: [
        {
          path: '/eos/chain',
          name: 'chaninfo',
          icon: 'smile',
          component: './eos/chain',
        },
        {
          path: '/eos/block',
          name: 'blockinfo',
          icon: 'smile',
          component: './eos/block',
        }
      ]
    },
    {
      path: '/antd',
      name: 'antd',
      icon: 'crown',
      //component: './antd',
      routes: [
        {
          path: '/antd/dataentry',
          name: 'dataentry',
          icon: 'smile',
          component: './antd/dataentry',
        },
        {
          path: '/antd/datadisplay',
          name: 'datadisplay',
          icon: 'smile',
          component: './antd/datadisplay',
        },
        {
          path: '/antd/react',
          name: 'react',
          icon: 'smile',
          component: './antd/react',
        }
      ],
    },
    {
      name: 'list.table-list',
      icon: 'table',
      path: '/list',
      component: './ListTableList',
    },
    {
      path: '/',
      redirect: '/welcome',
    },
    {
      component: './404',
    },

  ],
  // Theme for antd: https://ant.design/docs/react/customize-theme-cn
  theme: {
    // ...darkTheme,
    'primary-color': defaultSettings.primaryColor,
  },
  // @ts-ignore
  title: false,
  ignoreMomentLocale: true,
  proxy: proxy[REACT_APP_ENV || 'dev'],
  manifest: {
    basePath: '/',
  },
  request: {
    dataField: '',
  }



});
