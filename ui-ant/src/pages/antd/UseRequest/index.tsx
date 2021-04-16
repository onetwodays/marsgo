import { useRequest } from "umi";
import { message } from "antd";
import React, { useState } from "react";



function changeUserName(username: string): Promise<{ success: boolean }>
{
    console.log("username=", username);
    return new Promise((resolve) =>
    {
        setTimeout(() =>
        {
            resolve({ success: true });

        }, 1000);
    });


}

export default () =>
{
    const [state, setState] = useState<string>('');
    const { loading, run } = useRequest(changeUserName, {
        manual: true,
        onSuccess: (result, params) =>
        {
            if (result.success)
            {
                setState("");
                message.success(`The username was changed to"${params[0]}"`)
            }
        },

    });


    return (
        <div>
            <input onChange={(e) => setState(e.target.value)}
                value={state}
                placeholder="Please enter username"
                style={{ width: 240, marginRight: 16 }}
            />
            <button disabled={loading} type="button" onClick={() => run(state)}>
                {loading ? 'loading' : 'Edit'}

            </button>
        </div>

    );
};
