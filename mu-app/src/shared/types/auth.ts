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
