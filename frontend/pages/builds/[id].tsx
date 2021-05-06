import { useRouter } from "next/router";
import { useEffect, useRef, useState } from "react";
import useSWR from "swr";
import { Text, useColorMode } from "@chakra-ui/react";
import DashboardLayout from "../../layouts/dashboard";
import Logs from "../../components/Logs";
import { GetServerSideProps } from "next";
import { fetchSSR } from "../../lib/fetch";
import { IApp, IBuild, IUser } from "../../types/haas";

interface IBuildEvent {
  Stream: "stdout" | "stderr" | "status";
  Output: string;
}

export default function BuildPage(props: {
  user: { user: IUser };
  build: { build: IBuild };
  app: { app: IApp };
}) {
  const router = useRouter();
  const { id } = router.query;

  const { colorMode } = useColorMode();

  const { data: build } = useSWR(`/builds/${id}`, { initialData: props.build });
  const { data: app } = useSWR(() => "/apps/" + build?.build.AppID, {
    initialData: props.app,
  });
  const [logs, setLogs] = useState<IBuildEvent[]>([]);

  const { data: user } = useSWR("/users/me", { initialData: props.user });

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
