import { PropsWithChildren } from "react";
import DashboardLayout from "./dashboard";

import { Heading } from "@chakra-ui/react";
import { IApp, ITeam, IUser } from "../types/haas";

export default function AppLayout({
  children,
  selected,
  app,
  user,
  team,
}: PropsWithChildren<{
  selected: string;
  app: IApp | undefined;
  user: IUser | undefined;
  team: ITeam | undefined;
}>) {
  return (
    <DashboardLayout
      title={app?.Name}
      user={user}
      sidebarSections={[
        {
          items: [
            {
              icon: "view-back",
              text: "Back",
              url:
                team?.Personal === false
                  ? `/teams/${app?.TeamID}`
                  : "/dashboard",
            },
          ],
        },
        {
          title: app?.Name,
          items: [
            {
              icon: "explore",
              text: "Dashboard",
              url: `/apps/${app?.ID}`,
              selected: selected == "Dashboard",
            },
            {
              icon: "search",
              text: "Logs",
              url: `/apps/${app?.ID}/logs`,
              selected: selected == "Logs",
            },
            {
              icon: "share",
              text: "Deploy",
              url: `/apps/${app?.ID}/deploy`,
              selected: selected == "Deploy",
            },
          ],
        },
      ]}
    >
      <Heading as="h2" pb={2} pt={1}>
        {selected}
      </Heading>

      {children}
    </DashboardLayout>
  );
}
