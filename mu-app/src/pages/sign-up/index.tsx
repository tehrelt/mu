import { LoginForm } from "@/components/login-form";
import { RegisterForm } from "@/components/register-form";

export const SignUpPage = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen ">
      <div className=" p-8 rounded-lg shadow-md">
        <RegisterForm />
      </div>
    </div>
  );
};
