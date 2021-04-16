import React from 'react';
import { Carousel, Card, Collapse, Empty, Image, Button, List, Divider, Typography, Popover } from 'antd';
import { Statistic, Row, Col, Table, Timeline, Tooltip, Alert, message, BackTop } from 'antd';
import { Form, Input, Checkbox } from 'antd';

import img1 from '../imag/3.jpeg';
import img2 from '../imag/4.jpeg';
import img3 from '@/assets/image/5.jpg';
/*
import { useModel } from 'umi';
import ProductList from '@/components/ProductList';


const Products = () => {
    const { dataSource, reload, deleteProducts } = useModel('useProductList');
    return (
        <div>
            <a onClick={() => reload()}>reload</a>
            <ProductList onDelete={deleteProducts} products={dataSource} />
        </div>
    );
};

export default Products;
*/

const contentStyle = {
        height: '160px',
        color: '#fff',
        lineHeight: '160px',
        textAlign: 'center',
        background: '#364d79',
};

const { Panel } = Collapse;

const text = `
  A dog is a type of domesticated animal.
  Known for its loyalty and faithfulness,
  it can be found as a welcome guest in many households across the world.
`;

const data = [
        'Racing car sprays burning fuel into crowd.',
        'Japanese princess to wed commoner.',
        'Australian walks 100km after outback crash.',
        'Man charged over missing wedding girl.',
        'Los Angeles battles huge wildfires.',
];

const content = (
        <div>
                <p>Content</p>
                <p>Content</p>
        </div>
);

const dataSource = [
        {
                key: '1',
                name: '胡彦斌',
                age: 32,
                address: '西湖区湖底公园1号',
        },
        {
                key: '2',
                name: '胡彦祖',
                age: 42,
                address: '西湖区湖底公园1号',
        },
];

const columns = [
        {
                title: '姓名',
                dataIndex: 'name',
                key: 'name',
        },
        {
                title: '年龄',
                dataIndex: 'age',
                key: 'age',
        },
        {
                title: '住址',
                dataIndex: 'address',
                key: 'address',
        },
];


const success = () =>
{
        message.success('This is a prompt message for success, and it will disappear in 10 seconds', 10);
};




export default () =>
{
        return (
                <div>
                        <Carousel autoplay >
                                <div>
                                        <Image style={contentStyle} src="/image/2.jpeg" />
                                </div>
                                <div>
                                        <Image style={contentStyle} src={img1} />
                                </div>
                                <div>
                                        <Image style={contentStyle} src={img2} />
                                </div>
                                <div>
                                        <Image style={contentStyle} src={img3} />
                                </div>
                        </Carousel>
                        <Alert message="Success Text" type="success" closable />
                        <Button onClick={success}>Customized display duration</Button>


                        <Card title="曹雪" bordered={false} style={{ width: 300 }}>
                                <p>Card content</p>
                                <p>Card content</p>
                                <p>Card content</p>
                        </Card>

                        <Collapse defaultActiveKey={['1']} onChange={() => { }}>
                                <Panel header="This is panel header 1" key="1">
                                        <p>{text}</p>
                                </Panel>
                                <Panel header="This is panel header 2" key="2">
                                        <p>{text}</p>
                                </Panel>
                                <Panel header="This is panel header 3" key="3" disabled>
                                        <p>{text}</p>
                                </Panel>
                        </Collapse>

                        <Empty
                                image="https://gw.alipayobjects.com/zos/antfincdn/ZHrcdLPrvN/empty.svg"
                                imageStyle={{
                                        height: 60,
                                }}
                                description={
                                        <span>
                                                Customize <a href="#API">Description</a>
                                        </span>
                                }
                        >
                                <Button type="primary">Create Now</Button>
                        </Empty>,

                        <Divider orientation="left">Default Size</Divider>
                        <List
                                header={<div>Header</div>}
                                footer={<div>Footer</div>}
                                bordered
                                dataSource={data}
                                renderItem={item => (
                                        <List.Item>
                                                <Typography.Text mark>[ITEM]</Typography.Text> {item}
                                        </List.Item>
                                )}
                        />
                        <Divider orientation="left">Small Size</Divider>

                        <Popover content={content} title="Title">
                                <Button type="primary">Hover me</Button>
                        </Popover>

                        <Timeline>
                                <Timeline.Item>创建服务现场 2015-09-01</Timeline.Item>
                                <Timeline.Item>初步排除网络异常 2015-09-01</Timeline.Item>
                                <Timeline.Item>技术测试异常 2015-09-01</Timeline.Item>
                                <Timeline.Item>网络异常正在修复 2015-09-01</Timeline.Item>
                        </Timeline>

                        <Tooltip title="prompt text">
                                <span>Tooltip will show on mouse enter.</span>
                        </Tooltip>

                        <BackTop />
                        Scroll down to see the bottom-right
                        <strong className="site-back-top-basic"> gray </strong>
                        button.


                </div>
        );

};
