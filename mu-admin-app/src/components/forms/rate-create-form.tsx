import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Form, FormField, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { routes } from "@/shared/routes";
import { rateService } from "@/shared/services/rates.service";
import { RateCreate, rateCreateSchema } from "@/shared/types/rate";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";

export function RateCreateForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const form = useForm<RateCreate>({
    resolver: zodResolver(rateCreateSchema),
    defaultValues: {
      initialRate: 0,
      measureUnit: "",
      name: "",
    },
  });

  const navigate = useNavigate();

  const { mutate: create, isPending } = useMutation({
    mutationKey: ["rate-create"],
    mutationFn: async (data: RateCreate) => await rateService.create(data),
    onSuccess: () => {
      toast.success("Тариф успешно создан");
      navigate(routes.home);
    },
    onError: (e) => {
      console.log(e);
      toast.error("Не удалось создать тариф", {
        description: "Проверьте правильность введенных данных",
      });
    },
  });

  const submit = async (data: RateCreate) => {
    create(data);
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Создание тарифа</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(submit, console.warn)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="name">Название</FormLabel>
                        <Input {...field} />
                      </div>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="measureUnit"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="measureUnit">
                          Ед. измерения
                        </FormLabel>
                        <Input {...field} />
                      </div>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="initialRate"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="initialRate">
                          Стоимость за 1 ед.
                        </FormLabel>
                        <Input {...field} type="number" step="0.01" />
                        <FormMessage />
                      </div>
                    )}
                  />
                  <Button type="submit" className="w-full" disabled={isPending}>
                    Создать
                  </Button>
                </div>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
