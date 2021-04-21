import Head from "next/head";
import Link from "next/link";
import styles from "../styles/Home.module.css";

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>Not Found | HaaS</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>Not Found</h1>
        <h5 className={styles.subtitle}>
          Would you like to go <Link href="/">back home</Link>?
        </h5>
      </main>
    </div>
  );
}
