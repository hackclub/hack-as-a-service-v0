import {
  Flex,
  Text,
  Heading
} from "@theme-ui/components";
import { SxProp } from "@theme-ui/core";

export function Stat({
  label,
  description,
  sx,
}: { label: string; description: string; } & SxProp) {
  return (
    <Flex sx={{ flexDirection: "column", ...sx }}>
      <Text sx={{ fontSize: "20px" }}>{label}</Text>
      <Heading sx={{ fontSize: "40px" }}>{description}</Heading>
    </Flex>
  );
}
