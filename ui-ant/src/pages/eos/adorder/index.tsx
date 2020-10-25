import React, { useRef } from 'react';
import { PlusOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { Button, Tag, Space, Tooltip } from 'antd';
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
        dataIndex: 'id',
        // valueType: 'indexBorder',
        width: 48,
        // render: (_) => <a>{_}</a>,
        sorter: (a, b) => a.id - b.id,
    },
    {
        title: '广告主',
        dataIndex: 'user',
        copyable: true,
        ellipsis: false,
        formItemProps: {
            rules: [
                {
                    required: true,
                    message: '此项为必填项',
                },
            ],
        },
        // width: '30%',
        search: true,
    },
    {
        title: '交易对',
        dataIndex: 'pair',
        copyable: true,
        // ellipsis: true,
        tip: '标题过长会自动收缩',
        formItemProps: {
            rules: [
                {
                    required: true,
                    message: '此项为必填项',
                },
            ],
        },
        // width: '30%',
        search: false,
    },
    {
        title: '买卖',
        dataIndex: 'side',
        filters: true,
        width: 48,
    },
    {
        title: '价格',
        dataIndex: 'price',
        // filters: true,
        // width: 48,
    },
    {
        title: '数量',
        dataIndex: 'amount',
        // filters: true,
        // width: 48,
    },

    {
        title: '剩余',
        dataIndex: 'left',
        // filters: true,
        // width: 48,
    },

    {
        title: '待放币',
        dataIndex: 'freeze',
        // filters: true,
        // width: 48,
    },
    {
        title: '状态',
        dataIndex: 'status',
        filters: true,
        width: 48,
    },
    {
        title: '状态描述',
        dataIndex: 'status_str',
        filters: true,
        ellipsis: true,

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
        search: false,
        ellipsis: true,
    },
    {
        title: '付款方式',
        dataIndex: 'pay_accounts',
        render: (_, row) => (
            <Space>
                {row.pay_accounts.map((id) => (
                    <Tag key={id}>
                        {id}
                    </Tag>
                ))}
            </Space>
        ),
        search: false,
        ellipsis: true,
    },
    {
        title: (
            <>
                创建时间
                <Tooltip placement="top" title="这是一段描述">
                    <QuestionCircleOutlined style={{ marginLeft: 4 }} />
                </Tooltip>
            </>
        ),
        key: 'since',
        dataIndex: 'ctime',
        valueType: 'dateTime',
        search: false,
    },
    {
        title: '更新时间',
        key: 'since1',
        dataIndex: 'utime',
        valueType: 'dateTime',
        search: false,
    },

    {
        title: '备注',
        dataIndex: 'source',
        copyable: true,
        ellipsis: true,
        tip: '标题过长会自动收缩',
        search: false,
        width: '20%',
    },


    {
        title: '操作',
        key: 'option',
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
            rowSelection={{}}
            tableAlertRender={({ selectedRowKeys, selectedRows, onCleanSelected }) => (
                <Space size={24}>
                    <span>
                        已选 {selectedRowKeys.length} 项
                    <a style={{ marginLeft: 8 }} onClick={onCleanSelected}>
                            取消选择
                    </a>
                    </span>
                    <span>{`总成交记录: ${selectedRows.reduce(
                        (pre, item) => pre + item.vec_deal.length,
                        0,
                    )} 个`}</span>

                </Space>
            )}
            tableAlertOptionRender={() => {
                return (
                    <Space size={16}>
                        <a>批量删除</a>
                        <a>导出数据</a>
                    </Space>
                );
            }}
            request={async () => {
                const res = await fetchAll('adorders', {
                    scope: 'adxcnymkask',
                });
                console.log("xxx:", res);
                return Promise.resolve({
                    data: res.rows,
                });
            }}
            onRequestError={(e: Error) => msgError(e.message)}
            rowKey="id"
            pagination={{
                showQuickJumper: true,
                pageSize: 20,
            }}
            search={{
                // filterType: 'light',

            }}
            scroll={{ x: 1300 }}
            options={false}
            dateFormatter="string"
            headerTitle="广告表格"
            toolBarRender={() => [
                <Button key="new" icon={<PlusOutlined />} type="primary">
                    新建
                </Button>,
            ]}
        />
    );
};
