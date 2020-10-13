import React from 'react';
import { useRequest } from 'umi';

//formatResult: res => res?.data umi要求后端返回的http.body 里面必须有data字段
const chaininfo: React.FC<{}> = () => {
    const { data, error, loading } = useRequest('/v1/chain/get_info', {
        //debounceInterval: 500,
        pollingInterval: 2000,
        pollingWhenHidden: false,
        formatResult: (res: any) => res,


    });

    if (loading) {
        return <div>loading...</div>;
    }

    if (error) {
        return <div>{error.message}</div>;
    }

    const {

        block_cpu_limit,
        block_net_limit,
        chain_id,
        fork_db_head_block_id,
        fork_db_head_block_num,
        head_block_id,
        head_block_num,
        head_block_producer,
        head_block_time,
        last_irreversible_block_id,
        last_irreversible_block_num,
        server_full_version_string,
        server_version,
        server_version_string,
        virtual_block_cpu_limit,
        virtual_block_net_limit,
    } = data;
    return <div>
        <h1>EOS Blockchain Info：EOS</h1>
        <h3><span>server_version：</span>{server_version}</h3>
        <h3><span>server_version_string：</span>{server_version_string}</h3>
        <h3><span>server_full_version_string：</span>{server_full_version_string}</h3>
        <h3><span>chain_id：</span>{chain_id}</h3>
        <h3><span>head_block_num：</span>{head_block_num}</h3>
        <h3><span>last_irreversible_block_num：</span>{last_irreversible_block_num}</h3>
        <h3><span>last_irreversible_block_id：</span>{last_irreversible_block_id}</h3>
        <h3><span>head_block_id：</span>{head_block_id}</h3>
        <h3><span>head_block_time：</span>{head_block_time}</h3>
        <h3><span>head_block_producer：</span>{head_block_producer}</h3>
        <h3><span>fork_db_head_block_id：</span>{fork_db_head_block_id}</h3>
        <h3><span>fork_db_head_block_num：</span>{fork_db_head_block_num}</h3>
        <h3><span>virtual_block_cpu_limit：</span>{virtual_block_cpu_limit}</h3>
        <h3><span>virtual_block_net_limit：</span>{virtual_block_net_limit}</h3>
        <h3><span>block_cpu_limit：</span>{block_cpu_limit}</h3>
        <h3><span>block_net_limit：</span>{block_net_limit}</h3></div>;
};
export default chaininfo