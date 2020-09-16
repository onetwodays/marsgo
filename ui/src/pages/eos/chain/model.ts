import { Effect, ImmerReducer, Reducer, Subscription, request } from 'umi';


export interface chaininfoProps {
    block_cpu_limit: number;
    block_net_limit: number;
    chain_id: string;
    fork_db_head_block_id: string;
    fork_db_head_block_num: number;
    head_block_id: string
    head_block_num: number
    head_block_producer: string
    head_block_time: string
    last_irreversible_block_id: string
    last_irreversible_block_num: number
    server_full_version_string: string
    server_version: string
    server_version_string: string
    virtual_block_cpu_limit: number
    virtual_block_net_limit: number
}


export interface ChainModelState {
    name: string
    chain_info: chaininfoProps



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

const default_chain_info = {
    block_cpu_limit: 0,
    block_net_limit: 0,
    chain_id: "",
    fork_db_head_block_id: "",
    fork_db_head_block_num: 0,
    head_block_id: "",
    head_block_num: 0,
    head_block_producer: "",
    head_block_time: "",
    last_irreversible_block_id: "",
    last_irreversible_block_num: 0,
    server_full_version_string: "",
    server_version: "",
    server_version_string: "",
    virtual_block_cpu_limit: 0,
    virtual_block_net_limit: 0,
};

const default_name = "test";

const ChainModel: ChainModelType = {
    namespace: 'chain',
    state: {
        name: default_name,
        chain_info: default_chain_info,
    },
    effects: {
        * query({ payload }, { call, put }) {
            const data = yield request("eos/v1/chain/get_info");

            console.log(data)
            console.log(1234)

            yield put({
                type: "save",
                payload: {
                    //data,
                    chain_info: data || default_chain_info,
                    name: "接口实时获取",
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
                    dispatch({
                        type: "query",
                    });

                }
            });
        }
    }
};
export default ChainModel;