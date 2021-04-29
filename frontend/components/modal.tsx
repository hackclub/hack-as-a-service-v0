import { Box, Flex, Heading } from "@theme-ui/components";
import { PropsWithChildren } from "react";

export default function Modal({
  children,
  title,
}: PropsWithChildren<{ title: string }>) {
  return (
    <Flex
      sx={{
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        alignItems: "center",
        justifyContent: "center",
      }}
      bg="rgba(0, 0, 0, 0.5)"
    >
      <Flex
        bg="background"
        sx={{ borderRadius: 10, flexDirection: "column", minWidth: 500 }}
        px={5}
        py={4}
      >
        <Heading as="h1" mb={4}>
          {title}
        </Heading>
        <Box>{children}</Box>
      </Flex>
    </Flex>
  );
}
