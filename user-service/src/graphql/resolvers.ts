import { Op } from 'sequelize';
import db from '@/models';
import { generateToken } from '@/utils/jwt';
import { requireAuth, requireRole } from '@/middleware/auth';
import {
  AuthContext,
  UserRole
} from '@/types';
import {
  Resolvers,
  MutationRegisterArgs,
  MutationLoginArgs,
  MutationUpdateProfileArgs,
  MutationChangePasswordArgs,
  MutationDeleteUserArgs,
  QueryUserArgs
} from '@/generated/graphql';

const { User } = db;

const resolvers: Resolvers = {
  Query: {
    me: async (_, __, context: AuthContext) => {
      return requireAuth(context);
    },

    users: async (_, __, context: AuthContext) => {
      requireRole(context, UserRole.MANAGER);
      return await User.findAll({
        order: [['created_at', 'DESC']]
      });
    },

    user: async (_, { id }: QueryUserArgs, context: AuthContext) => {
      requireAuth(context);
      const user = await User.findByPk(id);

      if (!user) {
        throw new Error('User not found');
      }

      // Only managers can view other users' details
      if (context.user.id !== id && context.user.role !== UserRole.MANAGER) {
        throw new Error('Access denied');
      }

      return user;
    }
  },

  Mutation: {
    register: async (_, { input }: MutationRegisterArgs) => {
      const { username, email, password, role } = input;

      // Check if user already exists
      const existingUser = await User.findOne({
        where: {
          [Op.or]: [{ email }, { username }]
        }
      });

      if (existingUser) {
        throw new Error('User with this email or username already exists');
      }

      // Validate password strength
      if (password.length < 6) {
        throw new Error('Password must be at least 6 characters long');
      }

      const userRole: any = role ?? UserRole.MEMBER;

      const user = await User.create({
        username,
        email,
        password_hash: password, // Will be hashed by the hook
        role: userRole
      });

      const token = generateToken({ userId: user.id });

      return {
        token,
        user
      };
    },

    login: async (_, { input }: MutationLoginArgs) => {
      const { email, password } = input;

      const user = await User.findOne({ where: { email } });

      if (!user) {
        throw new Error('Invalid email or password');
      }

      const isValidPassword = await user.checkPassword(password);

      if (!isValidPassword) {
        throw new Error('Invalid email or password');
      }

      const token = generateToken({ userId: user.id });

      return {
        token,
        user
      };
    },

    updateProfile: async (_, { input }: MutationUpdateProfileArgs, context: AuthContext) => {
      const user = requireAuth(context);

      // Check if username or email is already taken by another user
      if (input.username || input.email) {
        const whereClause: any = {};
        if (input.username) whereClause.username = input.username;
        if (input.email) whereClause.email = input.email;

        const existingUser = await User.findOne({
          where: {
            ...whereClause,
            id: { [Op.ne]: user.id }
          }
        });

        if (existingUser) {
          throw new Error('Username or email already taken');
        }
      }

      // Only managers can change roles
      if (input.role && user.role !== UserRole.MANAGER) {
        throw new Error('Access denied. Only managers can change roles');
      }

      await user.update(input);
      await user.reload();

      return user;
    },

    changePassword: async (_, { input }: MutationChangePasswordArgs, context: AuthContext) => {
      const user = requireAuth(context);
      const { currentPassword, newPassword } = input;

      const isValidPassword = await user.checkPassword(currentPassword);

      if (!isValidPassword) {
        throw new Error('Current password is incorrect');
      }

      if (newPassword.length < 6) {
        throw new Error('New password must be at least 6 characters long');
      }

      await user.update({ password_hash: newPassword });

      return true;
    },

    deleteUser: async (_, { id }: MutationDeleteUserArgs, context: AuthContext) => {
      const currentUser = requireRole(context, UserRole.MANAGER);

      if (currentUser.id === id) {
        throw new Error('You cannot delete your own account');
      }

      const user = await User.findByPk(id);

      if (!user) {
        throw new Error('User not found');
      }

      await user.destroy();

      return true;
    }
  }
};

export default resolvers;
