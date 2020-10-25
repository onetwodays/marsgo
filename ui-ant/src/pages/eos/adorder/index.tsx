import React, { useRef } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Tag, Space } from 'antd';
import ProTable, { ProColumns, TableDropdown, ActionType } from '@ant-design/pro-table';
import { fetchAll } from '@/services/eos/fetch';
import { msgError } from '@/utils/notify';

interface Adorder {
    "id": number;
    "user": string;
    "pair": string;
    "ctime": string;
    "utime": string;
    "status": number;
    "status_str": string;
    "side": number;
    "type": number;
    "role": number;
    "price": string;
    "amount": string;
    "amount_min": string;
    "amount_max": string;
    "maker_ask_fee_rate": string;
    "maker_bid_fee_rate": string;
    "left": string;
    "freeze": string;
    "deal_fee": string;
    "deal_stock": string;
    "deal_money": string;
    "source": string;
    "vec_deal": number[];
    "pay_accounts": number[];

}

const columns: ProColumns<Adorder>[] = [
    {
        title: 'id',
        dataIndex: 'id', //第一列是序号列，没有标题
        valueType: 'indexBorder',
        width: 48,
    },
    {
        title: '广告主',
        dataIndex: 'user',
        copyable: true,
        ellipsis: true,
        tip: '标题过长会自动收缩',
        formItemProps: {
            rules: [
                {
                    required: true,
                    message: '此项为必填项',
                },
            ],
        },
        //width: '30%',
        search: false,
    },
    {
        title: '交易对',
        dataIndex: 'pair',
        copyable: true,
        ellipsis: true,
        tip: '标题过长会自动收缩',
        formItemProps: {
            rules: [
                {
                    required: true,
                    message: '此项为必填项',
                },
            ],
        },
        //width: '30%',
        search: false,
    },
    {
        title: '状态',
        dataIndex: 'status',
        filters: true,
    },
    {
        title: '成交记录',
        dataIndex: 'vec_deal',
        render: (_, row) => (
            <Space>
                {row.vec_deal.map((id) => (
                    <Tag key={id}>
                        {id}
                    </Tag>
                ))}
            </Space>
        ),
    },
    {
        title: '创建时间',
        key: 'since',
        dataIndex: 'ctime',
        valueType: 'dateTime',
    },
    {
        title: '更新时间',
        key: 'since1',
        dataIndex: 'utime',
        valueType: 'dateTime',
    },

    {
        title: '备注',
        dataIndex: 'source',
        copyable: true,
        ellipsis: true,
        tip: '标题过长会自动收缩',
        formItemProps: {
            rules: [
                {
                    required: true,
                    message: '此项为必填项',
                },
            ],
        },
        //width: '30%',
        search: false,
    },


    {
        title: '操作',
        valueType: 'option',
        render: (text, row, _, action) => [
            <Button title='手动下架'>下架</Button>,
            <TableDropdown
                key="actionGroup"
                onSelect={() => action.reload()}
                menus={[
                    { key: 'copy', name: '复制' },
                    { key: 'delete', name: '删除' },
                ]}
            />,
        ],
    },
];

export default () => {
    const actionRef = useRef<ActionType>();

    return (
        <ProTable<Adorder>
            columns={columns}
            actionRef={actionRef}
            request={async () => {
                let res = await fetchAll('adorders', {
                    scope: 'adxcnymkask',
                });
                console.log("xxx:", res);
                return Promise.resolve({
                    data: res.rows,
                });
            }}
            onRequestError={(e: Error) => msgError(e.message)}
            rowKey="id"
            dateFormatter="string"
            headerTitle="高级表格"
            toolBarRender={() => [
                <Button key="button" icon={<PlusOutlined />} type="primary">
                    新建
        </Button>,
            ]}
        />
    );
};
