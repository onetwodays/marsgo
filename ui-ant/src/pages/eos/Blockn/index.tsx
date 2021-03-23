import React from 'react';
import { useRequest } from 'umi';
import { fn_get_block } from '@/services/eos/get'
import { GetBlockResult } from 'eosjs/dist/eosjs-rpc-interfaces';

export default {
    const { data } = useRequest<GetBlockResult>(() => { return fn_get_block(1); });
    console.log(data);

    return(
    <p> 111</p >

     );




};
