import React from 'react';
import { fn_get_block } from '@/utils/eosjs';

import { useRequest } from 'umi';

import { Button, Card } from 'antd'


//formatResult: res => res?.data umi要求后端返回的http.body 里面必须有data字段
const blockinfo: React.FC<{}> = () => {
    const { data, error, loading, run } = useRequest(fn_get_block, {
        manual: true,
        formatResult: (res: any) => res,
    });

    console.log("111");

    if (loading) {
        return <div>loading...</div>;
    }

    if (error) {
        return <div>{error.message}</div>;
    }
    return (
        <div>
            <Button onClick={() => run(1)}>区块1</Button>
            <p>{data}</p>
        </div >
    );


};
export default blockinfo;