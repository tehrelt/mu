import { AppSidebar } from "@/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useTitle } from "@/shared/hooks/use-title";
import { accountStore } from "@/shared/store/account-store";
import { Outlet } from "react-router-dom";

export const DashboardLayout = () => {
  useTitle("MU: Dashboard");

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="p-6">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
};
export const AccountCheck = () => {
  const account = accountStore((s) => s.account);

  return account ? <Outlet /> : <NoAccount />;
};

const NoAccount = () => {
  return (
    <div className="flex justify-center items-center h-screen">
      <p className="text-xl font-black">Выберите счёт</p>
    </div>
  );
};
