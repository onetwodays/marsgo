import { rpc } from './config';


export const fn_get_block = async (block_num: number) => {
    let res = await rpc.get_block(block_num);
    return {
        "timestamp": res.timestamp,
        "producer": res.producer,
        "confirmed": res.confirmed,
        "previous": res.previous,
        "transaction_mroot": res.transaction_mroot,
        "action_mroot": res.action_mroot,
        "schedule_version": res.schedule_version,
        "producer_signature": res.producer_signature,
        "id": res.id,
        "block_num": res.block_num,
        "ref_block_prefix": res.ref_block_prefix,

    }
};



