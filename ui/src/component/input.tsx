import * as React from 'react';
import { PropsWithChildren } from 'react';

interface ITextInput extends React.InputHTMLAttributes<HTMLInputElement> {}

export const TextInput: React.FC<PropsWithChildren<ITextInput>> = ({ children, ...props}) =>
  <input className="rounded-sm border-2 px-3 py-1.5 focus: outline-none" {...props}>{children}</input>
