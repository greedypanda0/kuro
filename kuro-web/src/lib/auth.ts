import NextAuth from "next-auth";
import GitHub from "next-auth/providers/github";
import { PrismaAdapter } from "@auth/prisma-adapter";
import { prisma } from "@/lib/prisma";
import { UAParser } from "ua-parser-js";
import { headers } from "next/headers";

export const { handlers, signIn, signOut, auth } = NextAuth({
  adapter: PrismaAdapter(prisma),
  providers: [
    GitHub({
      clientId: process.env.GITHUB_CLIENT_ID!,
      clientSecret: process.env.GITHUB_CLIENT_SECRET!,
      profile: (profile) => ({
        name: profile.name ?? profile.login,
        email: profile.email,
        image: profile.avatar_url,
        username: profile.login,
      }),
    }),
  ],
  callbacks: {
    session: async ({ session, user }) => {
      const dbSession = await prisma.session.findUnique({
        where: { sessionToken: session.sessionToken },
      });
      if (!dbSession) throw new Error("Session not found");

      if (!dbSession.deviceName) {
        const { os, browser } = await detectDeviceFromHeaders();
        await prisma.session.update({
          where: { id: dbSession.id },
          data: {
            userAgent: os,
            deviceName: browser,
          },
        });
      }

      return {
        ...session,
        user: {
          ...session.user,
          id: user.id,
        },
      };
    },
  },
  events: {
    async createUser({ user }) {
      await prisma.user.update({
        where: { id: user.id },
        data: {
          username: user.username ?? undefined,
        },
      });
    },
  },
  pages: {
    signIn: "/auth/signin",
  },
});

export async function detectDeviceFromHeaders() {
  const header = await headers();
  const ua = header.get("user-agent") ?? "";
  const parsed = UAParser(ua);

  const browser = parsed.browser.name ?? "Browser";
  const os = parsed.os.name ?? "OS";

  return {
    browser,
    os,
  };
}
