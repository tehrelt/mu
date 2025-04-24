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
};
