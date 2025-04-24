import { z } from "zod";

export const userSchema = z.object({
  id: z.string().uuid(),
  lastName: z.string(),
  firstName: z.string(),
  middleName: z.string(),
  email: z.string().email(),
  phone: z.string(),
  createdAt: z.preprocess(
    (arg) => (typeof arg === "string" ? new Date(arg) : undefined),
    z.date()
  ),
  updatedAt: z
    .preprocess(
      (arg) => (typeof arg === "string" ? new Date(arg) : undefined),
      z.date()
    )
    .optional(),
});

export type User = z.infer<typeof userSchema>;

export const userListSchema = z.object({
  users: z.array(userSchema),
});
export type UserList = z.infer<typeof userListSchema>;
