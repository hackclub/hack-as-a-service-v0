import AppLayout from "../../../../layouts/app";

import { Stat } from "../../../../components/Stat";
import { Flex, Button } from "@theme-ui/components";
import { useRouter } from "next/router";

export default function AppAddonOverview() {
  return (
    <AppLayout selected="Dashboard">
      <Flex>
        <Stat
          sx={{ marginRight: "100px" }}
          label="Active Addons"
          description="1"
        />
        <Stat
          sx={{ marginRight: "100px" }}
          label="Storage"
          description="2.4 GB"
        />
      </Flex>
      {devAddons.map((addon) => (
        <Addon {...addon} />
      ))}
    </AppLayout>
  );
}

// TODO: Make component look good lol
function Addon({ name, activated, description, id }: IAddon) {
  const router = useRouter();
  const link = router.pathname + "/id";
  return (
    <Flex key={id}>
      <h3>{name}</h3>
      <h5>{description}</h5>
      <Button onClick={() => router.push(link)}>
        {activated ? "Manage" : "Enable"}
      </Button>
    </Flex>
  );
}

export interface IAddon {
  name: string;
  activated: boolean;
  description: string;
  id: string;
}

export const devAddons: IAddon[] = [
  {
    id: "0",
    name: "PostgresSQL",
    activated: false,
    description:
      "PostgreSQL is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. Postgres >>>",
  },
  {
    id: "1",
    name: "MongoDB",
    activated: true,
    description:
      "MongoDB is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. MongoDB >>>",
  },
  {
    id: "1",
    name: "Redis",
    activated: false,
    description:
      "Redis is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. Redis >>>",
  },
];
