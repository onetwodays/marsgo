// 相对路径，会从 src/pages 开始找起 ./user/login
// 如果指向 src 目录的文件，可以用 @,比如 component: '@/layouts/basic'
// 路由管理 通过约定的语法根据在 config.ts 中配置路由。

export default [
    {
        path: '/user',
        layout: false,
        routes: [
            {
                name: 'login',
                path: '/user/login',
                component: './user/Login',
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
            },
            {
                path: '/eos/code',
                name: 'code',
                icon: 'smile',
                component: './eos/code',
            },
            {
                path: '/eos/adorder',
                name: 'ad',
                icon: 'smile',
                component: './eos/adorder',
            },
        ]
    },
    {
        path: '/antd',
        name: 'antd',
        icon: 'crown',
        //component: './antd',
        routes: [
            {   //菜单跳转到外部地址
                path: 'https://beta-pro.ant.design/docs/router-and-nav-cn',
                target: '_blank', // 点击新窗口打开
                icon: 'smile',
                name: "ant.pro文档",
            },
            {   //菜单跳转到外部地址
                path: 'https://ahooks.js.org/zh-CN/hooks/async/',
                target: '_blank', // 点击新窗口打开
                icon: 'GiftTwoTone',
                name: "useRequest文档",
            },
            {   //菜单跳转到外部地址
                path: 'https://umijs.org/zh-CN/docs',
                target: '_blank', // 点击新窗口打开
                icon: 'MessageTwoTone',
                name: "umi文档",
            },
            {   //菜单跳转到外部地址
                path: 'https://ant.design/components/overview-cn/',
                target: '_blank', // 点击新窗口打开
                icon: 'EyeTwoTone',
                name: "antd文档",
            },
            {   //菜单跳转到外部地址
                path: 'https://procomponents.ant.design/components/',
                target: '_blank', // 点击新窗口打开
                icon: 'EyeTwoTone',
                name: "antdpro文档",
            },



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
            },
            {
                path: '/antd/productlist',
                name: 'productlist',
                icon: 'smile',
                component: './antd/product',
            },


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
        redirect: './welcome',
    },
    {
        component: './404',
    },

]
