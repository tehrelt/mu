import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
} from "lucide-react";
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

// This is sample data.
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
          title: "Запросы на доступ",
          url: "#",
        },
        {
          title: "Список",
          url: "#",
        },
      ],
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="text-xl text-center py-4 font-bold">
        Admin MoiUslugi
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
