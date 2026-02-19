/*
  Warnings:

  - Added the required column `expires_at` to the `auth_tokens` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "auth_tokens" ADD COLUMN     "expires_at" TIMESTAMP(3) NOT NULL;
