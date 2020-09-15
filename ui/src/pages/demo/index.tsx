import React,{ FC,useState,useEffect } from 'react';

import {ConnectProps,DemoModelState,connect } from 'umi'
import {Button,Card,Row,Col} from 'antd'
import styles from './index.less';
import moment from 'moment'


const namespace = 'demo';

interface PageProps extends  ConnectProps{
    demo: DemoModelState;
}


const Demo:FC<PageProps> = ({demo,dispatch}) => {


  const { rows=[]} =demo;
  
  const onButtonPushDown= ()=>{
    dispatch!({
      type:`${namespace}/fetch`,
      payload:{
        rows:[{id:2,name:'zhouhao',desc:'周浩',url:'www.baidu.com'},{id:1,name:'caoxue',desc:'曹雪',url:''}],
        filterKey:"zhouhao"
      }

    });

  };

  const [nowTime,setNowTime] = useState<String>(moment().format('YYYY年MM月DD日 ddd HH:mm'));
  let   [count,  setCount]   = useState<Number>(0);

  
  useEffect(()=>{
    const timer = setInterval(()=>{setNowTime(moment().format('YYYY年MM月DD日 ddd HH:mm:ss'))},1000);
    document.title = `曹雪${nowTime}`;
    return ()=>{clearInterval(timer);}
  },[nowTime]);


  return(


    
      <div>
        <h1 className={styles.title}>{nowTime}</h1>
        <div>{count}</div>
        <Button onClick={() => { setCount(count + 1); }}>
          点击
        </Button>


        <Button onClick={onButtonPushDown}>
          添加
        </Button>

       
        <Card title={'表信息'}>
          <Row>
          {rows.reverse().map(item => (
          <Col key={item.id} span={3} >
            <p>{item.name}</p>
          </Col>
        ))}
          </Row>
        </Card>
       
      </div>
  );
}

export default connect(
  ({demo}:{demo:DemoModelState})=>({demo}) 
  )(Demo); //连接页面和models
