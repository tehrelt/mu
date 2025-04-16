import {
  ActionFunctionArgs,
  LoaderFunctionArgs,
  redirect,
} from "@remix-run/node";

export const action = async ({ request }: ActionFunctionArgs) =>
  logout(request);

export const loader = async ({ request }: LoaderFunctionArgs) => {
  redirect("/");
};
