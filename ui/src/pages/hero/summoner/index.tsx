import React,{FC} from 'react';
import styles from './index.less';
import {connect,SummonerModelState,ConnectProps} from 'umi';
import {Row,Col} from 'antd';




interface PagePros extends  ConnectProps{
  summoner:SummonerModelState;
}

const Summoner:FC<PagePros>=({summoner})=>{
  const { summoners } = summoner;
  return (
    <div>
      <Row>
        {summoners.reverse().map(item => (
          <Col key={item.summoner_id} span={3} className={styles.heroitem}>
            <img src={`https://game.gtimg.cn/images/yxzj/img201606/summoner/${item.summoner_id}.jpg`} />
            <p>{item.summoner_name}</p>
          </Col>
        ))}
      </Row>
    </div>
  );

};

export default connect(({summoner}:{summoner:SummonerModelState})=>({summoner}))(Summoner);







