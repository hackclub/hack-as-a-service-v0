import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import useSWR from "swr";
import { Text, useColorMode } from "@chakra-ui/react";
import DashboardLayout from "../../layouts/dashboard";
import Logs from "../../components/Logs";
import { GetServerSideProps } from "next";
import { fetchSSR } from "../../lib/fetch";
import { IApp, IBuild, IUser } from "../../types/haas";

interface IBuildEvent {
  Timestamp: number;
  Date: Date;
  Stream: "stdout" | "stderr" | "status";
  Output: string;
}

// https://stackoverflow.com/a/17415677
function toIsoString(date: Date) {
  var tzo = -date.getTimezoneOffset(),
    dif = tzo >= 0 ? "+" : "-",
    pad = function (num: number) {
      var norm = Math.floor(Math.abs(num));
      return (norm < 10 ? "0" : "") + norm;
    };

  return (
    date.getFullYear() +
    "-" +
    pad(date.getMonth() + 1) +
    "-" +
    pad(date.getDate()) +
    "T" +
    pad(date.getHours()) +
    ":" +
    pad(date.getMinutes()) +
    ":" +
    pad(date.getSeconds()) +
    dif +
    pad(tzo / 60) +
    ":" +
    pad(tzo % 60)
  );
}

function parseEvent(s: string): IBuildEvent {
  const event = JSON.parse(s);
  event.Date = new Date(event.Timestamp / 1000000);
  return event;
}

export default function BuildPage(props: {
  user: { user: IUser };
  build: { build: IBuild };
  app: { app: IApp };
}) {
  const router = useRouter();
  const { id } = router.query;

  const { colorMode } = useColorMode();

  const { data: build, mutate: mutateBuild } = useSWR(`/builds/${id}`, {
    initialData: props.build,
  });
  const { data: app } = useSWR(() => "/apps/" + build?.build.AppID, {
    initialData: props.app,
  });
  const [logs, setLogs] = useState<IBuildEvent[]>([]);

  const { data: user } = useSWR("/users/me", { initialData: props.user });

  useEffect(() => {
    if (!build) return;
    const events = ((build.build.Events ?? []) as string[]).map(parseEvent);
    events.sort((a, b) => a.Timestamp - b.Timestamp);
    setLogs(events);
  }, [build]);

  useEffect(() => {
    if (!build) return;

    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_API_BASE.replace("http", "ws")}/builds/${
        build.build.ID
      }/logs`
    );

    ws.onmessage = (e) => {
      const newEv = parseEvent(e.data);
      setLogs((old) => {
        const newLogs = old.concat(newEv);
        newLogs.sort((a, b) => a.Timestamp - b.Timestamp);
        console.log(`New event ts=${newEv.Timestamp}`);
        return newLogs;
      });
      if (newEv.Stream == "status") {
        // build ended, refresh through SWR for proper logs
        mutateBuild();
      }
    };

    return () => {
      ws.close();
    };
  }, [build]);

  return (
    <DashboardLayout
      user={user?.user}
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
        logs={logs}
        keyer={(log) => log.Output}
        render={(i) => (
          <>
            <Text
              color={colorMode == "dark" ? "snow" : "grey"}
              my={0}
              as="span"
            >
              {toIsoString(i.Date)}{" "}
            </Text>
            {i.Stream != "status" ? (
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
            )}
          </>
        )}
      />
    </DashboardLayout>
  );
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const [user, build] = await Promise.all(
      ["/users/me", `/builds/${ctx.params.id}`].map((i) => fetchSSR(i, ctx))
    );

    const app = await fetchSSR(`/apps/${build.build.AppID}`, ctx);

    return {
      props: {
        user,
        build,
        app,
      },
    };
  } catch (e) {
    if (e.url == "/users/me") {
      return {
        redirect: {
          destination: "/login",
          permanent: false,
        },
      };
    } else {
      return {
        notFound: true,
      };
    }
  }
};
