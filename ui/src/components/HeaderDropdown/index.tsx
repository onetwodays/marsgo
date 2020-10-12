/*
使用js来动态判断是否为组件添加class（类名）
<Button className={classnames({
    //这里可以根据各属性动态添加，如果属性值为true则为其添加该类名，如果值为false，则不添加。这样达到了动态添加class的目的
      base: true,
      inProgress: this.props.store.submissionInProgress,
      error: this.props.store.errorOccurred,
      disabled: this.props.form.valid,
    })}>
<Button/>
*/

import { DropDownProps } from 'antd/es/dropdown';
import { Dropdown } from 'antd';
import React from 'react';
import classNames from 'classnames'; //react官方推荐的classnames库
import styles from './index.less';

declare type OverlayFunc = () => React.ReactNode;

export interface HeaderDropdownProps extends Omit<DropDownProps, 'overlay'> {
    overlayClassName?: string;
    overlay: React.ReactNode | OverlayFunc | any;
    placement?: 'bottomLeft' | 'bottomRight' | 'topLeft' | 'topCenter' | 'topRight' | 'bottomCenter';
}

const HeaderDropdown: React.FC<HeaderDropdownProps> = ({ overlayClassName: cls, ...restProps }) => (
    <Dropdown overlayClassName={classNames(styles.container, cls)} {...restProps} />
);

export default HeaderDropdown;
