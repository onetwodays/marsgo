
import styles from './index.less';
import React, { FC } from 'react';
import { ChainModelState, ConnectProps, connect } from 'umi';

interface PageProps extends ConnectProps {
  chain: ChainModelState;
}
const ChainPage: FC<PageProps> = ({ chain, dispatch }) => {
  const { name, chain_info } = chain;
  const {
    server_version, chain_id, head_block_num, last_irreversible_block_num,
    last_irreversible_block_id, head_block_id, head_block_time, head_block_producer,
    virtual_block_cpu_limit, virtual_block_net_limit, block_cpu_limit, block_net_limit,
    server_version_string,
  } = JSON.parse(chain_info);


  return <div >
    <h1>EOS Blockchain Info：</h1>
    <h3><span>server_version：</span>{server_version}</h3>
    <h3><span>chain_id：</span>{chain_id}</h3>
    <h3><span>head_block_num：</span>{head_block_num}</h3>
    <h3><span>last_irreversible_block_num：</span>{last_irreversible_block_num}</h3>
    <h3><span>last_irreversible_block_id：</span>{last_irreversible_block_id}</h3>
    <h3><span>head_block_id：</span>{head_block_id}</h3>
    <h3><span>head_block_time：</span>{head_block_time}</h3>
    <h3><span>head_block_producer：</span>{head_block_producer}</h3>
    <h3><span>virtual_block_cpu_limit：</span>{virtual_block_cpu_limit}</h3>
    <h3><span>virtual_block_net_limit：</span>{virtual_block_net_limit}</h3>
    <h3><span>block_cpu_limit：</span>{block_cpu_limit}</h3>
    <h3><span>block_net_limit：</span>{block_net_limit}</h3>
    <h3><span>server_version_string：</span>{server_version_string}</h3>
  </div>;
};


export default connect(({ chain }: { chain: ChainModelState }) => ({
  chain
}))(ChainPage);