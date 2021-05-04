import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import AppLayout from "../../../layouts/app";
import { Box, Heading, Spinner, Text } from "@chakra-ui/react";

interface ILog {
  stream: "stdout" | "stderr";
  log: string;
}

function useLogs(appId: string): { logs: ILog[]; error: string | undefined } {
  const [logs, setLogs] = useState<ILog[]>([]);

  useEffect(() => {
    if (!appId) return;

    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_API_BASE.replace(
        "http",
        "ws"
      )}/apps/${appId}/logs`
    );

    ws.onopen = () => {
      setLogs([]);
    };

    ws.onmessage = (e) => {
      setLogs((old) => old.concat(JSON.parse(e.data)));
    };

    return () => {
      ws.close();
    };
  }, [appId]);

  return { logs, error: undefined };
}

export default function AppDashboardPage() {
  const router = useRouter();
  const { id } = router.query;

  const logsElement = useRef(null);

  const { logs } = useLogs(id as string);

  useEffect(() => {
    if (logsElement.current) {
      logsElement.current.scroll({
        top: logsElement.current.scrollHeight,
        behavior: "smooth",
      });
    }
  }, [logs]);

  return (
    <AppLayout selected="Logs">
      <Box
        bg="sunken"
        borderRadius="10px"
        height="500px"
        overflow="auto"
        p={1}
        ref={logsElement}
      >
        {logs.map((i) => (
          <pre
            key={i.log}
            style={{ margin: "5px 0", padding: 0, fontSize: "inherit" }}
          >
            <Text
              color={i.stream == "stdout" ? "green" : "red"}
              my={0}
              as="span"
            >
              [{i.stream}]
            </Text>{" "}
            <Text my={0} as="span">
              {i.log}
            </Text>
          </pre>
        ))}
      </Box>
    </AppLayout>
  );
}
