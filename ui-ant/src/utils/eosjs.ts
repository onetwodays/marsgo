import { Api, JsonRpc } from 'eosjs';
import { JsSignatureProvider } from 'eosjs/dist/eosjs-jssig';  // development only

const privateKeys = ["5JBdZDCH4NKsGXR4DJkiaSQQ4kBE9vnePA1Be2RLeaSFvDqBzLg"];

const signatureProvider = new JsSignatureProvider(privateKeys);
const rpc = new JsonRpc('http://127.0.0.1:18888'); //required to read blockchain state
const api = new Api({ rpc, signatureProvider }); //required to submit transactions

export async function fn_get_block(block_num: number) {
    return await rpc.get_block(block_num);
}

