import { Outlet } from 'react-router-dom';
import {Title} from "./title.tsx";
import LogViewer from './logviewer.tsx';

export const Layout = () => {
  return <section className="flex gap-5 h-screen">
    <div className="p-5 flex flex-col gap-5 flex-1 h-full">
      <Title>XDB</Title>
      <Outlet />
    </div>
    <LogViewer />
  </section>;
};
