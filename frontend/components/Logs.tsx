import React, { Ref, useEffect, useRef } from "react";
import { Box, useColorMode } from "@chakra-ui/react";

export type LogsProps<T> = {
  render: (item: T) => React.ReactElement;
  keyer: (item: T) => React.Key;
  logs: T[];
};

export default function Logs<T>({ logs, render, keyer }: LogsProps<T>) {
  const { colorMode } = useColorMode();

  const logsElement = useRef(null);

  useEffect(() => {
    if (logsElement.current) {
      logsElement.current.scroll({
        top: logsElement.current.scrollHeight,
        behavior: "smooth",
      });
    }
  }, [logs]);

  return (
    <Box
      ref={logsElement}
      bg={colorMode == "dark" ? "steel" : "sunken"}
      borderRadius="10px"
      maxWidth="100%"
      height="500px"
      overflowY="auto"
      p={1}
      fontSize="sm"
    >
      {logs.map((log) => (
        <pre
          key={keyer(log)}
          style={{
            margin: "5px 0",
            padding: 0,
            fontSize: "inherit",
            overflow: "visible",
          }}
        >
          {render(log)}
        </pre>
      ))}
    </Box>
  );
}
