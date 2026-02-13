"use client";

import axios from "axios";
import { useSession } from "next-auth/react";
import { useEffect, useState } from "react";

export default function Home() {
  const [data, setData] = useState({});
  const session = useSession();

  useEffect(() => {
    const fetch = async () => {
      const data = await axios
        .get("http://localhost:8080/api/ping", {
          withCredentials: true,
        })
        .then((res) => res.data)
        .catch(console.error);
      setData(data);
    };

    fetch();
  }, [setData]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      hello {JSON.stringify(data)}
      <pre>{JSON.stringify(session, null, 2)}</pre>
      lol
    </div>
  );
}
