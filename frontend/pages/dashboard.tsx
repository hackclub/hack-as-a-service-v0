import useSWR from "swr";
import { Heading } from "@chakra-ui/react";
import App from "../components/App";
import DashboardLayout, {
  ISidebarItem,
  ISidebarSection,
} from "../layouts/dashboard";
import { GetServerSideProps } from "next";
import fetchApi from "../lib/fetch";
import { ITeam, IUser } from "../types/haas";
import Head from "next/head";

export default function Dashboard(props: {
  user: { user: IUser };
  teams: { teams: ITeam[] };
  personalTeam: { team: ITeam };
}) {
  const { data: teams } = useSWR("/users/me/teams", {
    initialData: props.teams,
  });
  const { data: personalTeam } = useSWR("/teams/me", {
    initialData: props.personalTeam,
  });
  const { data: user } = useSWR("/users/me", { initialData: props.user });

  const teamList = teams.teams
    .filter((i: ITeam) => !i.Personal)
    .map(
      (i: ITeam): ISidebarItem => ({
        icon: "person",
        image: i.Avatar || undefined,
        text: i.Name,
        url: `/teams/${i.ID}`,
      })
    );

  const sidebarSections: ISidebarSection[] = [
    {
      title: "Teams",
      items: teamList
        ? teamList.length > 0
          ? teamList
          : [{ text: "You're not a part of any teams." }]
        : [],
    },
  ];

  return (
    <>
      <Head>
        <title>HaaS Dashboard</title>
      </Head>
      <DashboardLayout
        title="Personal Apps"
        sidebarSections={sidebarSections}
        user={user.user}
      >
        {personalTeam.team.Apps.length > 0 ? (
          personalTeam.team.Apps.map((app: any) => {
            return (
              <App
                url={`/apps/${app.ID}`}
                name={app.Name}
                shortName={app.ShortName}
                key={app.ID}
              />
            );
          })
        ) : (
          <Heading as="h3" size="sm" fontWeight="normal" mt={1}>
            You don't have any personal apps quite yet. 😢
          </Heading>
        )}
      </DashboardLayout>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const user = await fetchApi("/users/me", {
      headers: ctx.req.headers as HeadersInit,
    });

    const teams = await fetchApi("/users/me/teams", {
      headers: ctx.req.headers as HeadersInit,
    });

    const personalTeam = await fetchApi("/teams/me", {
      headers: ctx.req.headers as HeadersInit,
    });

    return {
      props: {
        user,
        teams,
        personalTeam,
      },
    };
  } catch (e) {
    return {
      redirect: {
        destination: "/login",
        permanent: false,
      },
    };
  }
};
