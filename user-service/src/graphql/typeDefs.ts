import { gql } from 'apollo-server-express';

const typeDefs = gql`
  scalar DateTime
  scalar UUID

  enum UserRole {
    MANAGER
    MEMBER
  }

  type User {
    id: UUID!
    username: String!
    email: String!
    role: UserRole!
    createdAt: DateTime!
    updatedAt: DateTime!
  }

  type AuthPayload {
    token: String!
    user: User!
  }

  input RegisterInput {
    username: String!
    email: String!
    password: String!
    role: UserRole = MEMBER
  }

  input LoginInput {
    email: String!
    password: String!
  }

  input UpdateUserInput {
    username: String
    email: String
    role: UserRole
  }

  input ChangePasswordInput {
    currentPassword: String!
    newPassword: String!
  }

  type Query {
    me: User
    users: [User!]!
    user(id: UUID!): User
  }

  type Mutation {
    register(input: RegisterInput!): AuthPayload!
    login(input: LoginInput!): AuthPayload!
    updateProfile(input: UpdateUserInput!): User!
    changePassword(input: ChangePasswordInput!): Boolean!
    deleteUser(id: UUID!): Boolean!
  }
`;

export default typeDefs;
