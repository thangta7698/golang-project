import { Request } from 'express';
import { verifyToken, extractTokenFromHeader } from '@/utils/jwt';
import db from '@/models';
import { AuthContext, UserRole } from '@/types';
import { User } from '@/models/user';

const { User: UserModel } = db;

export const authenticate = async (req: Request): Promise<AuthContext> => {
  const token = extractTokenFromHeader(req.headers.authorization);

  if (!token) {
    return { user: null, isAuthenticated: false, req };
  }

  try {
    const decoded = verifyToken(token);
    const user = await UserModel.findByPk(decoded.userId);

    if (!user) {
      return { user: null, isAuthenticated: false, req };
    }

    return { user, isAuthenticated: true, req };
  } catch (error) {
    return {
      user: null,
      isAuthenticated: false,
      error: error instanceof Error ? error.message : 'Authentication failed',
      req
    };
  }
};

export const requireAuth = (context: AuthContext): User => {
  if (!context.isAuthenticated || !context.user) {
    throw new Error('Authentication required');
  }
  return context.user;
};

export const requireRole = (context: AuthContext, requiredRole: UserRole): User => {
  const user = requireAuth(context);

  if (user.role !== requiredRole) {
    throw new Error(`Access denied. Required role: ${requiredRole}`);
  }

  return user;
};
