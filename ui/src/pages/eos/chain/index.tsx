
import styles from './index.less';
import React, { FC } from 'react';
import { ChainModelState, ConnectProps, connect } from 'umi';

interface PageProps extends ConnectProps {
  chain: ChainModelState;
}
const ChainPage: FC<PageProps> = ({ chain, dispatch }) => {


  console.log(chain);

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
  } = chain.chain_info;


  return (
    <div >
      <h1>EOS Blockchain Info：{chain.name}</h1>
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
      <h3><span>block_net_limit：</span>{block_net_limit}</h3>




    </div>);
};


export default connect(({ chain }: { chain: ChainModelState }) => ({
  chain
}))(ChainPage);