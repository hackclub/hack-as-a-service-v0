import { GetServerSideProps } from "next";
import { useRouter } from "next/router";
import useSWR from "swr";
import AppLayout from "../../../layouts/app";
import fetchApi, { fetchSSR } from "../../../lib/fetch";
import {
  Text,
  Button,
  Flex,
  IconButton,
  Input,
  useToast,
} from "@chakra-ui/react";
import { IApp, ITeam, IUser } from "../../../types/haas";
import { useState } from "react";
import Icon from "@hackclub/icons";
import { ErrorToast, SuccessToast } from "../../../components/Toast";

function EnvVar({
  envVar,
  value,
  onRemove,
  onKeyChange,
  onValueChange,
  disabled,
}: {
  envVar: string;
  value: string;
  onRemove: () => void;
  onKeyChange: (key: string) => void;
  onValueChange: (value: string) => void;
  disabled?: boolean;
}) {
  return (
    <Flex my={1}>
      <Input
        flex={1}
        placeholder="Key"
        value={envVar}
        mr={1}
        onChange={(e) => onKeyChange(e.target.value)}
        disabled={disabled}
      />
      <Input
        flex={2}
        placeholder="Value"
        value={value}
        mr={2}
        onChange={(e) => onValueChange(e.target.value)}
        disabled={disabled}
      />
      <IconButton
        aria-label="Remove environment variable"
        icon={<Icon glyph="delete" />}
        disabled={disabled}
        onClick={() => onRemove()}
      />
    </Flex>
  );
}

export default function AppDashboardPage(props: {
  user: { user: IUser };
  app: { app: IApp };
  team: { team: ITeam };
  env: { env: { [key: string]: string } };
}) {
  const router = useRouter();
  const { id } = router.query;

  const { data: user } = useSWR("/users/me", { initialData: props.user });
  const { data: app } = useSWR(`/apps/${id}`, { initialData: props.app });
  const { data: team } = useSWR(() => "/teams/" + app.app.TeamID, {
    initialData: props.team,
  });

  const [env, setEnv] = useState<{ key: string; value: string; id: string }[]>(
    Object.entries(props.env.env).map(([key, value]) => ({
      key,
      value,
      id: Math.random().toString(),
    }))
  );
  const [loading, setLoading] = useState<string | null>(null);

  const toast = useToast();

  return (
    <AppLayout
      selected="Environment"
      user={user.user}
      app={app.app}
      team={team.team}
    >
      <Flex mb={2}>
        <Button
          px={2}
          mr={1}
          isDisabled={!!loading}
          onClick={() =>
            setEnv((e) =>
              e.concat({ key: "", value: "", id: Math.random().toString() })
            )
          }
        >
          Add Pair
        </Button>
        <Button
          px={2}
          variant="cta"
          isLoading={!!loading}
          loadingText={loading}
          onClick={async () => {
            setLoading("Saving...");

            try {
              await fetchApi(`/apps/${id}/env`, {
                method: "PUT",
                body: JSON.stringify({
                  Env: env.reduce((acc, { key, value }) => {
                    acc[key] = value;

                    return acc;
                  }, {}),
                }),
              });

              setLoading("Restarting app...");

              await fetchApi(`/apps/${id}/restart`, {
                method: "POST",
              });

              toast({
                status: "success",
                duration: 5000,
                position: "top",
                render: () => (
                  <SuccessToast text="Your app's environment has successfully been updated." />
                ),
              });
            } catch (e) {
              toast({
                status: "error",
                duration: 5000,
                position: "top",
                render: () => (
                  <ErrorToast text={e.message || "Something went wrong."} />
                ),
              });
            } finally {
              setLoading(null);
            }
          }}
        >
          Save
        </Button>
      </Flex>
      <Text color="gray" mt={0} mb={2}>
        Changing environment variables requires an app restart, which may take
        up to 60 seconds.
      </Text>
      {env.map(({ key, value, id }) => (
        <EnvVar
          key={id}
          envVar={key}
          value={value}
          disabled={!!loading}
          onKeyChange={(e) => {
            setEnv(
              env.map((x) => {
                if (x.id == id) {
                  x.key = e;
                }

                return x;
              })
            );
          }}
          onValueChange={(e) => {
            setEnv(
              env.map((x) => {
                if (x.id == id) {
                  x.value = e;
                }

                return x;
              })
            );
          }}
          onRemove={() => {
            setEnv((e) => e.filter((x) => x.id != id));
          }}
        />
      ))}
    </AppLayout>
  );
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const [user, app, env] = await Promise.all(
      [
        "/users/me",
        `/apps/${ctx.params.id}`,
        `/apps/${ctx.params.id}/env`,
      ].map((i) => fetchSSR(i, ctx))
    );

    const team = await fetchSSR(`/teams/${app.app.TeamID}`, ctx);

    return {
      props: {
        user,
        app,
        team,
        env,
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
