import { Box, Button, Input, Text } from "@chakra-ui/react";
import { GetServerSideProps } from "next";
import { useRouter } from "next/router";
import { FormEvent, useRef } from "react";
import useSWR from "swr";
import AppLayout from "../../../layouts/app";
import fetchApi, { fetchSSR } from "../../../lib/fetch";
import { IApp, ITeam, IUser } from "../../../types/haas";

export default function AppDeployPage(props: {
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

  const repoUrlRef = useRef<HTMLInputElement>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const url = repoUrlRef.current.value;
    const res = await fetchApi(`/apps/${id}/deploy`, {
      method: "POST",
      body: JSON.stringify({
        GitRepository: url,
      }),
    });
    router.push(`/builds/${res.build.ID}`);
  }

  return (
    <AppLayout
      selected="Deploy"
      user={user.user}
      app={app.app}
      team={team.team}
    >
      <Box as="form" onSubmit={onSubmit}>
        <Text htmlFor="repoUrl" as="label" my={0}>
          Git repository URL
          <br />
          <Text color="grey" size="xs" my={0}>
            Must be a public repository
          </Text>
        </Text>
        <Input name="repoUrl" type="url" required ref={repoUrlRef} />
        <Button variant="ctaLg" mt={2} type="submit">
          Deploy
        </Button>
      </Box>
    </AppLayout>
  );
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const [user, app] = await Promise.all(
      ["/users/me", `/apps/${ctx.params.id}`].map((i) => fetchSSR(i, ctx))
    );

    const team = await fetchSSR(`/teams/${app.app.TeamID}`, ctx);

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
