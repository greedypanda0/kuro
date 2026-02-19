"use client";

import { Prisma } from "@@/prisma/generated/prisma/browser";
import axios from "axios";
import { signIn, useSession } from "next-auth/react";
import { useEffect, useState } from "react";

export default function Home() {
  const [data, setData] = useState({});
  const session = useSession();

  const fetch = async () => {
    const data = await axios
      .post("http://localhost:3000/api/auth/token", {
        // withCredentials: true,
        expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24).toISOString(),
      })
      .then((res) => res.data)
      .catch(console.error);

    setData(data);
  };

  return (
    <div className="flex flex-col min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <button onClick={() => signIn("github")}>SignIN</button>
      Json {session.data?.user?.name}
      <pre>{JSON.stringify(data, null, 2)}</pre>
      <div>
        <button onClick={() => fetch()}>Create Token</button>
      </div>
    </div>
  );
}
