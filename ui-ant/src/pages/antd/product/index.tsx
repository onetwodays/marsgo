import React from 'react';
import { useModel } from 'umi'

import ProductList from '@/components/ProductList'

const Products = () => {
    const { dataSource, reload, deleteProducts } = useModel('useProductList');
    return (
        <div>
            <a onClick={() => reload()}>reload</a>
            <ProductList onDelete={deleteProducts} products={dataSource} />
        </div>
    );
};

export default Products

