import React from 'react';
import { fn_get_block } from '@/services/eos/get';

import { useRequest } from 'umi';

import { useState } from 'react';
import { Button, Card } from 'antd';
import { GetBlockResult } from 'eosjs/dist/eosjs-rpc-interfaces';


//formatResult: res => res?.data umi要求后端返回的http.body 里面必须有data字段
const blockinfo: React.FC<{}> = () => {
    const [state, setState] = useState<GetBlockResult>();
    const { loading, run } = useRequest(fn_get_block, {
        manual: true,
        formatResult: (res: any) => res,
        onSuccess: (res, params) => {
            setState(res);
        },

    });


    if (loading) {
        return <div>loading...</div>;
    }

    const { timestamp,
        producer,
        confirmed,
        previous,
        transaction_mroot,
        action_mroot,
        schedule_version,
        producer_signature,
        id,
        block_num,
        ref_block_prefix
    } = state;









    return (
        <div>
            <Button onClick={() => run(1)}>区块1</Button>

            <Card>{id}</Card>
            <Card>{schedule_version}</Card>
            <Card>{block_num}</Card>
            <Card>{ref_block_prefix}</Card>
            <Card>{timestamp}</Card>
            <Card>{producer}</Card>
            <Card>{producer_signature}</Card>
            <Card>{previous}</Card>
            <Card>{confirmed}</Card>
            <Card>{timestamp}</Card>
            <Card>{transaction_mroot}</Card>
            <Card>{action_mroot}</Card>


        </div >

    );



};
export default blockinfo;