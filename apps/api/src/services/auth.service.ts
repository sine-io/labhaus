import bcrypt from 'bcrypt';
import { User, AuthTokens, RegisterRequest, LoginRequest } from '@labhaus/types';
import { query } from '../db.js';
import {
  generateAccessToken,
  generateRefreshToken,
  verifyRefreshToken,
  TOKEN_EXPIRES_IN,
} from '../lib/jwt.js';

const SALT_ROUNDS = 10;

export class AuthService {
  /**
   * Register a new user
   */
  async register(data: RegisterRequest): Promise<{ user: User; tokens: AuthTokens }> {
    // Check if user exists
    const existing = await query('SELECT id FROM users WHERE email = $1', [data.email]);
    if (existing.rows.length > 0) {
      throw new Error('User already exists');
    }

    // Hash password
    const password_hash = await bcrypt.hash(data.password, SALT_ROUNDS);

    // Create user
    const result = await query(
      `INSERT INTO users (email, password_hash, name)
       VALUES ($1, $2, $3)
       RETURNING *`,
      [data.email, password_hash, data.name || null]
    );

    const user = result.rows[0] as User;

    // Generate tokens
    const tokens = await this.generateTokens(user);

    return { user, tokens };
  }

  /**
   * Login with email and password
   */
  async login(data: LoginRequest): Promise<{ user: User; tokens: AuthTokens }> {
    // Find user
    const result = await query('SELECT * FROM users WHERE email = $1', [data.email]);
    if (result.rows.length === 0) {
      throw new Error('Invalid credentials');
    }

    const user = result.rows[0] as User;

    // Check password
    if (!user.password_hash) {
      throw new Error('Invalid credentials');
    }

    const valid = await bcrypt.compare(data.password, user.password_hash);
    if (!valid) {
      throw new Error('Invalid credentials');
    }

    // Generate tokens
    const tokens = await this.generateTokens(user);

    return { user, tokens };
  }

  /**
   * Refresh access token
   */
  async refreshToken(refreshToken: string): Promise<AuthTokens> {
    // Verify refresh token
    const payload = verifyRefreshToken(refreshToken);

    // Check if refresh token exists in database
    const result = await query(
      'SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()',
      [refreshToken]
    );

    if (result.rows.length === 0) {
      throw new Error('Invalid refresh token');
    }

    // Get user
    const userResult = await query('SELECT * FROM users WHERE id = $1', [payload.user_id]);
    if (userResult.rows.length === 0) {
      throw new Error('User not found');
    }

    const user = userResult.rows[0] as User;

    // Generate new tokens
    return this.generateTokens(user);
  }

  /**
   * Get user by ID
   */
  async getUserById(userId: string): Promise<User | null> {
    const result = await query('SELECT * FROM users WHERE id = $1', [userId]);
    return result.rows[0] || null;
  }

  /**
   * Google OAuth login/register
   */
  async googleAuth(googleId: string, email: string, name?: string, avatarUrl?: string) {
    // Check if user exists with google_id
    let result = await query('SELECT * FROM users WHERE google_id = $1', [googleId]);

    let user: User;

    if (result.rows.length > 0) {
      // Existing user
      user = result.rows[0] as User;
    } else {
      // Check if email exists
      result = await query('SELECT * FROM users WHERE email = $1', [email]);

      if (result.rows.length > 0) {
        // Link Google account to existing email user
        result = await query(
          'UPDATE users SET google_id = $1, email_verified = true, updated_at = NOW() WHERE email = $2 RETURNING *',
          [googleId, email]
        );
        user = result.rows[0] as User;
      } else {
        // Create new user
        result = await query(
          `INSERT INTO users (email, google_id, name, avatar_url, email_verified)
           VALUES ($1, $2, $3, $4, true)
           RETURNING *`,
          [email, googleId, name || null, avatarUrl || null]
        );
        user = result.rows[0] as User;
      }
    }

    // Generate tokens
    const tokens = await this.generateTokens(user);

    return { user, tokens };
  }

  /**
   * Generate access and refresh tokens
   */
  private async generateTokens(user: User): Promise<AuthTokens> {
    const access_token = generateAccessToken(user);
    const refresh_token = generateRefreshToken(user);

    // Store refresh token in database
    await query(
      `INSERT INTO refresh_tokens (user_id, token, expires_at)
       VALUES ($1, $2, NOW() + INTERVAL '7 days')`,
      [user.id, refresh_token]
    );

    return {
      access_token,
      refresh_token,
      token_type: 'Bearer',
      expires_in: TOKEN_EXPIRES_IN,
    };
  }

  /**
   * Remove user's password_hash from response
   */
  sanitizeUser(user: User): Omit<User, 'password_hash'> {
    const { password_hash, ...sanitized } = user;
    return sanitized;
  }
}
