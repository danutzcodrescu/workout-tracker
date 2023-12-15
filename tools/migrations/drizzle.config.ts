import 'dotenv/config';
import type { Config } from 'drizzle-kit';
import path from 'path';

export default {
  schema: path.resolve(__dirname, './schema.ts'),
  out: path.resolve(__dirname, './sql/'),
  driver: 'pg',
  dbCredentials: {
    host: process.env.DB_HOST as string,
    user: process.env.DB_USER as string,
    password: process.env.DB_PASSWORD as string,
    database: process.env.DB_NAME as string,
    ssl: true
  },
  verbose: true,
  strict: true,
  schemaFilter: 'workouts'
} satisfies Config;
