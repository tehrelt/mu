import { localizeTicketStatus, TicketStatusEnum } from "@/shared/types/ticket";
import { Badge } from "./ui/badge";

type Props = {
  val: TicketStatusEnum;
};

export const TicketStatus = ({ val }: Props) => {
  const variant =
    val === "TicketStatusApproved"
      ? "success"
      : val === "TicketStatusRejected"
      ? "destructive"
      : "pending";

  return <Badge variant={variant}>{localizeTicketStatus(val)}</Badge>;
};
