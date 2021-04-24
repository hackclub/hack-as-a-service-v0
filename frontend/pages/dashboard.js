import Head from "next/head";
import { Text, Card, Box, Heading, Button, Grid, Badge } from "theme-ui";
import Icon from "@hackclub/icons";
import { useState } from "react";

function isOdd(num) {
  return num % 2;
}

function AppBadge({ setOpenBadge, x, openBadge }) {
  return (
    <Badge
      sx={{
        borderRadius: 0,
        mb: 0,
        width: "100%",
        p: 2,
        px: 3,
        fontSize: "18px",
        bg: openBadge == x ? "elevated" : "sunken",
        cursor: "pointer",
        color: openBadge == x ? "text" : "secondary",

        ":focus,:hover": {
          bg: "elevated",
          color: "text",
        },
      }}
      onClick={() => setOpenBadge(openBadge == x ? "1000000" : x)}
    >
      <Text sx={{ fontSize: "0.6em", mr: "2px" }}>▶︎</Text> sarthak-discord
      {openBadge == x ? (
        <Box sx={{ my: 2 }}>
          <Box>
            <Text sx={{ fontSize: "0.6em", mr: "2px", color: "elevated" }}>
              ▶︎
            </Text>{" "}
            <Text sx={{ verticalAlign: "middle" }}>
              <Icon
                glyph="docs"
                size={18}
                style={{ transform: "translateY(3.5px)", marginRight: "4px" }}
              />{" "}
              <Text sx={{ height: "18px", fontWeight: "500" }}>Overview</Text>{" "}
            </Text>
          </Box>
          <Box>
            <Text sx={{ fontSize: "0.6em", mr: "2px", color: "elevated" }}>
              ▶︎
            </Text>{" "}
            <Text sx={{ verticalAlign: "middle" }}>
              <Icon
                glyph="private-fill"
                size={18}
                style={{ transform: "translateY(4px)", marginRight: "4px" }}
              />{" "}
              <Text sx={{ height: "18px", fontWeight: "500" }}>Secrets</Text>{" "}
            </Text>
          </Box>
          <Box
            sx={{
              width: '100%',
              ":focus,:hover": {
                textDecoration: "underline",
              },
            }}
          >
            <Text sx={{ fontSize: "0.6em", mr: "2px", color: "elevated", textDecorationColor: "elevated" }}>
              ▶︎
            </Text>{" "}
            <Text sx={{ verticalAlign: "middle" }}>
              <Icon
                glyph="link"
                size={18}
                style={{ transform: "translateY(3.5px)", marginRight: "4px" }}
              />{" "}
              <Text
                sx={{
                  height: "18px",
                  fontWeight: "500",
                }}
              >
                Domains
              </Text>{" "}
            </Text>
          </Box>
          <Box>
            <Text sx={{ fontSize: "0.6em", mr: "2px", color: "elevated" }}>
              ▶︎
            </Text>{" "}
            <Text sx={{ verticalAlign: "middle" }}>
              <Icon
                glyph="settings"
                size={18}
                style={{ transform: "translateY(3.5px)", marginRight: "4px" }}
              />{" "}
              <Text sx={{ height: "18px", fontWeight: "500" }}>Settings</Text>{" "}
            </Text>
          </Box>
        </Box>
      ) : (
        ""
      )}
    </Badge>
  );
}

export default function Home() {
  const [openBadge, setOpenBadge] = useState("");
  return (
    <div>
      <Grid columns="1fr 4fr" sx={{ height: "100vh" }}>
        <Box
          sx={{
            bg: "sunken",
            height: "100vh",
            overflow: "scroll",
            p: 4,
            px: 0,
          }}
        >
          <Heading sx={{ px: 3, pb: 3 }} as="h1">
            HaaS
          </Heading>
          {Array.from(Array(100).keys()).map((x) => (
            <AppBadge setOpenBadge={setOpenBadge} x={x} openBadge={openBadge} />
          ))}
        </Box>
        <Box></Box>
      </Grid>
    </div>
  );
}
