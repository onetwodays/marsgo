import { Api, JsonRpc } from 'eosjs';
import { JsSignatureProvider } from 'eosjs/dist/eosjs-jssig';  // development only


//import ScatterJS from 'scatterjs-core'
//import ScatterEOS from 'scatterjs-plugin-eosjs2'

//导入钱包相关的

const appName = 'otcexchange';
const contract = 'otcexchange';

//jungle testnet
/*
const network = {
    blockchain: 'eos',
    protocol: 'https',
    host: 'jungle2.cryptolions.io',
    port: 443,
    chainId: 'e70aaab8997e1dfce58fbfac80cbbb8fecec7b99cf982a9444273cbc64c41473',
};
*/



// mainnet
// const network = {
//   blockchain: 'eos',
//   protocol: 'https',
//   host: 'api.eosnewyork.io',
//   port: 80,
//   chainId: 'aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906',
// };


// local
const network = {
    blockchain: 'eos',
    protocol: 'http',
    host: 'zhongyingying.qicp.io',
    port: 38000,
    chainId: '215331c6c3863ad74a1e1680234f2237374d1d46ab2a3d7534126d01d815a488',
};

//ScatterJS.plugins(new ScatterEOS());

//const signatureProvider = ScatterJS.scatter.eosHook(network, null, true);




const privateKeys = ["5JBdZDCH4NKsGXR4DJkiaSQQ4kBE9vnePA1Be2RLeaSFvDqBzLg"];
const signatureProvider = new JsSignatureProvider(privateKeys);

const url = network.protocol + '://' + network.host + ':' + network.port;

const rpc = new JsonRpc(url, { fetch })
const api = new Api({
    rpc,
    signatureProvider,
    chainId: network.chainId,
    textDecoder: new TextDecoder(),
    textEncoder: new TextEncoder(),
});

export { api, rpc, network, appName, contract }