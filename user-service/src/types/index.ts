export enum UserRole {
  MANAGER = 'MANAGER',
  MEMBER = 'MEMBER'
}

export interface UserAttributes {
  id: string;
  username: string;
  email: string;
  password_hash: string;
  role: UserRole;
  created_at?: Date;
  updated_at?: Date;
  deleted_at?: Date | null;
}

export interface UserCreationAttributes extends Omit<UserAttributes, 'id' | 'created_at' | 'updated_at'> {
  id?: string;
}

export interface AuthContext {
  user: any | null;
  isAuthenticated: boolean;
  error?: string;
  req: any;
}

export interface JWTPayload {
  userId: string;
  iat?: number;
  exp?: number;
}

export interface RegisterInput {
  username: string;
  email: string;
  password: string;
  role?: UserRole;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface UpdateUserInput {
  username?: string;
  email?: string;
  role?: UserRole;
}

export interface ChangePasswordInput {
  currentPassword: string;
  newPassword: string;
}

export interface AuthPayload {
  token: string;
  user: any;
}

// Additional types for future Team integration
export interface TeamAttributes {
  id: string;
  name: string;
  description?: string;
  created_by_id: string;
  created_at?: Date;
  updated_at?: Date;
  deleted_at?: Date | null;
}

export interface TeamUserAttributes {
  id: string;
  team_id: string;
  user_id: string;
  joined_at?: Date;
}
