import { Outlet } from 'react-router-dom';
import {Title} from "./title.tsx";

export const Layout = () => {
  return <section className="p-5 flex-1 h-full min-h-[90vh]">
    <Title>XDB</Title>
    <Outlet />
  </section>;
};
