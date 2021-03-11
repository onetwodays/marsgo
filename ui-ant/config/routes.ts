// 相对路径，会从 src/pages 开始找起 ./user/login
// 如果指向 src 目录的文件，可以用 @,比如 component: '@/layouts/basic'


export default  [
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
      redirect: '@/pages/welcome',
    },
    {
      component: './404',
    },

  ]