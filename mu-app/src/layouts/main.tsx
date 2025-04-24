import { Outlet } from "react-router-dom";
import { Header } from "../components/header";

export const MainLayout = () => {
  return (
    <div>
      <Header />
      <main>
        <Outlet />
      </main>
      <footer>
        <p>&copy; 2023 My App</p>
      </footer>
    </div>
  );
};
