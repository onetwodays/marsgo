import React, { Component } from 'react';
//导入antd布局相关的组件
import { Layout,
         Menu,
         Switch, 
         Divider 
        } from 'antd';


//导入antd-design-icons的图标

import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  UserOutlined,
  VideoCameraOutlined,
  UploadOutlined,
  DesktopOutlined,
} from '@ant-design/icons';

//导入umi的Link组件
import {Link} from 'umi'

// Header, Footer, Sider, Content组件在Layout组件模块下 对象的解析赋值
const { Header, Footer, Sider, Content } = Layout; 
const SubMenu = Menu.SubMenu;


const menuData = [
  { route:'/hello',                   name:'欢迎'},            //0
  { route:'/dashboard/base',          name:'基本信息'},        //1
  { route:'/dashboard/monitor',       name:'监控'},           //2
  { route:'/dashboard/workplace',     name:'工作台'},         //3
  { route:'/demo',                    name:'示例'},          //4
  { route:'/test',                    name:'test'},          //4

]

//定义多个菜单数组
const subMenuHero = [
  { route:'/hero/hero',                        name:'英雄'},            //0
  { route:'/hero/item',                         name:'局内道具'},            //0
  { route:'/hero/summoner',                     name:'召唤师技能'},            //0

]







//umi 会自动使用BasicLayout包裹页面，并传入如下 props
//{
//  match?: match<P>;
//  location: Location<S>;
//  history: History;
//  route: IRoute;
//}

class BasicLayout extends Component {
  constructor(props){
    super(props);
    this.state ={
      collapsed:false,
      mode:  'inline',
      theme: 'dark',
    };
  }
  
  

  //菜单是不是可收缩的
  onCollapse = collapsed =>{
    this.setState({collapsed})
  };

  //菜单是水平还是垂直
  onChangeMode= value=>{
    this.setState({
      mode:value?'vertical':'inline',
    });
  };

  //改变菜单的主题
  onChangeTheme= value=>{
    this.setState({theme:value?'dark':'light',});
  };
  
  // location:{pathname} 解析赋值，location是模式，pathname才是要赋值的变量
  render() {
    const {
      location:{pathname}, 
      children,
    } = this.props;

    //const location = this.props.location;
    //const pathname = location.pathname; //note:pathname是个数组


    return (
      <Layout>
        <Sider width={256} style={{ minHeight: '100vh'}} collapsible collapsed={this.state.collapsed} onCollapse={this.onCollapse}>
          <div style={{ height: '32px', background: 'rgba(255,255,255,.2)', margin: '16px' }} />
          
          <Menu theme={this.state.theme} mode={this.state.mode} defaultSelectedKeys={[pathname]}>


            <SubMenu key="subHero" title={<span><DesktopOutlined/><span>王者荣耀资料库</span></span>}>
              {subMenuHero.map(menu=>(
                <Menu.Item key={`/${menu.route}`}>
                  <Link to={menu.route}>{menu.name}</Link>
                </Menu.Item>
                
              ))}
            </SubMenu>

            <Menu.Item icon={<UserOutlined />} key={`/${menuData[0].route}`}>
            <Link to={menuData[0].route}>{menuData[0].name}</Link> 
            </Menu.Item>



            <SubMenu key="sub1" title={<span><DesktopOutlined/><span>仪表盘</span></span>} >
              <Menu.Item key={`/${menuData[1].route}`}><Link to={menuData[1].route}>{menuData[1].name}</Link></Menu.Item>
              <Menu.Item key={`/${menuData[2].route}`}><Link to={menuData[2].route}>{menuData[2].name}</Link></Menu.Item>
              <Menu.Item key={`/${menuData[3].route}`}><Link to={menuData[3].route}>{menuData[3].name}</Link></Menu.Item>
            </SubMenu>

            <Menu.Item icon={<MenuFoldOutlined />} key={`/${menuData[4].route}`}>
                <Link to={menuData[4].route}>{menuData[4].name}</Link>
            </Menu.Item>

            <Menu.Item  key={`/${menuData[5].route}`}>
                <Link to={menuData[5].route}>{menuData[5].name}</Link>
            </Menu.Item>

          </Menu>
        </Sider>
        <Layout>
          <Header  style={{ background: '#fff', textAlign: 'center', padding: 24 }}>hello world!</Header>
          <Content style={{ margin: '24px 16px 0' }}>
            <div style={{ padding: 24, background: '#fff', minHeight: 360 }}>
              {children}
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>ZHworker Design ©2018 Created by iliu</Footer>
        </Layout>
      </Layout>
    )
  }
}

export default BasicLayout;