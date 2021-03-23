/**
 * 在生产环境 代理是无法生效的，所以这里没有生产环境的配置
 * The agent cannot take effect in the production environment
 * so there is no configuration of the production environment
 * For details, please see
 * https://pro.ant.design/docs/deploy
 * 代理不会修改控制台的 url，它的所有操作都在 node.js 中进行。
 * //target: 'http://127.0.0.1:18888',
 */
export default {
	dev: {
		'/api/': {
			target: 'https://preview.pro.ant.design',
			changeOrigin: true,
			pathRewrite: { '^': '' },
		},
		'/v1/': {
			target: 'http://zhongyingying.qicp.io:38000',
			changeOrigin: true,
			pathRewrite: { '^': '' },
		},
	},
	test: {
		'/api/': {
			target: 'https://preview.pro.ant.design',
			changeOrigin: true,
			pathRewrite: { '^': '' },
		},
	},
	pre: {
		'/api/': {
			target: 'your pre url',
			changeOrigin: true,
			pathRewrite: { '^': '' },
		},
	},
};
