import { useRequest, useModel } from 'umi';
import React from 'react';



function getGenegis(): Promise<string>
{

        const { substrateApi } = useModel('@@initialState')

        return new Promise((resolve) => { resolve(substrateApi.genesisHash.toHex()) });

}

const Info = () =>
{


        const { data } = useRequest(getGenegis);
        console.log(data)

        return (
                <div>genesisHash:{data}</div>
        );
}

export default Info;
