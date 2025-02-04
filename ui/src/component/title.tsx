import { PropsWithChildren } from 'react';
import { classNames } from '../utils';

interface ITitle {
  capitalize?: boolean;
  className?: string;
}

export const Title: React.FC<PropsWithChildren<ITitle>> = ({children, capitalize, className}) =>
  <h1 className={classNames(`font-semibold text-4xl text-gray-800 flex items-center ${className}`, capitalize ? "capitalize" : "normal-case")}>{children}</h1>
