import { LoaderFunctionArgs } from "@remix-run/node";
import { Outlet, useLoaderData } from "@remix-run/react";
import { Header } from "~/components/header";
import { AuthService } from "~/services/auth";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const cookie = request.headers.get("Cookie");
  const auther = new AuthService(cookie);
  const profile = await auther.profile();
  return { profile };
};

export default function Layout() {
  const { profile } = useLoaderData<typeof loader>();

  return (
    <div>
      <Header profile={profile} />
      <main>
        <Outlet />
      </main>
    </div>
  );
}
