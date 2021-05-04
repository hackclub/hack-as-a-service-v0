import Icon from "@hackclub/icons";
import Link from "next/link";
import { useRouter } from "next/router";
import { PropsWithChildren, ReactElement, useEffect } from "react";
import useSWR from "swr";
import {
  Avatar,
  Box,
  Flex,
  Heading,
  IconButton,
  SystemStyleObject,
} from "@chakra-ui/react";
import { Glyph } from "../types/glyph";
import ColorSwitcher from "../components/ColorButton";

function SidebarItem({
  image,
  icon,
  children,
  url,
  sx,
  selected,
}: PropsWithChildren<ISidebarItem> & { sx?: SystemStyleObject }) {
  let imageComponent: ReactElement;

  if (image) {
    imageComponent = (
      <Avatar src={image} sx={{ borderRadius: 8 }} bg="sunken" mr={15} />
    );
  } else if (icon) {
    imageComponent = (
      <Flex
        width="48px"
        height="48px"
        borderRadius={8}
        alignItems="center"
        justifyContent="center"
        boxShadow="0 4px 12px 0 rgba(0,0,0,.1)"
        backgroundImage={
          selected ? "linear-gradient(-45deg, #ec3750, #ff8c37)" : null
        }
        mr={15}
        bg="sunken"
      >
        <Icon glyph={icon} color={selected ? "white" : null} />
      </Flex>
    );
  }

  const item = (
    <Flex
      alignItems="center"
      sx={{
        alignItems: "center",
        ...(url ? { cursor: "pointer" } : {}),
        ...sx,
      }}
      my="10px"
    >
      {(image || icon) && imageComponent}
      <Heading
        as="h3"
        size="md"
        sx={{
          fontWeight: "normal",
          ...(image || icon
            ? { whiteSpace: "nowrap", overflow: "hidden" }
            : {}),
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
      {title && (
        <Heading size="md" mb={3} mt={6}>
          {title}
        </Heading>
      )}
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
    <Flex alignItems="center" position="sticky" top={0} py="24px" px="50px">
      <Avatar src={avatar} />
      <Box flexGrow={1} />
      <ColorSwitcher />
      <IconButton mx="5px" aria-label="Controls" background="inherit">
        <Icon glyph="controls" size={32} />
      </IconButton>
      <IconButton mx="5px" aria-label="Log out" background="inherit">
        <Link href="/logout">
          <Icon glyph="door-leave" size={32} />
        </Link>
      </IconButton>
    </Flex>
  );
}

export interface ISidebarSection {
  title?: string;
  items: ISidebarItem[];
}

export interface ISidebarItem {
  image?: string;
  icon?: Glyph;
  text: string;
  url?: string;
  selected?: boolean;
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
  const { data: user, error: userError } = useSWR("/users/me");

  useEffect(() => {
    if (userError && process.env.NODE_ENV !== "development") {
      router.push("/login");
    }
  }, [userError]);

  return (
    <Flex minHeight="100vh">
      <Box
        flexBasis={400}
        flexShrink={0}
        flexGrow={0}
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
      <Box flex={"1 1 auto"} px="50px" py="35px">
        <Flex alignItems="center" position="sticky" top={0} py={3}>
          {image && (
            <Avatar size="md" src={image} borderRadius={8} bg="sunken" mr={4} />
          )}

          <Heading as="h1" fontSize={50}>
            {title}
          </Heading>
        </Flex>

        {children}
      </Box>
    </Flex>
  );
}
