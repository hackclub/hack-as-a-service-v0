import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import AppLayout from "../../../layouts/app";
import { Text, useColorMode } from "@chakra-ui/react";
import Logs from "../../../components/Logs";
import { GetServerSideProps } from "next";
import fetchApi from "../../../lib/fetch";
import { IApp, ITeam, IUser } from "../../../types/haas";
import useSWR from "swr";

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

export default function AppDashboardPage(props: {
  user: { user: IUser };
  app: { app: IApp };
  team: { team: ITeam };
}) {
  const router = useRouter();
  const { id } = router.query;

  const { data: user } = useSWR("/users/me", { initialData: props.user });
  const { data: app } = useSWR(`/apps/${id}`, { initialData: props.app });
  const { data: team } = useSWR(() => "/teams/" + app.app.TeamID, {
    initialData: props.team,
  });

  const { colorMode } = useColorMode();

  const { logs } = useLogs(id as string);

  return (
    <AppLayout selected="Logs" user={user.user} app={app.app} team={team.team}>
      <Logs
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

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const [user, app] = await Promise.all(
      ["/users/me", `/apps/${ctx.params.id}`].map((i) =>
        fetchApi(i, { headers: ctx.req.headers as HeadersInit })
      )
    );

    const team = await fetchApi(`/teams/${app.app.TeamID}`, {
      headers: ctx.req.headers as HeadersInit,
    });

    return {
      props: {
        user,
        app,
        team,
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
