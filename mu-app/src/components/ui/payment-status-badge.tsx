import { PaymentStatus } from "@/shared/types/payment";
import { Badge } from "./badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./tooltip";

export const PaymentStatusBadge = ({ status }: { status: PaymentStatus }) => {
  const variant =
    status === "pending"
      ? "pending"
      : status === "success"
      ? "success"
      : "destructive";

  const tooltip =
    status === "pending"
      ? "Ожидание платежа"
      : status === "success"
      ? "Платёж принят"
      : "Платёж отклонен";

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>
          <Badge className="p-1" variant={variant} />
        </TooltipTrigger>
        <TooltipContent>{tooltip}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};
