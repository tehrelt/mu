import { Outlet } from "react-router-dom";
import { Header } from "../components/header";

export const MainLayout = () => {
  return (
    <div className="h-screen flex flex-col">
      <Header />
      <main className="flex-2">
        <Outlet />
      </main>
      <footer className="flex justify-center py-4">
        <p className="italic text-muted-foreground">
          &copy; 2025 ПИ-21а Евтеев Дмитрий
        </p>
      </footer>
    </div>
  );
};
