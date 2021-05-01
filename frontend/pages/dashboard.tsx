import { useRouter } from "next/router";
import { useEffect } from "react";
import useSWR from "swr";
import { Button, Text, Box, Flex, Heading } from "theme-ui";
import DashboardLayout, {
  ISidebarItem,
  ISidebarSection,
} from "../layouts/dashboard";
import fetchApi from "../lib/fetch";

function App({ name, shortName }: { name: string; shortName: string }) {
  return (
    <Box bg="sunken" my={10} p={30} sx={{ borderRadius: 10 }}>
      <Flex sx={{ alignItems: "center" }}>
        <Heading as="h2" sx={{ fontWeight: "normal" }} mr={3}>
          {name}
        </Heading>
        <Text color="muted">({shortName})</Text>
      </Flex>
    </Box>
  );
}

export default function Dashboard() {
  const { data: teams } = useSWR("/users/me/teams", fetchApi);
  const { data: apps } = useSWR("/users/me/apps", fetchApi);

  const personalTeam = teams
    ? teams.teams.find((team: any) => team.Personal)
    : null;

  // Get personal apps
  const personalApps =
    personalTeam && apps
      ? apps.apps.filter((app: any) => app.TeamID == personalTeam.ID)
      : null;

  const sidebarSections: ISidebarSection[] = [
    {
      items: [
        {
          // image: user?.user.Avatar,
          icon: "home",
          text: "Personal Apps",
        },
      ],
    },
    {
      title: "Teams",
      items: teams
        ? teams.teams
            .filter((i: any) => !i.Personal)
            .map((i: any) => ({
              icon: "person",
              image: i.Avatar || undefined,
              text: i.Name,
              url: `/team/${i.ID}`,
            }))
        : [],
    },
  ];

  return (
    <DashboardLayout title="Personal Apps" sidebarSections={sidebarSections}>
      {personalApps &&
        personalApps.map((app: any) => {
          return <App name={app.Name} shortName={app.ShortName} key={app.ID} />;
        })}
    </DashboardLayout>
  );
}
