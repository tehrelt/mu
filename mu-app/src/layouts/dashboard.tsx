import { AppSidebar } from "@/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { accountStore } from "@/shared/store/account-store";
import { Outlet } from "react-router-dom";

export const DashboardLayout = () => {
  const account = accountStore((s) => s.account);

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="p-6">{account ? <Outlet /> : <NoAccount />}</div>
      </SidebarInset>
    </SidebarProvider>
  );
};

const NoAccount = () => {
  return (
    <div className="flex justify-center items-center h-screen">
      <p className="text-xl font-black">Выберите аккаунт</p>
    </div>
  );
};
