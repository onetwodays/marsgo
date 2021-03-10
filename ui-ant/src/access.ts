// src/access.ts
export default function access(initialState: { currentUser?: API.CurrentUser | undefined }) {
  const { currentUser } = initialState || {};
  return {
    canAdmin: currentUser && currentUser.access === 'admin', //false,该条路由将会被禁用,该条路由将会被禁用，并且从左侧 layout 菜单中移除，如果直接从 URL 访问对应路由，将看到一个 403 页面
    canReadFoo: true,
    canUpdateFoo: () => false,
    canDeleteFoo: (data: any) => data?.status < 1, //按业务需求自己任意定义鉴权函数
  };
}

/*
权限的定义依赖于初始数据，初始数据需要通过 @umijs/plugin-initial-state 生成。

生成完初始化数据后，就可以开始定义权限了。首先新建 src/access.ts ，在该文件中 export default 一个函数，定义用户拥有的权限

该文件需要返回一个 function，返回的 function 会在应用初始化阶段被执行，
执行后返回的对象将会被作为用户所有权限的定义。
对象的每个 key 对应一个 boolean 值，只有 true 和 false，代表用户是否有该权限
*/
