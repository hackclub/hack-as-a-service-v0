import Link from "next/link";
import useSWR from "swr";
import { Box, Flex, Heading, Text } from "@chakra-ui/react";
import DashboardLayout, { ISidebarSection } from "../layouts/dashboard";

function App({
  name,
  shortName,
  url,
}: {
  name: string;
  shortName: string;
  url: string;
}) {
  return (
    <Link href={url}>
      <Box
        bg="sunken"
        my={10}
        p={30}
        sx={{ borderRadius: 10, cursor: "pointer" }}
      >
        <Flex sx={{ alignItems: "center" }}>
          <Heading as="h2" sx={{ fontWeight: "normal" }} mr={3}>
            {name}
          </Heading>
          <Text color="muted">({shortName})</Text>
        </Flex>
      </Box>
    </Link>
  );
}

export default function Dashboard() {
  const { data: teams } = useSWR("/users/me/teams");
  const { data: personalTeam } = useSWR("/teams/me");

  const teamList = teams?.teams
    .filter((i: any) => !i.Personal)
    .map((i: any) => ({
      icon: "person",
      image: i.Avatar || undefined,
      text: i.Name,
      url: `/teams/${i.ID}`,
    }));

  const sidebarSections: ISidebarSection[] = [
    {
      items: [
        {
          // image: user?.user.Avatar,
          icon: "home",
          text: "Personal Apps",
          url: "/dashboard",
        },
      ],
    },
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
    <DashboardLayout title="Personal Apps" sidebarSections={sidebarSections}>
      {personalTeam &&
        (personalTeam.team.Apps.length > 0 ? (
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
          <Heading as="h3" sx={{ fontWeight: "normal" }} mt={3}>
            You don't have any personal apps quite yet. 😢
          </Heading>
        ))}
    </DashboardLayout>
  );
}
