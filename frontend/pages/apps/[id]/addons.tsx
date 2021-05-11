import AppLayout from "../../../layouts/app";

import { Stat } from "../../../components/Stat";
import { Addon } from "../../../components/Addon";

import { Flex } from "@chakra-ui/react";
import { devAddons } from "../../../lib/dummyData";

export default function AppAddonOverview() {
  return (
    <AppLayout selected="Addons">
      <Flex>
        <Stat
          style={{ marginRight: "100px" }}
          label="Active Addons"
          description="1"
        />
        <Stat
          style={{ marginRight: "100px" }}
          label="Storage"
          description="2.4 GB"
        />
      </Flex>
      {devAddons.map((addon) => (
        <Addon {...addon} key={addon.id} />
      ))}
    </AppLayout>
  );
}
