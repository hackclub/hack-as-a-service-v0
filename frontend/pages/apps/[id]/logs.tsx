import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import AppLayout from "../../../layouts/app";
import { Text, useColorMode } from "@chakra-ui/react";
import Logs from "../../../components/Logs";

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

  const { colorMode } = useColorMode();

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
      <Logs
        ref={logsElement}
        logs={logs}
        keyer={(log) => log.log}
        render={(i) => (
          <>
            <Text
              color={i.stream == "stdout" ? "green" : "red"}
              my={0}
              as="span"
            >
              [{i.stream}]
            </Text>{" "}
            <Text
              my={0}
              as="span"
              color={colorMode == "dark" ? "background" : "text"}
            >
              {i.log}
            </Text>
          </>
        )}
      />
    </AppLayout>
  );
}
