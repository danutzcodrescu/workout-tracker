import {
  timestamp,
  jsonb,
  bigint,
  pgSchema,
  doublePrecision,
  smallserial,
  varchar,
  text,
  integer,
} from 'drizzle-orm/pg-core';

export const workoutsSchema = pgSchema('workouts');

export const activitiesTable = workoutsSchema.table('activities', {
  date: timestamp('date').primaryKey(),
  laps: jsonb('laps').notNull(),
  distance: bigint('distance', { mode: 'number' }).notNull(),
  duration: bigint('duration', { mode: 'number' }).notNull(),
  calories: bigint('calories', { mode: 'number' }).notNull(),
  pace: doublePrecision('pace').notNull(),
  groupId: integer('group_id').references(() => groupsTable.id),
});

export const groupsTable = workoutsSchema.table('groups', {
  id: smallserial('id').primaryKey().notNull(),
  name: varchar('name').notNull(),
  description: text('description'),
});
