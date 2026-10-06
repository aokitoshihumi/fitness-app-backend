CREATE TABLE `exercises` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`name` text NOT NULL,
	`body_part` text NOT NULL
);
--> statement-breakpoint
CREATE TABLE `workout_exercises` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`session_id` integer NOT NULL,
	`exercise_id` integer NOT NULL,
	`position` integer NOT NULL,
	FOREIGN KEY (`session_id`) REFERENCES `workout_sessions`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`exercise_id`) REFERENCES `exercises`(`id`) ON UPDATE no action ON DELETE no action
);
--> statement-breakpoint
CREATE INDEX `workout_exercises_session_idx` ON `workout_exercises` (`session_id`);--> statement-breakpoint
CREATE TABLE `workout_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`started_at` text NOT NULL,
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE INDEX `workout_sessions_user_started_idx` ON `workout_sessions` (`user_id`,`started_at`);--> statement-breakpoint
CREATE TABLE `workout_sets` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`workout_exercise_id` integer NOT NULL,
	`position` integer NOT NULL,
	`weight_kg` real NOT NULL,
	`reps` integer NOT NULL,
	FOREIGN KEY (`workout_exercise_id`) REFERENCES `workout_exercises`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE INDEX `workout_sets_workout_exercise_idx` ON `workout_sets` (`workout_exercise_id`);--> statement-breakpoint
ALTER TABLE `users` ADD `timezone` text DEFAULT 'Asia/Tokyo' NOT NULL;