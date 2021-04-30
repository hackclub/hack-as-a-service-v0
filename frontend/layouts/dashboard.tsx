import Icon from "@hackclub/icons";
import Link from "next/link";
import { useRouter } from "next/router";
import { PropsWithChildren, useEffect } from "react";
import useSWR from "swr";
import { Avatar, Box, Container, Flex, Heading, SxProp, Text } from "theme-ui";
import fetchApi from "../lib/fetch";

import { Glyph } from "../types/glyph";

function SidebarItem({
  image,
  icon,
  children,
  url,
  sx,
}: PropsWithChildren<ISidebarItem> & SxProp) {
  const item = (
    <Flex sx={{ alignItems: "center", cursor: "pointer", ...sx }} my={10}>
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
          <Icon glyph={icon as Glyph} />
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

  if (url) {
    return <Link href={url}>{item}</Link>;
  }

  return item;
}

function SidebarSection({
  title,
  items,
}: {
  title?: string;
  items: ISidebarItem[];
}) {
  return (
    <Box mt={4}>
      {title && <Heading mb={3}>{title}</Heading>}
      {items.map((item) => {
        return (
          <SidebarItem key={item.text} {...item}>
            {item.text}
          </SidebarItem>
        );
      })}
    </Box>
  );
}

function SidebarHeader({ avatar }: { avatar?: string }) {
  return (
    <Flex
      sx={{ alignItems: "center", position: "sticky", top: 0 }}
      py="24px"
      px="50px"
      bg="background"
    >
      <Avatar src={avatar} />
      <Box sx={{ flexGrow: 1 }} />
      <Icon glyph="controls" size={32} style={{ margin: "0 10px" }} />
      <Link href="/logout">
        <Icon
          glyph="door-leave"
          size={32}
          style={{ margin: "0 10px", cursor: "pointer" }}
        />
      </Link>
    </Flex>
  );
}

export interface ISidebarSection {
  title?: string;
  items: ISidebarItem[];
}

export interface ISidebarItem {
  image?: string;
  icon?: string;
  text: string;
  url?: string;
}

export default function DashboardLayout({
  title,
  image,
  sidebarSections,
  children,
}: PropsWithChildren<{
  title: string;
  image?: string;
  sidebarSections: ISidebarSection[];
}>) {
  const router = useRouter();
  const { data: user, error: userError } = useSWR("/users/me", fetchApi);

  useEffect(() => {
    if (userError && NODE_ENV !== "development") {
      router.push("/");
    }
  }, [userError]);

  return (
    <Flex sx={{ minHeight: "100vh" }}>
      <Box
        sx={{ flexBasis: 400, flexShrink: 0, flexGrow: 0 }}
        // px="50px"
        py="30px"
      >
        <SidebarHeader avatar={user?.user.Avatar} />
        <Box mt="40px" px="50px">
          {sidebarSections.map((v, i) => {
            return <SidebarSection key={i} title={v.title} items={v.items} />;
          })}
        </Box>
      </Box>
      <Box sx={{ flex: "1 1 auto" }} px="50px" py="35px">
        <Flex
          sx={{ alignItems: "center", position: "sticky", top: 0 }}
          py={3}
          bg="background"
        >
          {image && (
            <Avatar
              size={64}
              src={image}
              sx={{ borderRadius: 8 }}
              bg="sunken"
              mr={4}
            />
          )}

          <Heading as="h1" sx={{ fontSize: 50 }}>
            {title}
          </Heading>
        </Flex>

        {children}
      </Box>
    </Flex>
  );
}
