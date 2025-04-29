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
  ticketType: ticketTypeEnum,
  ticketStatus: ticketStatusEnum,
  createdBy: z.string(),
});
export type TicketHeader = z.infer<typeof ticketHeaderSchema>;

export const ticketNewAccount = ticketHeaderSchema.extend({
  ticketType: z.literal("TicketTypeAccount"),
  address: z.string(),
});
export type TicketNewAccount = z.infer<typeof ticketNewAccount>;

export const ticketConnectService = ticketHeaderSchema.extend({
  ticketType: z.literal("TicketTypeConnectService"),
  serviceId: z.string(),
  accountId: z.string(),
});
export type TicketConnectService = z.infer<typeof ticketConnectService>;

export type Ticket = TicketNewAccount | TicketConnectService;
