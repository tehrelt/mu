import { z } from "zod";

export const loginRequestSchema = z.object({
  login: z.string().min(2).max(100).trim(),
  password: z
    .string()
    .min(8, { message: "Пароль должен содержать не менее 8 символов" })
    .max(100),
});

export const loginResponseSchema = z.object({
  accessToken: z.string().jwt(),
  refreshToken: z.string().jwt(),
});

export type LoginRequestSchema = z.infer<typeof loginRequestSchema>;
export type LoginResponseSchema = z.infer<typeof loginResponseSchema>;
export type TokenPair = z.infer<typeof loginResponseSchema>;
