declare namespace API {
  export interface CurrentUser {
    avatar?: string;
    name?: string;
    title?: string;
    group?: string;
    signature?: string;
    tags?: {
      key: string;
      label: string;
    }[];
    userid?: string;
    access?: 'user' | 'guest' | 'admin';
    unreadCount?: number;
  }

  export interface LoginStateType {
    status?: 'ok' | 'error';
    type?: string;
  }

  export interface NoticeIconData {
    id: string;
    key: string;
    avatar: string;
    title: string;
    datetime: string;
    type: string;
    read?: boolean;
    description: string;
    clickClose?: boolean;
    extra: any;
    status: string;
  }




}

/*http://json2ts.com/
在 Pro 推荐在 src/services/API.d.ts 中定义接口数据的类型，以用户基本信息为例
在使用时我们就可以很简单的来使用, d.ts 结尾的文件会被 TypeScript 默认导入到全局，但是其中不能使用 import 语法，如果需要引用需要使用三斜杠。
*/