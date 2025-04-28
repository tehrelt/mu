import z from "zod";

export const tokenPairSchema = z.object({
  accessToken: z.string().jwt(),
});

export const loginSchema = z.object({
  login: z.string().email(),
  password: z.string().min(8).max(100),
});

export const profileSchema = z.object({
  id: z.string().uuid(),
  lastName: z.string().min(2).max(100),
  firstName: z.string().min(2).max(100),
  middleName: z.string().min(2).max(100).optional(),
  email: z.string().email(),
});

export type Profile = z.infer<typeof profileSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type AccessToken = z.infer<typeof tokenPairSchema>;

export const registerSchema = z.object({
  lastName: z.string().min(2).max(100).trim(),
  firstName: z.string().min(2).max(100).trim(),
  middleName: z.string().min(2).max(100).trim().optional(),
  email: z.string().email().trim(),
  phone: z
    .string()
    .length(10)
    .trim()
    .transform((v) => `7${v}`),
  passport: z.object({
    series: z.preprocess(
      (s) => (typeof s === "string" ? +s : s),
      z.number().min(999).max(9999)
    ),
    number: z.preprocess(
      (s) => (typeof s === "string" ? +s : s),
      z.number().min(99999).max(999999)
    ),
  }),
  snils: z.string().length(12),
  password: z.string().min(8).max(100),
});
export type RegisterInput = z.infer<typeof registerSchema>;
