import React,{ FC } from 'react';
import styles from './index.less';

import {ConnectProps,BaseModelStateType ,connect   } from 'umi'
import {Button} from 'antd'



interface PageProps extends  ConnectProps{
  base : BaseModelStateType;
}

const Base:FC<PageProps> = ({base,dispatch}) => {
  const onButtonPushDown= ()=>{
    dispatch!({
      type:'base/fetch',
      payload:{
        name:'iliu'
      }

    });

  };

  return(
    
      <div>
        <h2> this is {base.name}</h2>
        <Button onClick={onButtonPushDown}>获取信息</Button>
      </div>
  );
}

export default connect(
  ({base}:{base:BaseModelStateType})=>({base}) 
  )(Base); //连接页面和models

