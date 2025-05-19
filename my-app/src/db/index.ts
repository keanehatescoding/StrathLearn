import { drizzle } from 'drizzle-orm/postgres-js';
import 'dotenv/config';
import postgres from 'postgres';
import * as schema from './schema';

const client = postgres(process.env.DATABASE_URL as string);
const db = drizzle(client, { schema });

export default db;