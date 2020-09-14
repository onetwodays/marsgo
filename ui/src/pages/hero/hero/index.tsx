import React,{FC} from 'react';
import styles from './index.less';
import {connect,HeroModelState,ConnectProps} from 'umi'; // 1 在这里看到，model文件中导出的类型 都可以通过umi导入

import { Row, Col, Radio, Card } from 'antd';
import { RadioChangeEvent } from 'antd/es/radio/interface';
import FreeHeroItem from '@/component/FreeHeroItem';


const RadioGroup =Radio.Group;





interface PageProps extends ConnectProps {
  hero:HeroModelState;
}

const heroType = [
  { key: 0, value: '全部' },
  { key: 1, value: '战士' },
  { key: 2, value: '法师' },
  { key: 3, value: '坦克' },
  { key: 4, value: '刺客' },
  { key: 5, value: '射手' },
  { key: 6, value: '辅助' },
];

const Hero:FC<PageProps>= (props) =>{//2
  console.log(props.hero);
  
  //
  const { heros = [], filterKey = 0, freeheros = [], itemHover = 0 } = props.hero;
  const { dispatch } = props;

  console.log(freeheros);


  const onChange = (e:RadioChangeEvent)=>{
    console.log(e.target.value);
    dispatch!({
      type:"hero/save",
      payload:{
        filterKey:e.target.value
      }
    });
  };


  const onItemHover = (index:number)=>{
    dispatch!({
      type:"hero/save",
      payload:{itemHover:index},
    });

  };
  
  return (
    <div className={styles.normal}>
      <div className={styles.info}>
        <Row className={styles.freehero}>
          <Col span={24}>
            <p>周免英雄</p>
            <div>
              {
                freeheros.map((data, index) => (
                  <FreeHeroItem
                    data={data}
                    itemHover={itemHover}
                    onItemHover={onItemHover}
                    thisIndex={index}
                    key={index}
                  />
                ))
              }
            </div>
            
          </Col>
        </Row>
      </div>
      <Card className={styles.radioPanel}>
        <RadioGroup onChange={onChange} value={filterKey}>
          {heroType.map(data => (
              <Radio value={data.key} key={`hero-rodio-${data.key}`}>
                {data.value}
              </Radio>
            ))}
        </RadioGroup>
        
      </Card>
      <Row>
        {heros.filter(item => filterKey === 0 || item.hero_type === filterKey).reverse().map(item => (
          <Col key={item.ename} span={3} className={styles.heroitem}>
            <img src={`https://game.gtimg.cn/images/yxzj/img201606/heroimg/${item.ename}/${item.ename}.jpg`} />
            <p>{item.cname}</p>
          </Col>
        ))}
      </Row>
    </div>
  );

}
//connect是用来连接前端的ui界面和和前端model的一个嫁接桥梁 ，通过使用connect将model里面定义的state，和dispatch，和histoey方法等传递到前端供前端使用
export default connect(({hero}:{hero:HeroModelState})=>({hero}))(Hero);//3





