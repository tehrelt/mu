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
import { useConnectedServices } from "@/shared/hooks/use-services";
import { routes } from "@/shared/routes";
import { NavUser } from "./nav-user";
import { NewTicketButton } from "./new-ticket-button";
import ThemeSwitcher from "./theme-switcher";

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
  const servicesQuery = useConnectedServices();
  if (servicesQuery.data) {
    data.navMain[1].items = servicesQuery.data.services.map((svc) => ({
      title: svc.name,
      url: routes.dashboard.cabinet.dashboard(svc.cabinetId),
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
                    {/*@ts-ignore  */}
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
        <ThemeSwitcher />
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
