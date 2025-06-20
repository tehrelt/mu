import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Form, FormField, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { routes } from "@/shared/routes";
import { rateService } from "@/shared/services/rates.service";
import {
  localizeServiceType,
  RateCreate,
  rateCreateSchema,
  serviceTypeSchema,
} from "@/shared/types/rate";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

export function RateCreateForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const form = useForm<RateCreate>({
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
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
            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-expect-error */}
            <form onSubmit={form.handleSubmit(submit, console.warn)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <FormField
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-expect-error
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
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-expect-error
                    control={form.control}
                    name="serviceType"
                    render={({ field }) => (
                      <div className="grid gap-3">
                        <FormLabel htmlFor="measureUnit">Вид услуги</FormLabel>
                        <Select onValueChange={(v) => field.onChange(v)}>
                          <SelectTrigger className="w-full">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {Object.entries(serviceTypeSchema.Enum).map(
                              ([_, t]) => (
                                <SelectItem value={t} key={t}>
                                  {localizeServiceType(t)}
                                </SelectItem>
                              )
                            )}
                          </SelectContent>
                        </Select>
                      </div>
                    )}
                  />
                  <FormField
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-expect-error
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
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-expect-error
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
