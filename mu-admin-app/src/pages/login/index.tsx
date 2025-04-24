import { LoginForm } from "@/components/login-form";

export const LoginPage = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen ">
      <div className=" p-8 rounded-lg shadow-md">
        <LoginForm />
      </div>
    </div>
  );
};
