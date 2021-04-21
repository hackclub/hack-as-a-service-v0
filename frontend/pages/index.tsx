import Head from "next/head";
import Link from "next/link";
import styles from "../styles/Home.module.css";

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>Coming Soon</h1>
        <h5 className={styles.subtitle}>
          Hack as a Service | A <a href="https://hackclub.com">Hack Club</a>{" "}
          project
        </h5>
      </main>
    </div>
  );
}
