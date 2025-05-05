import { z } from "zod";

export const userSnippetSchema = z.object({
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

export type UserSnippet = z.infer<typeof userSnippetSchema>;

export const userSchema = userSnippetSchema.extend({
  passportSeries: z.number(),
  passportNumber: z.number(),
  snils: z.string(),
});

export type User = z.infer<typeof userSchema>;

export const userListSchema = z.object({
  users: z.array(userSnippetSchema),
});
export type UserList = z.infer<typeof userListSchema>;

export const houseSchema = z.object({
  id: z.string().uuid(),
  address: z.string(),
});

export const accountInfoSchema = z.object({
  id: z.string().uuid(),
  userId: z.string().uuid(),
  house: houseSchema,
  balance: z.number().multipleOf(0.01),
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
export type AccountInfo = z.infer<typeof accountInfoSchema>;

export const userAccountsSchema = z.object({
  accounts: z.array(accountInfoSchema),
});
export type UserAccounts = z.infer<typeof userAccountsSchema>;
