import z from "zod";

export const paymentStatusSchema = z.enum(["pending", "success", "canceled"]);
export type PaymentStatus = z.infer<typeof paymentStatusSchema>;

export const paymentSchema = z.object({
  id: z.string().uuid(),
  amount: z.number(),
  status: paymentStatusSchema,
  message: z.string(),
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
