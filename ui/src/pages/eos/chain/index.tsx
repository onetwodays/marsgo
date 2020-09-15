
import styles from './index.less';
import React, { FC } from 'react';
import { ChainModelState, ConnectProps,  connect } from 'umi';

interface PageProps extends ConnectProps {
  chain: ChainModelState;
}
const ChainPage: FC<PageProps> = ({ chain, dispatch }) => {
  const { name } = chain;
  return <div >Hello {name}</div>;
};


export default connect(({ chain }: { chain: ChainModelState }) => ({
  chain
}))(ChainPage);