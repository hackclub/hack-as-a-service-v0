import useSWR from "swr";
import { Button, Text, Box, Flex, Heading } from "theme-ui";
import DashboardLayout, { SidebarSection } from "../layouts/dashboard";
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
  const { data: user } = useSWR("/users/me", fetchApi);
  const { data: teams } = useSWR("/users/me/teams", fetchApi);
  const { data: apps } = useSWR("/users/me/apps", fetchApi);

  // Get personal apps
  const personalApps =
    teams && apps
      ? apps.apps.filter(
          (app: any) =>
            app.TeamID == teams.teams.find((team) => team.Personal).ID
        )
      : null;

  const sidebarSections: SidebarSection[] = [
    {
      items: [
        {
          image: user?.user.Avatar,
          icon: "person",
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
              text: i.Name,
            }))
        : [],
    },
  ];

  return (
    <DashboardLayout title="Personal Apps" sidebarSections={sidebarSections}>
      {personalApps &&
        personalApps.map((v, i) => {
          return <App name={v.Name} shortName={v.ShortName} key={i} />;
        })}
    </DashboardLayout>
  );
}
