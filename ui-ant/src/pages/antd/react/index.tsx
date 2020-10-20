import React from 'react';
import { useModel } from 'umi';
import { useAccess, Access } from 'umi'; //页面内的权限控制 useAccess是hook来获取权限定义，Access组件用于页面的元素显示和隐藏的控制
import { Button, notification } from 'antd';

const openNotification = () => {
    notification.open({
        message: 'Notification Title',
        description:
            'This is the content of the notification. This is the content of the notification. This is the content of the notification.',
        onClick: () => {
            console.log('Notification Clicked!');
        },
    });
};
export default () => {
    const { counter, add, minus } = useModel('counter', (ret) => ({
        add: ret.increment,
        minus: ret.decrement,
        counter: ret.counter
    }));

    const access = useAccess();//access 实例的成员: canReadFoo, canUpdateFoo, canDeleteFoo
    if (access.canReadFoo) {
        console.log('页面内权限控制!');
    }
    const foo = 3;

    return (<div>
        <p>{counter}</p>
        <Button onClick={add}>add by 1</Button>
        <Button onClick={minus}>minus by 1</Button>
        <Access accessible={access.canReadFoo} fallback={<div>Can not read foo content.</div>}>
            Foo content.
      </Access>
        <Access accessible={access.canUpdateFoo()} fallback={<div>Can not update foo.</div>}>
            Update foo.
      </Access>
        <Access accessible={access.canDeleteFoo(foo)} fallback={<div>Can not delete foo.</div>}>
            Delete foo.
      </Access>
    </div>);;
}

/*
useModel 可以接受一个可选的第二个参数，可以用于性能优化。
当组件只需要消费 model 中的部分参数，而对其他参数的变化并不关心时，
可以传入一个函数用于过滤。函数的返回值将取代 model 的返回值，成为 useModel 的最终返回值
*/