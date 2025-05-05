import { HouseAccount } from "../types/account";
import { create } from "zustand";
import { persist, devtools } from "zustand/middleware";

interface AccountStoreState {
  account?: HouseAccount;
}

interface AccountStoreAction {
  select: (account: HouseAccount) => void;
  clear: () => void;
}

type AccountStore = AccountStoreState & AccountStoreAction;

export const accountStore = create<AccountStore>()(
  persist(
    devtools((set) => ({
      account: undefined,
      select: (account: HouseAccount) => set(() => ({ account })),
      clear: () => set(() => ({ account: undefined })),
    })),
    { name: "account" }
  )
);
