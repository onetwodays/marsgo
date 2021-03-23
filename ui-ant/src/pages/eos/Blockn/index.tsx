import React from 'react';
import { useRequest } from 'umi';
import { fn_get_block } from '@/services/eos/get'


//异常没有处理
const Blockn = () => {

    const { data, loading, run } = useRequest(fn_get_block, { manual: true, });
    console.log(data);
    console.log(Object.prototype.toString.call(data));

    return (
        <div>
            <p>Enter Block number</p>
            <input placeholder={"1"} onChange={(e) => run(e.target.value)}></input>
            {loading ? (<p>loading</p>) : (<p>{JSON.stringify(data, null, 4)}</p>)}

        </div>

    );
};

export default Blockn;
