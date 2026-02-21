/*
  Warnings:

  - A unique constraint covering the columns `[user_id,name]` on the table `repositories` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "repositories_user_id_name_key" ON "repositories"("user_id", "name");
