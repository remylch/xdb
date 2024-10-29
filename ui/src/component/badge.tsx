import * as React from 'react';
import { PropsWithChildren } from 'react';
import { classNames } from '../utils.ts';

interface IBadge {
  variant?: 'neutral' | 'success' | 'warning' | 'danger';
  onClick?: () => void;
}

const color = {
  neutral: 'text-gray-600 border-gray-600 bg-gray-100',
  success: 'text-green-600 border-green-600 bg-green-100',
  warning: 'text-orange-600 border-orange-600 bg-orange-100',
  danger: 'text-red-600 border-red-600 bg-red-100'
};

export const Badge: React.FC<PropsWithChildren<IBadge>> = ({ variant = 'neutral', children, onClick }) => {
  return (
    <div
      onClick={onClick}
      className={classNames(color[variant], 'px-2 py-1 border rounded-xl font-semibold capitalize w-fit', onClick ? "cursor-pointer": "cursor-default")}>
      {children}
    </div>
  );
};
