
// https://github.com/umijs/umi-request/blob/master/README_zh-CN.md#%E4%B8%AD%E9%97%B4%E4%BB%B6
import { Context } from "umi-request";

// Context 上下文对象，包括 req 和 res 对象
export const demo1Middleware = async (ctx: Context, next: Function) => {
    console.log('request1');
    const { req } = ctx;
    const { url, options } = req;

    // 判断是否需要添加前缀，如果是统一添加可通过 prefix、suffix 参数配置
    if (url.indexOf('/api') !== 0) {
        ctx.req.url = `/api/v1/${url}`;
    }
    ctx.req.options = {
        ...options,
        foo: 'foo',
    };

    await next();

    const { res } = ctx;
    const { success = false } = res; // 假设返回结果为 : { success: false, errorCode: 'B001' }
    if (!success) {
        // 对异常情况做对应处理
    }

    console.log('response1');
};

export const demo2Middleware = async (ctx: Context, next: Function) => {
    console.log('request2');
    await next();
    console.log('response2');
};
