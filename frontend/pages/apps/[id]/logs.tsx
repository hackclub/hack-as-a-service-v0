import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import AppLayout from "../../../layouts/app";
import { Box, Heading, Spinner, Text } from "theme-ui";

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
        sx={{ borderRadius: 10, height: 500, overflow: "auto" }}
        p={3}
        ref={logsElement}
      >
        {logs.map((i) => (
          <pre key={i.log} style={{ margin: "5px 0" }}>
            <Text color={i.stream == "stdout" ? "green" : "red"}>
              [{i.stream}]
            </Text>{" "}
            <span>{i.log}</span>
          </pre>
        ))}
      </Box>
    </AppLayout>
  );
}
