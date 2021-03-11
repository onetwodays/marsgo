// 随着应用的发展，你会需要在多个页面分享 UI 元素(或在一个页面使用多次) ，在 umi 里你可以把这部分抽成 component 。

// 我们来编写一个 ProductList component，这样就能在不同的地方显示产品列表了


import React from 'react';
import { Table, Popconfirm, Button } from 'antd';

const ProductList: React.FC<{ products: { name: string }[]; onDelete: (id: string) => void }> = ({
    onDelete,
    products,
}) => {
    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
        },
        {
            title: 'Actions',
            render: (text, record) => {
                return (
                    <Popconfirm title="Delete?" onConfirm={() => onDelete(record.id)}>
                        <Button>Delete</Button>
                    </Popconfirm>
                );
            },
        },
    ];
    return <Table dataSource={products} columns={columns} />;
};

export default ProductList;