import Head from "next/head";
import useSWR from "swr";
import styles from "../styles/Home.module.css";

export default function Home() {
  const { data, error } = useSWR("/users/me", (key: string) =>
    fetch(process.env.NEXT_PUBLIC_API_BASE + key, {
      credentials: "include",
    }).then((resp) => {
      if (!resp.ok) {
        throw new Error("error");
      }
      return resp.json();
    })
  );

  return (
    <div className={styles.container}>
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        {data && (
          <div style={{ display: "flex", alignItems: "center" }}>
            <img
              src={data.user.Avatar}
              style={{
                width: "50px",
                borderRadius: "50%",
                marginRight: "10px",
              }}
            />
            <h3>{data.user.Name}</h3>
          </div>
        )}
        <h1 className={styles.title}>Coming Soon</h1>
        <h5 className={styles.subtitle}>
          Hack as a Service | A <a href="https://hackclub.com">Hack Club</a>{" "}
          project
        </h5>
      </main>
    </div>
  );
}
