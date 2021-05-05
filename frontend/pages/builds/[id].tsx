import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import useSWR from "swr";
import { Text, useColorMode } from "@chakra-ui/react";
import DashboardLayout from "../../layouts/dashboard";
import Logs from "../../components/Logs";

interface IBuildEvent {
  Stream: "stdout" | "stderr" | "status";
  Output: string;
}

export default function BuildPage() {
  const router = useRouter();
  const { id } = router.query;

  const { colorMode } = useColorMode();

  const logsElement = useRef(null);

  const { data: build } = useSWR(`/builds/${id}`);
  const { data: app } = useSWR(() => `/apps/${build?.build.AppID}`);
  const [logs, setLogs] = useState<IBuildEvent[]>([]);

  useEffect(() => {
    if (!build) return;
    setLogs(build.build.Events.map((e: string) => JSON.parse(e)));
  }, [build]);

  useEffect(() => {
    if (!build) return;

    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_API_BASE.replace("http", "ws")}/builds/${
        build.build.ID
      }/logs`
    );

    ws.onmessage = (e) => {
      setLogs((old) => old.concat(JSON.parse(e.data)));
    };

    return () => {
      ws.close();
    };
  }, [build]);

  useEffect(() => {
    if (logsElement.current) {
      logsElement.current.scroll({
        top: logsElement.current.scrollHeight,
        behavior: "smooth",
      });
    }
  }, [logs]);

  return (
    <DashboardLayout
      title={`Build ${build?.build.ID} for app ${app?.app.Name}`}
      sidebarSections={
        app
          ? [
              {
                items: [
                  {
                    text: "Back",
                    icon: "view-back",
                    url: `/apps/${app.app.ID}`,
                  },
                ],
              },
            ]
          : []
      }
    >
      <Logs
        ref={logsElement}
        logs={logs}
        keyer={(log) => log.Output}
        render={(i) =>
          i.Stream != "status" ? (
            <>
              <Text
                color={i.Stream == "stdout" ? "green" : "red"}
                my={0}
                as="span"
              >
                [{i.Stream}]
              </Text>{" "}
              <Text
                my={0}
                as="span"
                color={colorMode == "dark" ? "background" : "text"}
              >
                {i.Output}
              </Text>
            </>
          ) : (
            <>
              <Text
                color={colorMode == "dark" ? "snow" : "grey"}
                my={0}
                as="span"
              >
                [Build exited with status {i.Output}]
              </Text>
            </>
          )
        }
      />
    </DashboardLayout>
  );
}
