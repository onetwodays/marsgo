import React,{ FC } from 'react';

import {ConnectProps,DemoModelStateType,connect } from 'umi'
import {Button,Card,Row,Col} from 'antd'
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
        <p>{demo.rows[0].name}</p>
        {
          demo.rows.map((item)=>{<Card>{item.name}</Card>})
        }
        <Button onClick={onButtonPushDown}>获取信息</Button>
        <Card title={'按钮'}>
        <Button onClick={onButtonPushDown}></Button>
        <Button type="primary">主按钮：用于主行动点，一个操作区域只能有一个主按钮</Button>
        <Button>默认按钮：用于没有主次之分的一组行动点。</Button>
        <Button type="dashed">虚线按钮：常用于添加操作</Button>
        <Button type="link">链接按钮：用于次要或外链的行动点</Button>

        <Row>
      <Col span={24}>col</Col>
    </Row>
    <Row>
      <Col span={12}>col-12</Col>
      <Col span={12}>col-12</Col>
    </Row>
    <Row>
      <Col span={8}>col-8</Col>
      <Col span={8}>col-8</Col>
      <Col span={8}>col-8</Col>
    </Row>
    <Row>
      <Col span={6}>col-6</Col>
      <Col span={6}>col-6</Col>
      <Col span={6}>col-6</Col>
      <Col span={6}>col-6</Col>
    </Row>

          
        </Card>
       
      </div>
  );
}

export default connect(
  ({demo}:{demo:DemoModelStateType})=>({demo}) 
  )(Demo); //连接页面和models
