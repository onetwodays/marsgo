/*
@umijs/plugin-model 是一种基于 hooks 范式的简单数据流方案，可以在一定情况下替代 dva 来进行中台的全局数据流。
我们约定在 src/models目录下的文件为项目定义的 model 文件。
每个文件需要默认导出一个 function，该 function 定义了一个 Hook，不符合规范的文件我们会过滤掉。
文件名则对应最终 model 的 name，你可以通过插件提供的 API 来消费 model 中的数据。
https://ant.design/docs/react/practical-projects-cn
*/
import { useRequest } from 'umi';
import { message } from 'antd';
//import { queryProductList } from '@/services/product';

export default function useProductList(params: { pageSize: number; current: number }) {
    const msg = useRequest(() => { data: [{ name: "ant" }, { name: "umi" }] }); //发出http请求

    const deleteProducts = async (id: string) => {
        try {
            //await removeProducts(id);
            message.success('success');
            msg.run();
        } catch (error) {
            message.error('fail');
        }
    };

    return {
        namespace: "useProductList",
        dataSource: msg.data,
        reload: msg.run,
        loading: msg.loading,
        deleteProducts,
    };
}