import {
  timestamp,
  jsonb,
  bigint,
  pgSchema,
  doublePrecision,
} from 'drizzle-orm/pg-core';

export const workoutsSchema = pgSchema('workouts');

export const activitiesTable = workoutsSchema.table('activities', {
  date: timestamp('date').primaryKey().notNull(),
  laps: jsonb('laps').notNull(),
  distance: bigint('distance', { mode: 'number' }).notNull(),
  duration: bigint('duration', { mode: 'number' }).notNull(),
  calories: bigint('calories', { mode: 'number' }).notNull(),
  pace: doublePrecision('pace').notNull(),
});
