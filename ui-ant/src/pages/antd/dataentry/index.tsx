import React from 'react';

import { Form, Input, Checkbox, Button, InputNumber, Mentions, Rate, PageHeader } from 'antd';

const layout = {
    labelCol: { span: 8 },
    wrapperCol: { span: 16 },
};
const tailLayout = {
    wrapperCol: { offset: 8, span: 16 },
};

function onChange(value) {
    console.log('changed', value);
}
function onSelect(option) {
    console.log('select', option);
}

const { Option } = Mentions;

const Demo = () => {
    const onFinish = values => {
        console.log('Success:', values);
    };

    const onFinishFailed = errorInfo => {
        console.log('Failed:', errorInfo);
    };

    return (

        <div>
            <PageHeader
                className="site-page-header"
                onBack={() => null}
                title="Title"
                subTitle="This is a subtitle"
            />
            <div>
                <Rate allowHalf defaultValue={2.5} />
            </div>

            <Mentions
                style={{ width: '100%' }}
                onChange={onChange}
                onSelect={onSelect}
                defaultValue="@afc163"
            >
                <Option value="afc163">afc163</Option>
                <Option value="zombieJ">zombieJ</Option>
                <Option value="yesmeck">yesmeck</Option>
            </Mentions>

            <Form
                {...layout}
                name="basic"
                initialValues={{ remember: true }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
            >
                <Form.Item
                    label="Username"
                    name="username"
                    rules={[{ required: true, message: 'Please input your username!' }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Password"
                    name="password"
                    rules={[{ required: true, message: 'Please input your password!' }]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item {...tailLayout} name="remember" valuePropName="checked">
                    <Checkbox>Remember me</Checkbox>
                </Form.Item>

                <Form.Item>
                    <InputNumber min={1} max={10} defaultValue={3} onChange={onChange} />
                </Form.Item>

                <Form.Item {...tailLayout}>
                    <Button type="primary" htmlType="submit">
                        Submit
          </Button>
                </Form.Item>
            </Form>
        </div>

    );
};

export default Demo;