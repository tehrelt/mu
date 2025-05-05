import { ChevronsUpDown, GalleryVerticalEnd } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useAccounts } from "@/shared/hooks/use-accounts";
import { accountStore } from "@/shared/store/account-store";
import { Balance } from "./ui/balance";
import { useNavigate } from "react-router-dom";
import { routes } from "@/shared/routes";

export function AccountSwitcher() {
  const accounts = useAccounts();

  const selectedAccount = accountStore((s) => s.account);
  const selectAccount = accountStore((s) => s.select);
  const clearAccount = accountStore((s) => s.clear);

  const navigate = useNavigate();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <SidebarMenuButton
            size="lg"
            className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <DropdownMenuTrigger className="flex w-full items-center gap-x-2">
              <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                <GalleryVerticalEnd className="size-4" />
              </div>
              <div className="flex flex-col gap-0.5 leading-none">
                {selectedAccount ? (
                  <div className="flex flex-col items-start">
                    <p className="font-medium">
                      {selectedAccount.house.address}
                    </p>
                  </div>
                ) : (
                  <span className="font-medium">Select Account</span>
                )}
              </div>
              <ChevronsUpDown className="ml-auto" />
            </DropdownMenuTrigger>
          </SidebarMenuButton>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width)"
            align="start"
          >
            {!accounts.isFetching &&
              accounts.data &&
              accounts.data.accounts
                .filter(
                  (acc) =>
                    (selectedAccount && selectedAccount.id !== acc.id) ||
                    !selectedAccount,
                )
                .map((acc) => (
                  <DropdownMenuItem
                    key={acc.id}
                    onSelect={() => {
                      selectAccount(acc);
                      navigate(routes.dashboard.index);
                    }}
                  >
                    {acc.house.address}
                  </DropdownMenuItem>
                ))}

            {selectedAccount && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onSelect={() => {
                    clearAccount();
                    navigate(routes.dashboard.index);
                  }}
                >
                  Clear
                </DropdownMenuItem>
              </>
            )}

            {accounts.isError && (
              <DropdownMenuItem disabled>
                Error when fetching accounts {accounts.error.message}
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
      {selectedAccount && (
        <SidebarMenuItem className="flex justify-end items-center gap-2">
          Баланс: <Balance balance={selectedAccount.balance} />
        </SidebarMenuItem>
      )}
    </SidebarMenu>
  );
}
