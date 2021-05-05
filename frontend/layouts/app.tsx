import { useRouter } from "next/router";
import { PropsWithChildren } from "react";
import useSWR from "swr";
import DashboardLayout from "./dashboard";

import { Heading } from "theme-ui";

export default function AppLayout({
  children,
  selected,
}: PropsWithChildren<{ selected: string }>) {
  const router = useRouter();
  const { id } = router.query;

  const { data: app } = useSWR(`/apps/${id}`);
  const { data: team } = useSWR(() => `/teams/${app.app.TeamID}`);

  return (
    <DashboardLayout
      title={app?.app.Name}
      sidebarSections={[
        {
          items: [
            {
              icon: "view-back",
              text: "Back",
              url:
                team?.team.Personal === false
                  ? `/teams/${app.app.TeamID}`
                  : "/dashboard",
            },
          ],
        },
        {
          title: app?.app.Name,
          items: [
            {
              icon: "explore",
              text: "Dashboard",
              url: `/apps/${id}`,
              selected: selected == "Dashboard",
            },
            {
              icon: "search",
              text: "Logs",
              url: `/apps/${id}/logs`,
              selected: selected == "Logs",
            },
            {
              icon: "share",
              text: "Deploy",
              url: `/apps/${id}/deploy`,
              selected: selected == "Deploy",
            },
          ],
        },
      ]}
    >
      <Heading as="h2" pb={4} pt={1}>
        {selected}
      </Heading>

      {children}
    </DashboardLayout>
  );
}
