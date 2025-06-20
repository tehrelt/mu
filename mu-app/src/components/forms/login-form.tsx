import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { routes } from "@/shared/routes";
import { authService } from "@/shared/services/auth.service";
import { LoginInput, loginSchema } from "@/shared/types/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { GalleryVerticalEnd } from "lucide-react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Form, FormField, FormLabel } from "../ui/form";

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
      navigate(routes.dashboard.index);
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
          <CardTitle className="text-xl">С возвращением!</CardTitle>
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
    </div>
  );
}
