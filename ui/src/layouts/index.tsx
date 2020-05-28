import React, { Component } from 'react';
import { Layout,
         Menu,
         Switch, 
         Divider 
        } from 'antd';
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  UserOutlined,
  VideoCameraOutlined,
  UploadOutlined,
  DesktopOutlined,
} from '@ant-design/icons';

import {Link} from 'umi'

// Header, Footer, Sider, Content组件在Layout组件模块下
const { Header, Footer, Sider, Content } = Layout;
const SubMenu = Menu.SubMenu;


const menuData = [
  { route:'/hello',              name:'欢迎'},        //0

  { route:'/dashboard/base',    name:'基本信息'},      //1
  { route:'/dashboard/monitor',  name:'监控'},        //2
  { route:'/dashboard/workplace',name:'工作台'},      //3
  { route:'/demo',               name:'示例'},      //4

]

//umi 会自动使用BasicLayout包裹页面，并传入如下 props
//{
//  match?: match<P>;
//  location: Location<S>;
//  history: History;
//  route: IRoute;
//}

class BasicLayout extends Component {
  state ={
    collapsed:false,
    mode:  'inline',
    theme: 'dark',
  };

  onCollapse = collapsed =>{
    this.setState({collapsed})
  };

  onChangeMode= value=>{
    this.setState({
      mode:value?'vertical':'inline',
    });
  };

  onChangeTheme= value=>{
    this.setState({theme:value?'dark':'light',});
  };
  

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

          </Menu>
        </Sider>
        <Layout>
          <Header style={{ background: '#fff', textAlign: 'center', padding: 24 }}>hello world!</Header>
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