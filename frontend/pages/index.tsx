import { Flex, Heading, Text } from "@chakra-ui/react";
import Head from "next/head";

export default function Home() {
  return (
    <>
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Flex
        as="main"
        height="100vh"
        justifyContent="center"
        alignItems="center"
        flexDirection="column"
      >
        <Heading as="h1" fontSize="8rem" lineHeight="1.15">
          Coming Soon
        </Heading>
        <Text my={1}>
          Hack as a Service | A <a href="https://hackclub.com">Hack Club</a>{" "}
          project
        </Text>
      </Flex>
    </>
  );
}
