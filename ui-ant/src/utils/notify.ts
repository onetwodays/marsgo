import { notification, message } from 'antd';

export const msgSuccess = (msg: string) => {
    message.success(msg);
}

export const msgError = (msg: string) => {
    message.error(msg);
}

export const msgTx = (txid: string) => {
    msgSuccess('excute success, tx_id:' + txid);
}

export const notify = (message: string, description: string) => {
    notification.error({
        message,
        description,
    });
};

export const isToday = (time: number) => {
    return new Date(time).toDateString() === new Date(new Date().valueOf() - 8 * 3600000).toDateString();
}
