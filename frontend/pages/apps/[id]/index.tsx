import { GetServerSideProps } from "next";
import { useRouter } from "next/router";
import useSWR from "swr";
import AppLayout from "../../../layouts/app";
import fetchApi, { fetchSSR } from "../../../lib/fetch";
import { IApp, ITeam, IUser } from "../../../types/haas";

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

  return (
    <AppLayout
      selected="Dashboard"
      user={user.user}
      app={app.app}
      team={team.team}
    >
      App dashboard
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
