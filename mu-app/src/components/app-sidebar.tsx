import * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import { AccountSwitcher } from "@/components/version-switcher";
import { routes } from "@/shared/routes";
import { NavUser } from "./nav-user";
import { NewTicketButton } from "./new-ticket-button";
import { title } from "process";
import { accountStore } from "@/shared/store/account-store";
import { useQuery } from "@tanstack/react-query";
import { accountService } from "@/shared/services/account.service";

// This is sample data.
const data = {
  navMain: [
    {
      title: "Счёт",
      url: routes.dashboard,
      items: [
        {
          title: "Главная",
          url: routes.dashboard.index,
        },
        {
          title: "Пополнить",
          url: routes.dashboard.addFunds,
        },
      ],
    },
    {
      title: "Услуги",
      url: "#",
      items: [],
    },
    {
      title: "Настройки",
      url: "#",
      items: [
        {
          title: "Интеграции",
          url: routes.dashboard.settings.integrations,
        },
      ],
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const accountId = accountStore((s) => s.account!.id);

  const servicesQuery = useQuery({
    queryKey: ["account", "services"],
    queryFn: async () => await accountService.services(accountId),
  });

  if (servicesQuery.data) {
    data.navMain[1].items = servicesQuery.data.services.map((svc) => ({
      title: svc.name,
      url: routes.dashboard.service.dashboard(svc.id),
    }));
  }
  return (
    <Sidebar {...props}>
      <SidebarHeader className="justify-center">
        <AccountSwitcher />
        <NewTicketButton />
      </SidebarHeader>
      <SidebarContent>
        {/* We create a SidebarGroup for each parent. */}
        {data.navMain.map((item) => (
          <SidebarGroup key={item.title}>
            <SidebarGroupLabel>{item.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild isActive={item.isActive}>
                      <a href={item.url}>{item.title}</a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
