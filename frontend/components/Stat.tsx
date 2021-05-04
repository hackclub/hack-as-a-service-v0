// import { Flex, Text, Heading } from "@theme-ui/components";
import { Flex, Text, Heading } from "@chakra-ui/react";
// import { SxProp } from "@theme-ui/core";

export function Stat({
  label,
  description,
  style,
}: {
  style?: object;
  label: string;
  description: string;
}) {
  return (
    <Flex style={{ flexDirection: "column", ...style }}>
      <Text m="0" style={{ fontSize: "20px" }}>
        {label}
      </Text>
      <Heading style={{ fontSize: "40px" }}>{description}</Heading>
    </Flex>
  );
}
