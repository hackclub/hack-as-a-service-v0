import { Box, Button, Flex, Image, Text } from "@chakra-ui/react";
import Head from "next/head";
import Link from "next/link";

export default function Home() {
  return (
    <Box backgroundColor="white" color="#000000" minHeight="100vh">
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Flex
        alignItems="center"
        justifyContent="space-between"
        px={2}
        pt={1}
        position="absolute"
        width="100%"
      >
        <Image src="/nav-vector.svg" />
        <Link href="/login">
          <Button backgroundColor="#EC3750" px={2}>
            Login
          </Button>
        </Link>
      </Flex>
      <Flex
        height="100vh"
        justifyContent="center"
        alignItems="center"
        flexDirection="column"
        backgroundImage="url(/landing-vector.svg)"
        backgroundRepeat="no-repeat"
        backgroundPosition="center"
      >
        <Text as="h1" fontSize="5.5rem" textAlign="center" py={3}>
          A managed platform for makers.
        </Text>
        <Link href="/dashboard">
          <Button backgroundColor="#EC3750" px={2}>
            Start Building
          </Button>
        </Link>
      </Flex>
    </Box>
  );
}
