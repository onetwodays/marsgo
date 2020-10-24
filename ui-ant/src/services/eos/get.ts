import { rpc } from './config';


export const fn_get_block = async (block_num: number) => {
    let res = await rpc.get_block(block_num);
    return res;
};



