import { PropsWithChildren } from 'react';
import { classNames } from '../utils';

interface ITitle {
  capitalize?: boolean;
}

export const Title: React.FC<PropsWithChildren<ITitle>> = ({children, capitalize}) =>
  <h1 className={classNames('font-semibold text-4xl text-gray-800 flex items-center', capitalize ? "capitalize" : "normal-case")}>{children}</h1>
