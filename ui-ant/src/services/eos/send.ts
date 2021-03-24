import { api, contract } from './config'

export const pushAction = async (actor: string, permission: string, action: string, data: any) => {
    const result = await api.transact({
        actions: [{
            account: contract,
            name: action,
            authorization: [{
                actor,
                permission,
            }],
            data,
        }]
    }, {
        blocksBehind: 3,
        expireSeconds: 30,
    });
    return result;
}
