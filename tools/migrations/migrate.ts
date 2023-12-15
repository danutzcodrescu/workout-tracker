import { drizzle } from 'drizzle-orm/postgres-js';
import { migrate } from 'drizzle-orm/postgres-js/migrator';
import postgres from 'postgres';
import path from 'path';
import 'dotenv/config';

const sql = postgres(process.env.DB_URL as string, {
  password: process.env.DB_PASSWORD as string,
  max: 1,
});

const db = drizzle(sql);

migrate(db, { migrationsFolder: path.resolve(__dirname, './sql') })
  .then(() => sql.end())
  .then(() => console.log('done'));
