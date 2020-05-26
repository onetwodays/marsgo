import React,{ FC } from 'react';

import {ConnectProps,DemoModelStateType,connect } from 'umi'
import {Button,Card} from 'antd'
import styles from './index.less';


const namespace = 'demo';

interface PageProps extends  ConnectProps{
    demo: DemoModelStateType;
}


const Demo:FC<PageProps> = ({demo,dispatch}) => {
  const onButtonPushDown= ()=>{
    dispatch!({
      type:'demo/fetch',
      payload:{
        rows:[{id:2,name:'zhouhao',desc:'周浩',url:'www.baidu.com'}],
        filterKey:"zhouhao"
      }

    });

  };

  return(
    
      <div>
        {
          demo.rows.map(()=>{<Card>1111</Card>})
        }
        <Button onClick={onButtonPushDown}>获取信息</Button>
      </div>
  );
}

export default connect(
  ({demo}:{demo:Demo1ModelStateType})=>({demo}) 
  )(Demo); //连接页面和models
