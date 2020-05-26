import React, { Component } from 'react';
import { Card ,Button } from 'antd';
import { connect } from 'dva';

const namespace = 'puzzlecards';

//mapStateToProps 这个函数的入参 state 其实是 dva 中所有 state 的总合。对于初学 js 的人可能会很疑惑：
//这个入参是谁给传入的呢？其实你不用关心，你只需知道 dva 框架会适时调用 mapStateToProps，并传入 dva model state 作为入参，
//我们再次提醒：传入的 state 是个 "完全体"，包含了 所有 namespace 下的 state！
//我们自己定义的 dva model state 就是以 namespace 为 key 的 state 成员。
//所以 const namespace = 'puzzlecards' 中的 puzzlecards 必须和 model 中的定义完全一致。dva 期待 mapStateToProps 函数返回一个 对象，
//这个对象会被 dva 并入到 props 中
const mapStateToProps = (state) => {
  const cardList = state[namespace].data;
  return {
    cardList,
  };
};

//以 dispatch 为入参，返回一个挂着函数的对象，这个对象上的函数会被 dva 并入 props，注入给组件使用
//dispatch 函数就是和 dva model 打交道的唯一途径。 dispatch 函数接受一个 对象 作为入参，在概念上我们称它为 action，
//唯一强制要包含的是 type 字段，string 类型，用来告诉 dva 我们想要干什么。我们可以选择给 action 附着其他字段，
//这里约定用 payload字段表示额外信息。对组件发消息

const mapDispatchToProps = (dispatch) => {
  return {
    onClickAdd: (newCard) => {
      const action = {
        type: `${namespace}/addNewCard`,
        payload: newCard,
      };
      dispatch(action);
    },

    onDidMount: () => {
      dispatch({
        type: `${namespace}/queryInitCards`,
      });
    },

  };
};

//connect 是连接 dva 和 React 两个平行世界的关键，一定要理解。

//connect 让组件获取到两样东西：1. model 中的数据；2. 驱动 model 改变的方法。
// 本质上只是一个 javascript 函数，通过 @ 装饰器语法使用，放置在组件定义的上方；
//connect 既然是函数，就可以接受入参，第一个入参是最常用的，它需要是一个函数，我们习惯给它命名叫做 mapStateToProps，顾名思义就是把 dva model 中的 state 通过组件的 props 注入给组件。通过实现这个函数，我们就能实现把 dva model 的 state 注入给组件。

@connect(mapStateToProps, mapDispatchToProps)
export default class PuzzleCardsPage extends Component {
  componentDidMount(){
    this.props.onDidMount();
  }
  render() {
    return (
      <div>
        {
          this.props.cardList.map(card => {
            return (
              <Card key={card.id}>
                <div>Q: {card.setup}</div>
                <div>
                  <strong>A: {card.punchline}</strong>
                </div>
              </Card>
            );
          })
        }
        <div>
          <Button onClick={() => this.props.onClickAdd({
            setup: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit',
            punchline: 'here we use dva',
          })}> 添加卡片 </Button>
        </div>
      </div>
    );
  }
}