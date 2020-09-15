import { Effect, ImmerReducer, Reducer, Subscription,request } from 'umi';



export interface ChainModelState {
  name: string;
}

export interface ChainModelType {
  namespace: 'chain';
  state: ChainModelState;
  effects: {
    query: Effect;
  };
  reducers: {
    save: Reducer<ChainModelState>;
    // 启用 immer 之后
    // save: ImmerReducer<ChainModelState>;
  };
  subscriptions: { setup: Subscription };
}
const ChainModel: ChainModelType = {
  namespace: 'chain',
  state: {
    name: '',
  },
  effects: {
    *query({ payload }, { call, put }) {
        const data = yield request("eos/v1/chain/get_info");
        yield put({
            type:"save",
            payload:{
                name:data,
            },
        });

    },
  },
  reducers: {
    save(state, action) {
      return {
        ...state,
        ...action.payload,
      };
    },
    // 启用 immer 之后
    // save(state, action) {
    //   state.name = action.payload;
    // },
  },
  subscriptions: {
    setup({ dispatch, history }) {
      return history.listen(({ pathname }) => {
        if (pathname === '/eos/chain') {
        
        }
      });
    }
  }
};
export default ChainModel;