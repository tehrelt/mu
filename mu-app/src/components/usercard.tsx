import { fio } from "@/lib/utils";
import { useLogout } from "@/shared/hooks/use-logout";
import { Profile } from "@/shared/types/auth";
import { ChevronDown, House, LogOut } from "lucide-react";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";

type UserCardProps = {
  user: Profile;
};

export const UserCard = ({ user }: UserCardProps) => {
  const { mutate: logout } = useLogout();

  return (
    <div className="flex items-center  gap-2">
      <div className="flex flex-col items-end">
        <span className="text-sm font-medium">{fio(user)}</span>
        <span className="text-xs text-muted-foreground">{user.email}</span>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <Button variant="ghost" size="icon">
            <ChevronDown />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuGroup>
            <DropdownMenuItem>
              <House />
              Личный кабинет
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem onClick={() => logout()}>
              <LogOut />
              Выйти
            </DropdownMenuItem>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
