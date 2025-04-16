import { LoginForm } from "~/components/login-form";
import { ActionFunctionArgs, MetaFunction, redirect } from "@remix-run/node";
import { loginRequestSchema } from "~/schemes/login";
import { useActionData } from "@remix-run/react";
import { sessionService } from "~/services/session.server";
import { AuthService } from "~/services/auth";

export const meta: MetaFunction = () => {
  return [{ title: "Вход в аккаунт" }, { description: "Вход в аккаунт" }];
};

export async function action({ request }: ActionFunctionArgs) {
  const formData = await request.formData();
  const raw = Object.fromEntries(formData);
  const result = loginRequestSchema.safeParse(raw);

  if (!result.success) {
    return Response.json({ errors: result.error.flatten().fieldErrors });
  }

  const data = result.data;

  try {
    const authService = new AuthService();
    const tokensPair = await authService.login(data);
    return redirect("/", {
      headers: {
        "Set-Cookie": await sessionService.set(tokensPair),
      },
    });
  } catch (error) {
    return Response.json({ errors: { email: ["Неверный email или пароль"] } });
  }

  return null;
}

export default function Login() {
  const actionData = useActionData<typeof action>();
  console.log("Action data", actionData);

  return (
    <div className="flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <LoginForm />
        {/* {actionData?.errors && (
          <div className="mt-4 text-red-500">{actionData.errors}</div>
        )} */}
      </div>
    </div>
  );
}
