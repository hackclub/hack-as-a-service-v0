import DashboardLayout, { ISidebarItem } from "../../layouts/dashboard";
import { useRouter } from "next/router";
import useSWR from "swr";
import fetchApi from "../../lib/fetch";
import { Flex, Text, Heading, Grid, Box, Avatar } from "@theme-ui/components";
import { SxProp } from "@theme-ui/core";

import Link from "next/link";

function Field({
  label,
  description,
  sx,
}: { label: string; description: string } & SxProp) {
  return (
    <Flex sx={{ flexDirection: "column", ...sx }}>
      <Text sx={{ fontSize: "20px" }}>{label}</Text>
      <Heading sx={{ fontSize: "40px" }}>{description}</Heading>
    </Flex>
  );
}

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
      <Flex
        sx={{
          alignItems: "flex-start",
          justifyContent: "center",
          flexDirection: "column",
          borderRadius: 10,
          cursor: "pointer",
        }}
        bg="sunken"
        p={30}
      >
        <Heading as="h2" sx={{ fontWeight: "normal" }}>
          {name}
        </Heading>
        <Text color="muted">({shortName})</Text>
      </Flex>
    </Link>
  );
}

function TeamMember({ name, avatar }: { name: string; avatar: string }) {
  return (
    <Flex sx={{ alignItems: "center" }} my="10px">
      <Avatar src={avatar} mr={3} />
      <Text sx={{ fontSize: 20 }}>{name}</Text>
    </Flex>
  );
}

export default function TeamPage() {
  const router = useRouter();
  const { id } = router.query;

  const { data: team } = useSWR(`/teams/${id}`, fetchApi);
  const { data: apps } = useSWR("/users/me/apps", fetchApi);

  const teamApps = apps
    ? apps.apps.filter((app: any) => app.TeamID == id)
    : null;

  return (
    <DashboardLayout
      title={team ? team.team.Name : ""}
      image={team?.team.Avatar || undefined}
      sidebarSections={[
        {
          items: [
            {
              text: "Back",
              icon: "view-back",
              url: "/dashboard",
            },
          ],
        },
        {
          title: "Apps",
          items: teamApps
            ? teamApps.map(
                (app: any): ISidebarItem => ({
                  text: app.Name,
                  icon: "code",
                  url: `/app/${app.ID}`,
                })
              )
            : [],
        },
      ]}
    >
      <Flex>
        <Field
          label="Apps"
          description={teamApps?.length}
          sx={{ marginRight: "100px" }}
        />
        <Field
          label="Users"
          description={team?.team.Users.length}
          sx={{ marginRight: "100px" }}
        />
        <Field label="Balance" description="15 HN" />
      </Flex>

      <Flex mt="35px" sx={{ flexWrap: "wrap", alignItems: "flex-start" }}>
        {teamApps && teamApps.length > 0 ? (
          <Grid
            columns="repeat(auto-fit, minmax(240px, 1fr))"
            sx={{ flex: "1 0 auto" }}
          >
            {teamApps &&
              teamApps.map((app: any) => {
                return (
                  <App
                    key={app.ID}
                    name={app.Name}
                    shortName={app.ShortName}
                    url={`/app/${app.ID}`}
                  />
                );
              })}
          </Grid>
        ) : (
          <Box sx={{ flex: 1 }}>No apps yet :(</Box>
        )}

        <Box sx={{ flex: "0 0 450px" }} p={4}>
          <Heading>Team members</Heading>

          {team &&
            team.team.Users.map((user: any) => {
              return (
                <TeamMember
                  name={user.Name}
                  avatar={user.Avatar}
                  key={user.ID}
                />
              );
            })}
        </Box>
      </Flex>
    </DashboardLayout>
  );
}
