import z from "zod";

export const ticketTypeEnum = z.enum(["account", "connect-service"]);
export type TicketType = z.infer<typeof ticketTypeEnum>;

export const ticketNewAccountInputSchema = z.object({
  address: z.string(),
});

export const ticketConnectServiceInputSchema = z.object({
  serviceId: z.string().uuid(),
  accountId: z.string().uuid(),
});

export const ticketNewPayload = z.union([
  ticketNewAccountInputSchema,
  ticketConnectServiceInputSchema,
]);
export const ticketNewInput = z.object({
  type: ticketTypeEnum,
  payload: ticketNewPayload,
});
export type NewTicketPayload = z.infer<typeof ticketNewPayload>;
export type NewTicketInput = z.infer<typeof ticketNewInput>;
export type TicketNewAccountInput = z.infer<typeof ticketNewAccountInputSchema>;
export type TicketConnectServiceInput = z.infer<
  typeof ticketConnectServiceInputSchema
>;
