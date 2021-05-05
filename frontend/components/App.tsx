import { LinkBox, Flex, Heading, Text, useColorMode } from "@chakra-ui/react";
import Link from "next/link";

export default function App({
  name,
  shortName,
  url,
}: {
  name: string;
  shortName: string;
  url: string;
}) {
  const { colorMode } = useColorMode();

  return (
    <Link href={url} passHref>
      <LinkBox>
        <Flex
          alignItems="flex-start"
          justifyContent="center"
          flexDirection="column"
          borderRadius="10px"
          cursor="pointer"
          bg={colorMode == "dark" ? "slate" : "sunken"}
          p="30px"
        >
          <Heading as="h2" sx={{ fontWeight: "normal" }}>
            {name}
          </Heading>
          <Text color={colorMode == "dark" ? "smoke" : "muted"} my={1}>
            ({shortName})
          </Text>
        </Flex>
      </LinkBox>
    </Link>
  );
}
