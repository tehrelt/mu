import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { LoginInput, loginSchema } from "@/shared/types/auth";
import { Form, FormField, FormLabel } from "./ui/form";
import { useMutation } from "@tanstack/react-query";
import { authService } from "@/shared/services/auth.service";
import { Link, useNavigate } from "react-router-dom";
import { routes } from "@/shared/routes";
import { toast } from "sonner";
import { GalleryVerticalEnd } from "lucide-react";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const form = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      login: "",
      password: "",
    },
  });

  const navigate = useNavigate();

  const { mutate: login, isPending } = useMutation({
    mutationKey: ["login"],
    mutationFn: async (data: LoginInput) => await authService.login(data),
    onSuccess: () => {
      navigate(routes.home);
    },
    onError: (e) => {
      console.log(e);
      toast.error("Не удалось войти", {
        description: "Проверьте правильность введенных данных",
      });
    },
  });

  const submit = async (data: LoginInput) => {
    login(data);
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Link
        to={routes.home}
        className="flex items-center gap-2 self-center font-medium"
      >
        <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <GalleryVerticalEnd className="size-4" />
        </div>
        Мои услуги
      </Link>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
          <CardDescription>
            lorem ipsum ojgewopg -3 3i2p ehg e lk'ge
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(submit, console.warn)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <FormField
                    control={form.control}
                    name="login"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="email">Email</FormLabel>
                        <Input {...field} />
                      </div>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <div className="flex items-center">
                          <Label htmlFor="password">Пароль</Label>
                          <Link
                            to={routes.forgotPassword}
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                          >
                            Забыли пароль?
                          </Link>
                        </div>
                        <Input {...field} type="password" />
                      </div>
                    )}
                  />
                  <Button type="submit" className="w-full" disabled={isPending}>
                    Login
                  </Button>
                </div>
                <div className="text-center text-sm">
                  Нет аккаунта?{" "}
                  <Link
                    to={routes.signUp}
                    className="underline underline-offset-4"
                  >
                    Зарегистрироваться
                  </Link>
                </div>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        By clicking continue, you agree to our <a href="#">Terms of Service</a>{" "}
        and <a href="#">Privacy Policy</a>.
      </div>
    </div>
  );
}
