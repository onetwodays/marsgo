import { useRequest } from 'umi';
import React, { useState } from 'react';
import { fn_get_code } from '@/services/eos/get';
import { msgError } from '@/utils/notify';
import { Card, Space } from 'antd'




export default () => {
    const [contract, setContract] = useState<string>('otcexchange');

    const { loading, data } = useRequest(() => fn_get_code(contract), {
        //manual: true,
        refreshDeps: [contract],

        onSuccess: (result, params) => {
            console.log(result);
            //setCode(JSON.stringify(result));

        },
        onError: (error: Error, params: any[]) => {
            console.log(error);

            msgError(error.message);
        }
    });

    return (
        <div>

            <select
                onChange={(e) => setContract(e.target.value)}
                value={contract}
                style={{ marginBottom: 16, width: 120 }}
            >
                <option value="otcexchange">otcexchange</option>
                <option value="otcsystem">otcsystem</option>
                <option value="otc.token">otc.token</option>
            </select>
            <p>
                code:{loading ? 'loading' : JSON.stringify(data)}
            </p>






        </div>
    );
};