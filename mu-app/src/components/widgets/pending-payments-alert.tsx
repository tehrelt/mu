import { routes } from "@/shared/routes";
import { AlertCircle } from "lucide-react";
import { Link } from "react-router-dom";
import { Alert, AlertDescription, AlertTitle } from "../ui/alert";
import { Button } from "../ui/button";

export const PendingPaymentsAlert = () => {
  return (
    <Alert variant="warn">
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>У вас есть неоплаченные платежи</AlertTitle>
      <AlertDescription>
        <Link to={routes.dashboard.addFunds}>
          <Button variant={"outline"}>Перейти к оплате</Button>
        </Link>
      </AlertDescription>
    </Alert>
  );
};
