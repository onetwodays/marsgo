import React, { FC } from 'react';
import { Loading, connect, BlockModelState, ConnectProps } from 'umi'; // 1 在这里看到，model文件中导出的类型 都可以通过umi导入
import { Card, Button } from 'antd';


interface PageProps extends ConnectProps {
    block: BlockModelState;
    loading: boolean;
}



const Block: FC<PageProps> = ({ block, dispatch }) => {
    console.log("block", block);

    const {
        timestamp,
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
    } = block.block_info;

    const onButtonClick = (event) => {

        dispatch!({
            type: "block/query",
            payload: { blockno: 1 },
        });

    };

    return (
        <div>
            <Button type="primary" onClick={onButtonClick}>请求块</Button>

            <Card title={`block ${block.block_no} info`}>
                <p>timestamp:{timestamp}</p>
                <p>producer:{producer}</p>
                <p>confirmed:{confirmed}</p>
                <p>previous:{previous}</p>
                <p>transaction_mroot:{transaction_mroot}</p>
                <p>action_mroot:{action_mroot}</p>
                <p>schedule_version:{schedule_version}</p>
                <p>producer_signature:{producer_signature}</p>
                <p>id:{id}</p>
                <p>block_num{block_num}</p>
                <p>ref_block_prefix:{ref_block_prefix}</p>
            </Card>
        </div>
    )

}


export default connect(
    ({ block, loading }: { block: BlockModelState; loading: Loading }) => ({ block, loading: loading.models.block })
)(Block);
