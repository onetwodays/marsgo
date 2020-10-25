import { GetCodeResult } from 'eosjs/dist/eosjs-rpc-interfaces';
import { rpc } from './config';


export const fn_get_block = async (block_num: number) => {
    let res = await rpc.get_block(block_num);
    return res;
};


export const fn_get_code = async (code: string): Promise<GetCodeResult> => {
    return rpc.get_code(code);
};





