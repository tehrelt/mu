import { z } from "zod";

export const profileSchema = z.object({
  id: z.string().uuid(),
  lastName: z.string().min(2).max(100),
  firstName: z.string().min(2).max(100),
  middleName: z.string().min(2).max(100).optional(),
  email: z.string().email(),
});
export type Profile = z.infer<typeof profileSchema>;
