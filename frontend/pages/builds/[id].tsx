import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import useSWR from "swr";
import { Box, Text } from "@chakra-ui/react";
import DashboardLayout from "../../layouts/dashboard";

interface IBuildEvent {
  Stream: "stdout" | "stderr" | "status";
  Output: string;
}

export default function BuildPage() {
  const router = useRouter();
  const { id } = router.query;

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
      <Box
        bg="sunken"
        sx={{ borderRadius: 10, height: 500, overflow: "auto" }}
        p={3}
        ref={logsElement}
      >
        {logs.map((i) => (
          <pre key={i.Output} style={{ margin: "5px 0" }}>
            {i.Stream != "status" ? (
              <>
                <Text color={i.Stream == "stdout" ? "green" : "red"}>
                  [{i.Stream}]
                </Text>{" "}
                <span>{i.Output}</span>
              </>
            ) : (
              <>
                <Text color="grey">[Build exited with status {i.Output}]</Text>
              </>
            )}
          </pre>
        ))}
      </Box>
    </DashboardLayout>
  );
}
