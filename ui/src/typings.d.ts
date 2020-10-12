/*
前端开发很大程序上就是与 Window 打交道，有时候我们不得不给 Window 增加参数，例如各种统计的代码。在 TypeScript 中提供一个方式来增加参数
在 /src/typings.d.ts 中做如下定义
*/
declare module 'slash2';
declare module '*.css';
declare module '*.less';
declare module '*.scss';
declare module '*.sass';
declare module '*.svg';
declare module '*.png';
declare module '*.jpg';
declare module '*.jpeg';
declare module '*.gif';
declare module '*.bmp';
declare module '*.tiff';
declare module 'omit.js';

// google analytics interface
interface GAFieldsObject {
    eventCategory: string;
    eventAction: string;
    eventLabel?: string;
    eventValue?: number;
    nonInteraction?: boolean;
}
interface Window {
    ga: (
        command: 'send',
        hitType: 'event' | 'pageview',
        fieldsObject: GAFieldsObject | string,
    ) => void;
    reloadAuthorized: () => void;
}
//如果不想在 Window 中增加，但是想要全局使用，比如通过 define 注入的参数，我们通过 declare 关键字在 /src/typings.d.ts 注入
declare let ga: Function;

// preview.pro.ant.design only do not use in your production ;
// preview.pro.ant.design 专用环境变量，请不要在你的项目中使用它。
declare let ANT_DESIGN_PRO_ONLY_DO_NOT_USE_IN_YOUR_PRODUCTION: 'site' | undefined;

declare const REACT_APP_ENV: 'test' | 'dev' | 'pre' | false;
