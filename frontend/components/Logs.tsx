import React, { Ref } from "react";
import { Box, useColorMode } from "@chakra-ui/react";

export type LogsProps<T> = {
  render: (item: T) => React.ReactElement;
  keyer: (item: T) => React.Key;
  logs: T[];
};

type WithRef<P, R> = P & { ref?: Ref<R> };

export default function Logs<T>(props: WithRef<LogsProps<T>, HTMLDivElement>) {
  const Component = React.forwardRef(
    ({ render, keyer, logs }: LogsProps<T>, ref: Ref<HTMLDivElement>) => {
      const { colorMode } = useColorMode();

      return (
        <Box
          bg={colorMode == "dark" ? "steel" : "sunken"}
          borderRadius="10px"
          maxWidth="100%"
          height="500px"
          overflowY="auto"
          p={1}
          ref={ref}
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
  );

  return <Component {...props} />;
}
