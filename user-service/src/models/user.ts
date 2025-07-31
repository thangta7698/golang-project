import {
  DataTypes,
  Model,
  Sequelize,
  CreationOptional,
  InferAttributes,
  InferCreationAttributes,
  Association,
  HasManyGetAssociationsMixin,
  HasManyAddAssociationMixin,
  HasManyCreateAssociationMixin
} from 'sequelize';
import { v4 as uuidv4 } from 'uuid';
import bcrypt from 'bcryptjs';
import { UserRole } from '@/types';

export class User extends Model<
  InferAttributes<User>,
  InferCreationAttributes<User>
> {
  declare id: CreationOptional<string>;
  declare username: string;
  declare email: string;
  declare password_hash: string;
  declare role: UserRole;
  declare created_at: CreationOptional<Date>;
  declare updated_at: CreationOptional<Date>;
  declare deleted_at: CreationOptional<Date | null>;

  // Association mixins - these will be properly typed when you add the related models
  declare getOwnedTeams: HasManyGetAssociationsMixin<any>;
  declare addOwnedTeam: HasManyAddAssociationMixin<any, string>;
  declare createOwnedTeam: HasManyCreateAssociationMixin<any>;

  declare getOwnedFolders: HasManyGetAssociationsMixin<any>;
  declare addOwnedFolder: HasManyAddAssociationMixin<any, string>;
  declare createOwnedFolder: HasManyCreateAssociationMixin<any>;

  declare getOwnedNotes: HasManyGetAssociationsMixin<any>;
  declare addOwnedNote: HasManyAddAssociationMixin<any, string>;
  declare createOwnedNote: HasManyCreateAssociationMixin<any>;

  // Associations - these will be defined when you add the related models
  declare static associations: {
    ownedTeams: Association<User, any>;
    ownedFolders: Association<User, any>;
    ownedNotes: Association<User, any>;
    folderShares: Association<User, any>;
    noteShares: Association<User, any>;
    teamUsers: Association<User, any>;
  };

  // Instance methods
  async checkPassword(password: string): Promise<boolean> {
    return await bcrypt.compare(password, this.password_hash);
  }

  toJSON(): object {
    const values = { ...this.get() };
    delete (values as any).password_hash;
    return values;
  }

  // Class method for associations
  static associate(models: any): void {
    // Define associations here when you add more models
    // User.hasMany(models.Team, { foreignKey: 'created_by_id', as: 'ownedTeams' });
    // User.hasMany(models.Folder, { foreignKey: 'owner_id', as: 'ownedFolders' });
    // User.hasMany(models.Note, { foreignKey: 'owner_id', as: 'ownedNotes' });
    // User.hasMany(models.FolderShare, { foreignKey: 'user_id', as: 'folderShares' });
    // User.hasMany(models.NoteShare, { foreignKey: 'user_id', as: 'noteShares' });
    // User.hasMany(models.TeamUser, { foreignKey: 'user_id', as: 'teamUsers' });
  }
}

export function UserFactory(sequelize: Sequelize): typeof User {
  User.init(
    {
      id: {
        type: DataTypes.UUID,
        primaryKey: true,
        defaultValue: () => uuidv4(),
        allowNull: false
      },
      username: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
        validate: {
          len: [3, 50],
          notEmpty: true
        }
      },
      email: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
        validate: {
          isEmail: true,
          notEmpty: true
        }
      },
      password_hash: {
        type: DataTypes.STRING,
        allowNull: false,
        field: 'password_hash'
      },
      role: {
        type: DataTypes.ENUM(...Object.values(UserRole)),
        allowNull: false,
        defaultValue: UserRole.MEMBER,
        validate: {
          isIn: [Object.values(UserRole)]
        }
      },
      created_at: {
        type: DataTypes.DATE,
        allowNull: false,
        defaultValue: DataTypes.NOW
      },
      updated_at: {
        type: DataTypes.DATE,
        allowNull: false,
        defaultValue: DataTypes.NOW
      },
      deleted_at: {
        type: DataTypes.DATE,
        allowNull: true
      }
    },
    {
      sequelize,
      modelName: 'User',
      tableName: 'users',
      underscored: true,
      paranoid: true,
      timestamps: true,
      hooks: {
        beforeCreate: async (user: User) => {
          if (user.password_hash) {
            user.password_hash = await bcrypt.hash(user.password_hash, 12);
          }
        },
        beforeUpdate: async (user: User) => {
          if (user.changed('password_hash')) {
            user.password_hash = await bcrypt.hash(user.password_hash, 12);
          }
        }
      }
    }
  );

  return User;
}
