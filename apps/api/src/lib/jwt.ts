import jwt from 'jsonwebtoken';
import { User } from '@labhaus/types';

const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key-change-in-production';
const JWT_REFRESH_SECRET =
  process.env.JWT_REFRESH_SECRET || 'your-refresh-secret-key-change-in-production';
const JWT_EXPIRES_IN = process.env.JWT_EXPIRES_IN || '1h';
const JWT_REFRESH_EXPIRES_IN = process.env.JWT_REFRESH_EXPIRES_IN || '7d';

export interface JwtPayload {
  user_id: string;
  email: string;
}

/**
 * Generate access token
 */
export function generateAccessToken(user: Pick<User, 'id' | 'email'>): string {
  const payload: JwtPayload = {
    user_id: user.id,
    email: user.email,
  };

  return jwt.sign(payload, JWT_SECRET, {
    expiresIn: JWT_EXPIRES_IN,
  } as jwt.SignOptions);
}

/**
 * Generate refresh token
 */
export function generateRefreshToken(user: Pick<User, 'id' | 'email'>): string {
  const payload: JwtPayload = {
    user_id: user.id,
    email: user.email,
  };

  return jwt.sign(payload, JWT_REFRESH_SECRET, {
    expiresIn: JWT_REFRESH_EXPIRES_IN,
  } as jwt.SignOptions);
}

/**
 * Verify access token
 */
export function verifyAccessToken(token: string): JwtPayload {
  try {
    return jwt.verify(token, JWT_SECRET) as JwtPayload;
  } catch (error) {
    throw new Error('Invalid or expired token');
  }
}

/**
 * Verify refresh token
 */
export function verifyRefreshToken(token: string): JwtPayload {
  try {
    return jwt.verify(token, JWT_REFRESH_SECRET) as JwtPayload;
  } catch (error) {
    throw new Error('Invalid or expired refresh token');
  }
}

/**
 * Parse expires_in to seconds
 */
export function getExpiresInSeconds(expiresIn: string): number {
  if (expiresIn.endsWith('h')) {
    return parseInt(expiresIn) * 3600;
  } else if (expiresIn.endsWith('d')) {
    return parseInt(expiresIn) * 86400;
  } else if (expiresIn.endsWith('m')) {
    return parseInt(expiresIn) * 60;
  }
  return parseInt(expiresIn);
}

export const TOKEN_EXPIRES_IN = getExpiresInSeconds(JWT_EXPIRES_IN);
