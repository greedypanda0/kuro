import { auth } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import crypto from "crypto";

export async function POST(request: Request) {
  const session = await auth();

  if (!session?.user?.id) {
    return Response.json({ error: "Unauthorized" }, { status: 401 });
  }

  const body = await request.json();
  const expiresAt = new Date(body.expiresAt);

  if (isNaN(expiresAt.getTime())) {
    return Response.json({ error: "Invalid expiration date" }, { status: 400 });
  }

  const rawToken = crypto.randomBytes(32).toString("hex");
  const hashedToken = crypto
    .createHash("sha256")
    .update(rawToken)
    .digest("hex");

  await prisma.authTokens.create({
    data: {
      id: hashedToken,
      userId: session.user.id,
      expiresAt,
    },
  });

  return Response.json({ token: rawToken }, { status: 201 });
}
