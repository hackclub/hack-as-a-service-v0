import Icon from "@hackclub/icons";
import { PropsWithChildren, useEffect } from "react";
import useSWR from "swr";
import { Avatar, Box, Container, Flex, Heading, SxProp, Text } from "theme-ui";
import fetchApi from "../lib/fetch";

function SidebarItem({
  image,
  icon,
  children,
  sx,
}: PropsWithChildren<{ image?: string; icon?: string }> & SxProp) {
  return (
    <Flex sx={{ alignItems: "center", ...sx }} my={10}>
      {image ? (
        <Avatar src={image} sx={{ borderRadius: 8 }} bg="sunken" mr={15} />
      ) : (
        <Flex
          sx={{
            height: 48,
            width: 48,
            borderRadius: 8,
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 4px 12px 0 rgba(0,0,0,.1)",
          }}
          mr={15}
          bg="sunken"
        >
          <Icon glyph={icon} />
        </Flex>
      )}
      <Heading
        as="h3"
        sx={{
          fontWeight: "normal",
          whiteSpace: "nowrap",
          overflow: "hidden",
          textOverflow: "ellipsis",
        }}
      >
        {children}
      </Heading>
    </Flex>
  );
}

function SidebarSection({
  title,
  items,
}: {
  title?: string;
  items: { image?: string; icon?: string; text: string }[];
}) {
  return (
    <Box mt={4}>
      {title && <Heading mb={3}>{title}</Heading>}
      {items.map((v, i) => {
        return (
          <SidebarItem key={i} icon={v.icon} image={v.image}>
            {v.text}
          </SidebarItem>
        );
      })}
    </Box>
  );
}

function SidebarHeader({ avatar }: { avatar?: string }) {
  return (
    <Flex sx={{ alignItems: "center" }}>
      <Avatar src={avatar} />
      <Box sx={{ flexGrow: 1 }} />
      <Icon glyph="controls" size={32} style={{ margin: "0 10px" }} />
      <Icon glyph="door-leave" size={32} style={{ margin: "0 10px" }} />
    </Flex>
  );
}

export interface SidebarSection {
  title?: string;
  items: { image?: string; icon?: string; text: string }[];
}

export default function DashboardLayout({
  title,
  sidebarSections,
  children,
}: PropsWithChildren<{ title: string; sidebarSections: SidebarSection[] }>) {
  const { data: user } = useSWR("/users/me", fetchApi);

  return (
    <>
      <Flex sx={{ minHeight: "100vh" }}>
        <Box
          sx={{ flexBasis: 400, flexShrink: 0, flexGrow: 0 }}
          px="50px"
          py="30px"
        >
          <SidebarHeader avatar={user?.user.Avatar} />
          <Box mt={50}>
            {sidebarSections.map((v, i) => {
              return <SidebarSection key={i} title={v.title} items={v.items} />;
            })}
          </Box>
        </Box>
        <Box sx={{ flex: "1 1 auto" }} px="50px" py="30px">
          <Heading as="h1" mb={5}>
            {title}
          </Heading>

          {children}
        </Box>
      </Flex>
    </>
  );
}
