import z from "zod";

export const houseAccountSchema = z.object({
  id: z.string().uuid(),
  houseId: z.string().uuid(),
  userId: z.string().uuid(),
  balance: z.number().min(0),
});

export type HouseAccount = z.infer<typeof houseAccountSchema>;
