import { Effect, Reducer } from 'umi';
import { fn_get_block } from '@/services/eos/get'



export interface BlockProps {
    timestamp: string;
    producer: string;
    confirmed: number;
    previous: string;
    transaction_mroot: string;
    action_mroot: string;
    schedule_version: number;
    producer_signature: string;
    id: string;
    block_num: number;
    ref_block_prefix: number;
}

export interface BlockModelState {
    block_no: number;
    block_info: BlockProps;
}

export interface BlockModelType {
    namespace: string;
    state: BlockModelState;
    effects: {
        query: Effect,

    };
    reducers: {
        save: Reducer<BlockModelState>,
    };
}

const default_block: BlockProps = {
    timestamp: '',
    producer: '',
    confirmed: 0,
    previous: '',
    transaction_mroot: '',
    action_mroot: '',
    schedule_version: 0,
    producer_signature: '',
    id: "",
    block_num: 0,
    ref_block_prefix: 0
};


const BlockModel: BlockModelType = {
    namespace: 'block',
    state: {
        block_no: 0,
        block_info: default_block,

    },
    reducers: {
        save(state, action) {
            return {
                ...state,
                ...action.payload,
            };
        },

    },
    effects: {
        *query({ payload: { blockno } }, { call, put, select }) {
            let res = yield fn_get_block(blockno);
            console.log("res:", res);
            yield put({
                type: "save",
                payload: {
                    block_no: blockno,
                    block_info: res,
                }
            });


        },


    },


};

export default BlockModel;


