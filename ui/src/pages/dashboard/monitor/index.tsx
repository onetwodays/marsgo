import React,{useState,useEffect} from 'react';
import {Button} from 'antd'
import moment from 'moment'
import styles from './index.less';

export default () => {
  const [nowTime,setNowTime] = useState<String>(moment().format('YYYY年MM月DD日 ddd HH:mm'));
  var [count, setCount] = useState<Number>(0);

  useEffect(()=>{
    const timer = setInterval(()=>{setNowTime(moment().format('YYYY年MM月DD日 ddd HH:mm:ss'))},1000);
    document.title = `雪${nowTime}`;
    return()=>{clearInterval(timer);}
  },[nowTime]);
  return (
    <div>
      <h1 className={styles.title}>Page dashboard/monitor/index</h1>
      <h1 className={styles.title}>{nowTime}</h1>
      <div>{count}</div>
      <Button onClick={() => { setCount(count + 1); }}>
        点击
      </Button>
      
    </div>
  );
}
