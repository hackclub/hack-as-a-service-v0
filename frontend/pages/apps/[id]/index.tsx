import { GetServerSideProps } from "next";
import { useRouter } from "next/router";
import useSWR from "swr";
import AppLayout from "../../../layouts/app";
import fetchApi, { fetchSSR } from "../../../lib/fetch";
import { Button, Text } from "@chakra-ui/react";
import { IApp, ILetsEncrypt, ITeam, IUser } from "../../../types/haas";

export default function AppDashboardPage(props: {
  user: { user: IUser };
  app: { app: IApp };
  letsEncrypt: { letsencrypt: ILetsEncrypt };
  team: { team: ITeam };
}) {
  const router = useRouter();
  const { id } = router.query;

  const { data: user } = useSWR("/users/me", { initialData: props.user });
  const { data: app } = useSWR(`/apps/${id}`, { initialData: props.app });
  const { data: letsEncrypt, mutate: mutateLetsEncrypt } = useSWR(
    `/apps/${id}/letsencrypt`,
    {
      initialData: props.letsEncrypt,
    }
  );
  const { data: team } = useSWR(() => "/teams/" + app.app.TeamID, {
    initialData: props.team,
  });

  async function enableLetsEncrypt() {
    mutateLetsEncrypt({
      ...letsEncrypt,
      letsencrypt: { ...letsEncrypt.letsencrypt, LetsEncryptEnabled: true },
    });
    await fetchApi(`/apps/${id}/letsencrypt/enable`, {
      method: "POST",
    });
    mutateLetsEncrypt();
  }

  return (
    <AppLayout
      selected="Dashboard"
      user={user.user}
      app={app.app}
      team={team.team}
    >
      App dashboard
      <Text my={1}>
        Lets Encrypt enabled:{" "}
        {letsEncrypt.letsencrypt.LetsEncryptEnabled ? "yes" : "no"}
        {!letsEncrypt.letsencrypt.LetsEncryptEnabled && (
          <>
            <br />
            <Button onClick={enableLetsEncrypt} px={2}>
              Enable Lets Encrypt
            </Button>
          </>
        )}
      </Text>
    </AppLayout>
  );
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const [user, app, letsEncrypt] = await Promise.all(
      [
        "/users/me",
        `/apps/${ctx.params.id}`,
        `/apps/${ctx.params.id}/letsencrypt`,
      ].map((i) => fetchSSR(i, ctx))
    );

    const team = await fetchSSR(`/teams/${app.app.TeamID}`, ctx);

    return {
      props: {
        user,
        app,
        letsEncrypt,
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
