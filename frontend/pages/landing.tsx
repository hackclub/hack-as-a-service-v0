import { Box, Button, Flex, Image, Text } from "@theme-ui/components";
import Head from "next/head";
import Link from "next/link";
import styles from "../styles/Home.module.css";

export default function Home() {
  return (
    <Box backgroundColor="white" color="#000000" sx={{ minHeight: "100vh" }}>
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Flex
        sx={{
          alignItems: "center",
          justifyContent: "space-between",
          px: "4",
          pt: "3",
          position: "absolute",
          width: "100%",
        }}
      >
        <Image src="/nav-vector.svg" />
        <Link href="/login">
          <Button backgroundColor="#EC3750" px="4">
            Login
          </Button>
        </Link>
      </Flex>
      <Box
        className={styles.container}
        sx={{
          backgroundImage: "url(/landing-vector.svg)",
          backgroundRepeat: "no-repeat",
          backgroundPosition: "center",
        }}
      >
        <Text
          as="h1"
          sx={{
            fontSize: "5.5rem",
            textAlign: "center",
            py: "6",
          }}
        >
          A managed platform for makers.
        </Text>
        <Link href="/dashboard">
          <Button backgroundColor="#EC3750" px="4">
            Start Building
          </Button>
        </Link>
      </Box>
    </Box>
  );
}
