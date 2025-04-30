export const routes = {
  home: "/",
  signIn: "/sign-in",
  signUp: "/sign-up",
  forgotPassword: "/forgot-password",

  dashboard: {
    index: "/dashboard",
    newTicket: "/dashboard/new-ticket",
    addFunds: "/dashboard/add-funds",
    account: {
      transactionHistory: "/dashboard/account/transaction-history",
    },
    settings: {
      integrations: "/dashboard/settings/integrations",
    },
  },

  billing: {
    index: "/billing",
    process: (id?: string) =>
      id ? `/billing/process/${id}` : "/billing/process/:id",
  },
};
