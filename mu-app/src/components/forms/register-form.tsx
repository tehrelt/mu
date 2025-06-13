import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { routes } from "@/shared/routes";
import { authService } from "@/shared/services/auth.service";
import { RegisterInput, registerSchema } from "@/shared/types/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { GalleryVerticalEnd } from "lucide-react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Form, FormField, FormLabel, FormMessage } from "../ui/form";

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const form = useForm<RegisterInput>({
    //@ts-ignore
    resolver: zodResolver(registerSchema),
  });

  const navigate = useNavigate();

  const { mutate: register, isPending } = useMutation({
    mutationKey: ["register"],
    mutationFn: async (data: RegisterInput) => await authService.register(data),
    onSuccess: () => {
      navigate(routes.dashboard.index);
    },
    onError: (e) => {
      console.log(e);
      toast.error("Не удалось зарегистрироваться", {
        description: "Проверьте правильность введенных данных",
      });
    },
  });

  const submit = async (data: RegisterInput) => {
    register(data);
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
          <CardTitle className="text-xl">Добро пожаловать</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            {/*@ts-ignore */}
            <form onSubmit={form.handleSubmit(submit, console.warn)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <div className="grid grid-cols-3 gap-x-2">
                    <FormField
                      //@ts-ignore
                      control={form.control}
                      name="lastName"
                      render={({ field }) => (
                        <div className="grid gap-3">
                          <FormLabel htmlFor="lastName">Фамилия</FormLabel>
                          <Input {...field} placeholder="Иванов" />
                          <FormMessage />
                        </div>
                      )}
                    />
                    <FormField
                      //@ts-ignore
                      control={form.control}
                      name="firstName"
                      render={({ field }) => (
                        <div className="grid gap-3">
                          <FormLabel htmlFor="firstName">Имя</FormLabel>
                          <Input {...field} placeholder="Иван" />
                          <FormMessage />
                        </div>
                      )}
                    />
                    <FormField
                      //@ts-ignore
                      control={form.control}
                      name="middleName"
                      render={({ field }) => (
                        <div className="grid gap-3">
                          <FormLabel htmlFor="middleName">Отчество</FormLabel>
                          <Input {...field} placeholder="Иванович" />
                          <FormMessage />
                        </div>
                      )}
                    />
                  </div>
                  <FormField
                    //@ts-ignore
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="email">Email</FormLabel>
                        <Input {...field} placeholder="example@mail.ru" />
                        <FormMessage />
                      </div>
                    )}
                  />
                  <FormField
                    //@ts-ignore
                    control={form.control}
                    name="phone"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="phone">Телефон (без +7)</FormLabel>
                        <div className="flex items-center">
                          <div className="flex h-full items-center px-2 rounded-l-md border bg-accent">
                            +7
                          </div>
                          <Input
                            {...field}
                            className="rounded-l-none border-l-0"
                            placeholder="94912341212"
                          />
                        </div>
                        <FormMessage />
                      </div>
                    )}
                  />
                  <div className="space-y-2">
                    <p className="font-bold">Паспорт</p>
                    <div className="grid grid-cols-2 gap-x-2">
                      <FormField
                        //@ts-ignore
                        control={form.control}
                        name="passport.series"
                        render={({ field }) => (
                          <div className="grid gap-3">
                            <FormLabel htmlFor="passport.series">
                              Серия (4 цифры)
                            </FormLabel>
                            <Input type="number" {...field} />
                            <FormMessage />
                          </div>
                        )}
                      />
                      <FormField
                        //@ts-ignore
                        control={form.control}
                        name="passport.number"
                        render={({ field }) => (
                          <div className="grid gap-3">
                            <FormLabel htmlFor="passport.number">
                              Номер (6 цифр)
                            </FormLabel>
                            <Input type="number" {...field} />
                            <FormMessage />
                          </div>
                        )}
                      />
                    </div>
                  </div>
                  <FormField
                    //@ts-ignore
                    control={form.control}
                    name="snils"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="snils">
                          Снилс (формат 217-2117-217)
                        </FormLabel>
                        <Input {...field} />
                        <FormMessage />
                      </div>
                    )}
                  />
                  <FormField
                    //@ts-ignore
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <Label htmlFor="password">Пароль</Label>
                        <Input {...field} type="password" />
                        <FormMessage />
                      </div>
                    )}
                  />
                  <Button type="submit" className="w-full" disabled={isPending}>
                    Зарегистрироваться
                  </Button>
                </div>
                <div className="text-center text-sm">
                  Есть аккаунт?{" "}
                  <Link
                    to={routes.signIn}
                    className="underline underline-offset-4"
                  >
                    Войти
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
