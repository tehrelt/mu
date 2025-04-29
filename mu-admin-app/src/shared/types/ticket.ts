import z from "zod";

export const ticketTypeEnum = z.enum([
  "TicketTypeAccount",
  "TicketTypeConnectService",
]);
export type TicketTypeEnum = z.infer<typeof ticketTypeEnum>;
export const localizeTicketType = (t: TicketTypeEnum): string => {
  switch (t) {
    case "TicketTypeAccount":
      return "Новый счет";
    case "TicketTypeConnectService":
      return "Подключение услуги";
    default:
      return "";
  }
};

export const ticketStatusEnum = z.enum([
  "TicketStatusPending",
  "TicketStatusRejected",
  "TicketStatusApproved",
]);
export type TicketStatusEnum = z.infer<typeof ticketStatusEnum>;
export const localizeTicketStatus = (t: TicketStatusEnum): string => {
  switch (t) {
    case "TicketStatusPending":
      return "В ожидании";
    case "TicketStatusRejected":
      return "Отклонено";
    case "TicketStatusApproved":
      return "Одобрено";
    default:
      return "";
  }
};

export const ticketHeaderSchema = z.object({
  id: z.string(),
  ticket_type: ticketTypeEnum,
  ticket_status: ticketStatusEnum,
});
export type TicketHeader = z.infer<typeof ticketHeaderSchema>;
