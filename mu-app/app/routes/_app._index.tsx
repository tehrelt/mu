import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { AuthService } from "~/services/auth";

export const meta: MetaFunction = () => {
  return [
    { title: "МоиУслуги" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const auther = new AuthService(request.headers.get("cookie"));
  const profile = await auther.profile();
  return { profile };
};

export default function Index() {
  const { profile } = useLoaderData<typeof loader>();

  return (
    <div className="">
      INDEX PAGE
      {profile && <div>for {profile.email}</div>}
    </div>
  );
}
