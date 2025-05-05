import z from "zod";

export const houseAccountSchema = z.object({
  id: z.string().uuid(),
  userId: z.string().uuid(),
  house: z.object({
    id: z.string().uuid(),
    address: z.string(),
  }),
  balance: z.number(),
});

export type HouseAccount = z.infer<typeof houseAccountSchema>;
