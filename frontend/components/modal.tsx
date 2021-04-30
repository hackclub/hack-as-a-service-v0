import { Box, Close, Flex, Heading } from "@theme-ui/components";
import { PropsWithChildren } from "react";

export default function Modal({
  children,
  title,
  visible,
  onClose,
}: PropsWithChildren<{
  title: string;
  visible: boolean;
  onClose: () => void;
}>) {
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
        display: visible ? "flex" : "none !important",
      }}
      bg="rgba(0, 0, 0, 0.5)"
    >
      <Flex
        bg="background"
        sx={{
          borderRadius: 10,
          flexDirection: "column",
          minWidth: 500,
          position: "relative",
        }}
        px={5}
        py={4}
      >
        <Close
          mr={2}
          sx={{ position: "absolute", left: 10, top: 10, cursor: "pointer" }}
          onClick={() => onClose()}
        />
        <Box>{children}</Box>
      </Flex>
    </Flex>
  );
}
