import { Settings as LayoutSettings } from '@ant-design/pro-layout';

export default {
	logo: 'https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg', //layout左上角logo的url
	title: 'together',  //layout左上角的title
	navTheme: 'light', // 拂晓蓝,导航主题 light/dark
	primaryColor: '#1890ff',
	layout: 'mix', // 顶部显示1级菜单，左侧显示二三级菜单
	splitMenus: true, // 需要注意，当 mix 模式时，需要添加splitMenus: true，顶部才可以正确展示一级菜单
	locale: true, //通过 layout 配置的 locale 配置开启国际化，开启后路由里配置的菜单名会被当作菜单名国际化的 key，插件会去 locales 文件中查找 menu.[key]对应的文案，默认值为改 key

	contentWidth: 'Fluid', //layout的内容模式，Fluid定宽，Fixed:自适应
	fixedHeader: false, //是否固定header到顶部
	fixSiderbar: true, //是否固定导航
	colorWeak: false,
	menu: {            //关于菜单的配置，目前只支持local，可以关闭menu自带的全球化
		locale: true,
	},
	pwa: false,

	iconfontUrl: '',
} as LayoutSettings & {
	pwa?: boolean;
};
