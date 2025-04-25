import z from "zod";

export const houseAccountSchema = z.object({
  id: z.string().uuid(),
  userId: z.string().uuid(),
  house: z.object({
    id: z.string().uuid(),
    address: z.string(),
  }),
  balance: z.number().min(0),
});

export type HouseAccount = z.infer<typeof houseAccountSchema>;

export const paymentStatusSchema = z.enum(["pending", "success", "canceled"]);
export type PaymentStatus = z.infer<typeof paymentStatusSchema>;

export const paymentSchema = z.object({
  id: z.string().uuid(),
  amount: z.number(),
  status: paymentStatusSchema,
  createdAt: z.preprocess((arg) => {
    if (typeof arg === "string" || arg instanceof Date) return new Date(arg);
  }, z.date()),
  updatedAt: z
    .preprocess((arg) => {
      if (typeof arg === "string" || arg instanceof Date) return new Date(arg);
    }, z.date())
    .optional(),
});
export type Payment = z.infer<typeof paymentSchema>;
