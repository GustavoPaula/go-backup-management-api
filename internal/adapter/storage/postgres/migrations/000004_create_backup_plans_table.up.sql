-- CreateEnum
CREATE TYPE "week_day_enum" AS ENUM ('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');

-- CreateTable
CREATE TABLE "backup_plans" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "backup_size_bytes" BIGINT NOT NULL,
    "device_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- CreateTable
CREATE TABLE "backup_plans_week_day" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "day" "week_day_enum" NOT NULL,
    "time_day" TIME NOT NULL,
    "backup_plan_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- AddForeignKey
ALTER TABLE "backup_plans" ADD CONSTRAINT "backup_plan_device_id_fkey"
FOREIGN KEY ("device_id") REFERENCES "devices"("id")
ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "backup_plans_week_day" ADD CONSTRAINT "backup_plan_week_day_backup_plan_id_fkey"
FOREIGN KEY ("backup_plan_id") REFERENCES "backup_plans"("id") 
ON DELETE RESTRICT ON UPDATE CASCADE;