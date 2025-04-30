import z from "zod";

export const cabinetSchema = z.object({
  id: z.string().uuid(),
  accountId: z.string().uuid(),
  serviceId: z.string().uuid(),
  consumed: z.number().default(0),
  createdAt: z.preprocess((arg) => {
    if (typeof arg === "string" || arg instanceof Date) return new Date(arg);
  }, z.date()),
  updatedAt: z
    .preprocess((arg) => {
      if (typeof arg === "string" || arg instanceof Date) return new Date(arg);
    }, z.date())
    .optional(),
});

export type Cabinet = z.infer<typeof cabinetSchema>;
