import { useRouter } from "next/router";
import { PropsWithChildren } from "react";
import useSWR from "swr";
import fetchApi from "../lib/fetch";
import DashboardLayout from "./dashboard";

export default function AppLayout({ children }: PropsWithChildren<{}>) {
  const router = useRouter();
  const { id } = router.query;

  const { data: app } = useSWR(`/apps/${id}`, fetchApi);
  const { data: team } = useSWR(() => `/teams/${app.app.TeamID}`, fetchApi);

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
                  ? `/team/${app.app.TeamID}`
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
              url: `/app/${id}`,
            },
            {
              icon: "search",
              text: "Logs",
              url: `/app/${id}/logs`,
            },
          ],
        },
      ]}
    >
      {children}
    </DashboardLayout>
  );
}
