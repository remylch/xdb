import * as React from 'react';
import { PropsWithChildren } from 'react';

interface IButton extends React.ButtonHTMLAttributes<HTMLButtonElement> {}

export const Button: React.FC<PropsWithChildren<IButton>> = ({children, ...props}) =>
  <button
    type="button"
    className="rounded-md bg-gray-800 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-gray-900 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-500"
    {...props}
  >
    {children}
  </button>
