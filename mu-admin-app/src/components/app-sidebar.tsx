import { Bot, SquareTerminal } from "lucide-react";
import * as React from "react";
import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { routes } from "@/shared/routes";

import ThemeSwitcher from "./theme-switcher";

const data = {
  navMain: [
    {
      title: "Тарифы",
      url: "#",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "Создать",
          url: routes.rate.create,
        },
        {
          title: "Список",
          url: routes.rate.list,
        },
      ],
    },
    {
      title: "Пользователи",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Список",
          url: routes.users.list,
        },
      ],
    },
    {
      title: "Заявки",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Список",
          url: routes.tickets.list,
        },
      ],
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="text-xl text-center py-4 font-bold">
        <div>
          <p>&quot;Мои услуги&quot;</p>
          <p>Админ панель</p>
        </div>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <ThemeSwitcher />
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
