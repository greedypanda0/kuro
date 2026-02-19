/*
  Warnings:

  - You are about to drop the `Repositories` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE "Repositories";

-- CreateTable
CREATE TABLE "repositories" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "repositories_pkey" PRIMARY KEY ("id")
);
