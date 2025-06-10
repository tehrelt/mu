import { Balance } from "@/components/ui/balance";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { datef } from "@/shared/lib/utils";
import { routes } from "@/shared/routes";
import { billingService } from "@/shared/services/billing.service";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Link, Navigate, useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";

export const ProcessPaymentPage = () => {
  const params = useParams();

  const id = params["id"];
  if (!id) {
    return <Navigate to={routes.billing.index} />;
  }

  return <ProcessPayment id={id} />;
};

const ProcessPayment = ({ id }: { id: string }) => {
  const paymentQuery = useQuery({
    queryKey: ["payment", { id }],
    queryFn: async () => await billingService.find(id),
  });

  const navigate = useNavigate();

  const goToDashboard = () => navigate(routes.dashboard.index);

  const { mutate: pay } = useMutation({
    mutationKey: ["payment", { id }, "pay"],
    mutationFn: async () => await billingService.process(id, "pay"),
    onSuccess: () => {
      toast.success("Платеж успешно отменен");
      goToDashboard();
    },
  });
  const { mutate: cancel } = useMutation({
    mutationKey: ["payment", { id }, "cancel"],
    mutationFn: async () => await billingService.process(id, "cancel"),
    onSuccess: () => {
      toast.success("Платеж успешно отменен");
      goToDashboard();
    },
  });

  if (paymentQuery.isLoading) {
    return <p>Loading...</p>;
  }

  if (!paymentQuery.data) {
    return (
      <p>
        При загрузке платежа произошла ошибка.{" "}
        <Link to={routes.dashboard.index}>Перейти на главную страницу</Link>
      </p>
    );
  }

  if (paymentQuery.data.status != "pending") {
    return (
      <p>
        Платеж уже обработан.{" "}
        <Link to={routes.dashboard.index}>Перейти на главную страницу</Link>
      </p>
    );
  }

  const payment = paymentQuery.data;

  return (
    <div className="min-h-screen flex justify-center items-center">
      <Card className="w-[480px]">
        <CardHeader>
          <CardTitle>Оплата платежа</CardTitle>
          <CardDescription className="flex justify-between">
            <div>id: {payment.id}</div>
            <div>{datef(payment.createdAt)}</div>
          </CardDescription>
        </CardHeader>
        <CardContent className="grid">
          <Separator className="w-full" />
          <div className="flex items-center gap-x-2">
            К оплате: <Balance balance={payment.amount} />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Button className="w-full" onClick={() => pay()}>
            Оплатить
          </Button>
          <Button
            className="w-full"
            variant={"outline"}
            onClick={() => cancel()}
          >
            Отмена
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};
