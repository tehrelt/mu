import { z } from "zod";

export const loginSchema = z.object({
  login: z.string({ required_error: "Введите логин" }),
  password: z
    .string()
    .min(7, { message: "Пароль должен содержать не менее 8 символов" }),
});

export const tokenSchema = z.object({
  accessToken: z.string().jwt(),
  refreshToken: z.string().jwt(),
});

export type LoginSchema = z.infer<typeof loginSchema>;
export type TokenSchema = z.infer<typeof tokenSchema>;
