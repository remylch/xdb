import * as React from 'react';
import { PropsWithChildren } from 'react';
import { useNavigate } from 'react-router-dom';

interface IBack {
  to: string
}

export const Back: React.FC<PropsWithChildren<IBack>> = ({children, to}) => {
  const navigate = useNavigate();
  return <div onClick={() => navigate(to)} className="flex gap-2 items-center hover:underline hover:underline-offset-2 cursor-pointer w-fit">
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="h-4 w-4 pt-0.5">
      <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
    </svg>
    <span>
      {children}
    </span>
  </div>
}
