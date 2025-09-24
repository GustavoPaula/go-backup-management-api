CREATE TABLE "devices" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "customer_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);

-- Foreign Key
ALTER TABLE "devices" ADD CONSTRAINT "devices_customer_id_fkey" 
FOREIGN KEY ("customer_id") REFERENCES "customers"("id") 
ON DELETE RESTRICT ON UPDATE CASCADE;

-- Índices para performance
CREATE INDEX "idx_devices_customer_id" ON "devices"("customer_id");