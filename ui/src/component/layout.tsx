import { Outlet } from 'react-router-dom';

export const Layout = () => {
  return <section className="p-5 flex-1 h-full min-h-[90vh]">
    <Outlet />
  </section>;
};
