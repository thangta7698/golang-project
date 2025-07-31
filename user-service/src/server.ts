import 'dotenv/config';
import express from 'express';
import { ApolloServer } from 'apollo-server-express';
import cors from 'cors';
import helmet from 'helmet';
import { GraphQLScalarType, Kind } from 'graphql';

import db from '@/models';
import typeDefs from '@/graphql/typeDefs';
import resolvers from '@/graphql/resolvers';
import { authenticate } from '@/middleware/auth';

// Custom scalar resolvers
const scalarResolvers = {
  DateTime: new GraphQLScalarType({
    name: 'DateTime',
    description: 'Date custom scalar type',
    serialize(value: any) {
      return value instanceof Date ? value.toISOString() : value;
    },
    parseValue(value: any) {
      return new Date(value);
    },
    parseLiteral(ast) {
      if (ast.kind === Kind.STRING) {
        return new Date(ast.value);
      }
      return null;
    },
  }),
  UUID: new GraphQLScalarType({
    name: 'UUID',
    description: 'UUID custom scalar type',
    serialize(value: any) {
      return value;
    },
    parseValue(value: any) {
      return value;
    },
    parseLiteral(ast) {
      if (ast.kind === Kind.STRING) {
        return ast.value;
      }
      return null;
    },
  }),
};

async function startServer(): Promise<void> {
  const app = express();

  // Security middleware
  app.use(helmet({
    contentSecurityPolicy: false,
    crossOriginEmbedderPolicy: false
  }));

  // CORS configuration
  app.use(cors({
    origin: process.env.CORS_ORIGIN || 'http://localhost:3000',
    credentials: true
  }));

  // Body parsing middleware
  app.use(express.json());
  app.use(express.urlencoded({ extended: true }));

  // Test database connection
  try {
    await db.sequelize.authenticate();
    console.log('Database connection established successfully.');

    // Sync database (use { force: true } to recreate tables)
    await db.sequelize.sync({ alter: true });
    console.log('Database synchronized successfully.');
  } catch (error) {
    console.error('Unable to connect to the database:', error);
    process.exit(1);
  }

  // Create Apollo Server
  const server = new ApolloServer({
    typeDefs,
    resolvers: {
      ...resolvers,
      ...scalarResolvers
    },
    context: async ({ req }) => {
      const auth = await authenticate(req);
      return auth;
    },
    introspection: process.env.NODE_ENV !== 'production',
    playground: process.env.NODE_ENV !== 'production'
  });

  await server.start();

  // Apply Apollo GraphQL middleware
  server.applyMiddleware({
    app,
    path: '/graphql',
    cors: false // We handle CORS above
  });

  // Health check endpoint
  app.get('/health', (req, res) => {
    res.json({
      status: 'OK',
      timestamp: new Date().toISOString(),
      service: 'user-service'
    });
  });

  // 404 handler
  app.use('*', (req, res) => {
    res.status(404).json({
      error: 'Not Found',
      message: `Route ${req.originalUrl} not found`
    });
  });

  // Error handler
  app.use((error: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
    console.error('Server Error:', error);
    res.status(500).json({
      error: 'Internal Server Error',
      message: process.env.NODE_ENV === 'production' ? 'Something went wrong' : error.message
    });
  });

  const PORT = process.env.PORT || 4000;

  app.listen(PORT, () => {
    console.log(`🚀 User Service running on http://localhost:${PORT}`);
    console.log(`🎯 GraphQL endpoint: http://localhost:${PORT}${server.graphqlPath}`);
    console.log(`📊 GraphQL Playground: http://localhost:${PORT}${server.graphqlPath}`);
  });
}

// Handle unhandled promise rejections
process.on('unhandledRejection', (err: Error) => {
  console.error('Unhandled Promise Rejection:', err);
  process.exit(1);
});

// Handle uncaught exceptions
process.on('uncaughtException', (err: Error) => {
  console.error('Uncaught Exception:', err);
  process.exit(1);
});

// Graceful shutdown
process.on('SIGTERM', async () => {
  console.log('SIGTERM received, shutting down gracefully');
  await db.sequelize.close();
  process.exit(0);
});

startServer().catch(error => {
  console.error('Failed to start server:', error);
  process.exit(1);
});
