import Head from "next/head";
import styles from "../styles/Home.module.css";
import { Text, Card, Box, Heading, Button } from "theme-ui";
import Icon from '@hackclub/icons'

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>Hack as a Service</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box>
        <Card
          sx={{
            width: "100%",
            maxWidth: "600px",
            margin: "auto",
            textAlign: "left",
            bg: "elevated",
            p: [0, 0, 0],
            borderRadius: "extra",
          }}
        >
          <Box sx={{ p: [3, 4, 4], pb: [2, 2, 2] }}>
            <Box pb={1}>
              <Heading>Welcome to</Heading>
              <Heading as="h1" sx={{ fontSize: "4em" }}>
                Hack as a Service
              </Heading>
            </Box>
            <Box pb={3}>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
              eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
              enim ad minim veniam, quis nostrud exercitation ullamco laboris
              nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in
              reprehenderit in voluptate velit esse cillum dolore eu fugiat
              nulla.
            </Box>
          </Box>
          <Button
            variant="cta"
            sx={{
              width: "100%",
              borderRadius: "extra",
              borderTopLeftRadius: "0",
              borderTopRightRadius: "0",
              textAlign: "left",
              ":focus,:hover": {
                boxShadow: "none",
                transform: "scale(1)",
              },
            }}
          >
            <Text sx={{ width: "100%", px: ["0px", "16px", "16px"], display: 'inline-flex' }}>
              <Text>Login with Slack</Text>
              <Text sx={{ width: 'max-content', textAlign: 'right', flexGrow: 1, height: '30px'}}><Icon glyph="enter" size={30} /></Text>
            </Text>
          </Button>
        </Card>
      </Box>
    </div>
  );
}
