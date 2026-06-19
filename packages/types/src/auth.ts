import { z } from 'zod';

export const userSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  password_hash: z.string().nullable(),
  name: z.string().nullable(),
  avatar_url: z.string().url().nullable(),
  google_id: z.string().nullable(),
  email_verified: z.boolean(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

export const registerSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  name: z.string().min(1).optional(),
});

export const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(1),
});

export const refreshTokenSchema = z.object({
  refresh_token: z.string(),
});

export const googleAuthSchema = z.object({
  id_token: z.string(),
});

export type User = z.infer<typeof userSchema>;
export type RegisterRequest = z.infer<typeof registerSchema>;
export type LoginRequest = z.infer<typeof loginSchema>;
export type RefreshTokenRequest = z.infer<typeof refreshTokenSchema>;
export type GoogleAuthRequest = z.infer<typeof googleAuthSchema>;

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

export interface AuthResponse {
  user: Omit<User, 'password_hash'>;
  tokens: AuthTokens;
}
