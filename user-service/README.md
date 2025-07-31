# User Service

Node.js service with Express, GraphQL, Sequelize ORM, and JWT authentication.

## Features

- **GraphQL API** with Apollo Server
- **JWT Authentication** with role-based access control
- **PostgreSQL** database with Sequelize ORM
- **User Management** (register, login, profile management)
- **Role-based permissions** (MANAGER, MEMBER)
- **Security** with Helmet and CORS
- **Environment configuration** with dotenv

## Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Environment configuration:**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Database setup:**
   ```bash
   # Make sure PostgreSQL is running
   npm run migrate
   ```

4. **Start development server:**
   ```bash
   npm run dev
   ```

## GraphQL Playground

Visit `http://localhost:4000/graphql` to access the GraphQL Playground.

## API Examples

### Register User
```graphql
mutation {
  register(input: {
    username: "johndoe"
    email: "john@example.com"
    password: "password123"
    role: MEMBER
  }) {
    token
    user {
      id
      username
      email
      role
    }
  }
}
```

### Login
```graphql
mutation {
  login(input: {
    email: "john@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      username
      email
      role
    }
  }
}
```

### Get Current User (requires authentication)
```graphql
query {
  me {
    id
    username
    email
    role
    createdAt
  }
}
```

## Authentication

Include JWT token in request headers:
```
Authorization: Bearer <your-jwt-token>
```

## Environment Variables

- `PORT`: Server port (default: 4000)
- `NODE_ENV`: Environment (development/production)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_NAME`: Database name
- `DB_USERNAME`: Database username
- `DB_PASSWORD`: Database password
- `JWT_SECRET`: JWT secret key
- `JWT_EXPIRES_IN`: JWT expiration time
- `CORS_ORIGIN`: CORS allowed origin
