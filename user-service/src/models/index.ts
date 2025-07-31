import { Sequelize } from 'sequelize';
import config from '@/config/database';
import { UserFactory } from './user';

const env = process.env.NODE_ENV || 'development';
const dbConfig = config[env as keyof typeof config];

const sequelize = new Sequelize(
  dbConfig.database,
  dbConfig.username,
  dbConfig.password,
  dbConfig
);

const db = {
  sequelize,
  Sequelize,
  User: UserFactory(sequelize)
};

// Define associations
Object.values(db).forEach((model: any) => {
  if (model.associate) {
    model.associate(db);
  }
});

export { sequelize };
export default db;
