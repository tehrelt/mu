"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, loginSchema } from "@/schemes/login";
import { Form, FormField, FormLabel, FormMessage } from "./ui/form";
import { login } from "@/actions/auth";
import { useFormStatus } from "react-dom";
import { signIn } from "next-auth/react";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { pending } = useFormStatus();

  const form = useForm<LoginSchema>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      login: "",
      password: "",
    },
  });

  const onSubmit = async (data: LoginSchema) => {
    try {
      await login(data);
    } catch (e) {
      console.log(e);
      throw new Error("Login failed");
    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Добро пожаловать</CardTitle>
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit, console.warn)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <FormField
                    control={form.control}
                    name="login"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="email">Email</FormLabel>
                        <Input
                          id="email"
                          type="email"
                          placeholder="m@example.com"
                          {...field}
                        />
                        <FormMessage />
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
                          <a
                            href="#"
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                          >
                            Забыли пароль?
                          </a>
                        </div>
                        <Input id="password" type="password" {...field} />
                      </div>
                    )}
                  />
                  <Button type="submit" className="w-full" disabled={pending}>
                    {pending ? "Входим..." : "Вход"}
                  </Button>
                </div>
                <div className="text-center text-sm">
                  Нет аккаунта?{" "}
                  <a href="#" className="underline underline-offset-4">
                    Зарегистрироваться
                  </a>
                </div>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
      <div className="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4">
        ПИ-21а Евтеев Дмитрий
      </div>
    </div>
  );
}
