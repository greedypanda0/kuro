/*
  Warnings:

  - You are about to drop the column `author` on the `repositories` table. All the data in the column will be lost.
  - You are about to drop the column `description` on the `repositories` table. All the data in the column will be lost.
  - Added the required column `user_id` to the `repositories` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "repositories" DROP COLUMN "author",
DROP COLUMN "description",
ADD COLUMN     "user_id" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "repositories" ADD CONSTRAINT "repositories_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
