export const routes = {
  home: "/",
  login: "/login",
  requestAccess: "/request-access",

  rate: {
    index: "/",
    list: "/rates",
    create: "/rates/create",
    detail: (id?: string) => (id ? `/rates/${id}` : "/rates/:id"),
  },

  users: {
    index: "/",
    list: "/users",
    detail: (id?: string) => (id ? `/users/${id}` : "/users/:id"),
  },

  accounts: {
    index: "/",
    list: "/accounts",
    detail: (id?: string) => (id ? `/accounts/${id}` : "/accounts/:id"),
  },
};
