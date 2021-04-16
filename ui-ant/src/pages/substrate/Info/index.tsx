import { useRequest, useModel } from 'umi';
import React from 'react';
import { ApiPromise, WsProvider } from '@polkadot/api';
import Inspector from 'react-json-inspector'

import JSONTree from 'react-json-tree'
const substrateUrl = 'ws://localhost:9944';
const ADDR = '5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY'
const theme = {
    scheme: 'monokai',
    author: 'wimer hazenberg (http://www.monokai.nl)',
    base00: '#272822',
    base01: '#383830',
    base02: '#49483e',
    base03: '#75715e',
    base04: '#a59f85',
    base05: '#f8f8f2',
    base06: '#f5f4f1',
    base07: '#f9f8f5',
    base08: '#f92672',
    base09: '#fd971f',
    base0A: '#f4bf75',
    base0B: '#a6e22e',
    base0C: '#a1efe4',
    base0D: '#66d9ef',
    base0E: '#ae81ff',
    base0F: '#cc6633'
};

async function getUsername(): Promise<string>
{
    const api: ApiPromise = await ApiPromise.create({ provider: new WsProvider(substrateUrl) });

    return new Promise((resolve) =>
    {
        const hash: string = api.genesisHash.toHex();

        console.log("hash=", hash)
        resolve(hash);
    });
}


const Info = () =>
{
    const { initialState } = useModel('@@initialState')
    const genesisHash = initialState?.substrateApi?.genesisHash.toHex();
    const libraryInfo = initialState?.substrateApi?.libraryInfo;
    const runtimeVersion = initialState?.substrateApi?.runtimeVersion.toString();
    const runtimeMetadata1 = initialState?.substrateApi?.runtimeMetadata;
    const runtimeMetadataJson = JSON.parse(runtimeMetadata1?.toString() || '');
    console.log(runtimeMetadataJson)






    return (
        <div>

            <div>genesisHash:{genesisHash}</div>
            <div>libraryInfo:{libraryInfo}</div>
            <div>runtimeVersion:{runtimeVersion}</div>
            <div>e:{JSON.stringify(initialState?.substrateApi?.runtimeMetadata.asLatest.toHuman(), null, 2)}</div>

            <div>
                <p>runtimeMetadataJson</p>
                <JSONTree data={runtimeMetadataJson} theme={{
                    extend: theme,
                    // underline keys for literal values
                    valueLabel: {
                        textDecoration: 'underline'
                    },
                    // switch key for objects to uppercase when object is expanded.
                    // `nestedNodeLabel` receives additional arguments `expanded` and `keyPath`
                    nestedNodeLabel: ({ style }, nodeType, expanded) => ({
                        style: {
                            ...style,
                            textTransform: expanded ? 'uppercase' : style.textTransform
                        }
                    })
                }} />
            </div>





        </div>

    );
}

export default Info;
