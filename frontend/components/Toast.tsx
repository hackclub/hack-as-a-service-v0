import { Flex, Text } from "@chakra-ui/react";
import Icon from "@hackclub/icons";

export function ErrorToast({ text }: { text: string }) {
  return (
    <Flex
      p={1}
      bg="red"
      color="white"
      borderRadius={10}
      fontSize={20}
      alignItems="center"
    >
      <Icon glyph="important" style={{ marginRight: 20 }} />
      <Text p={0} m={0} fontSize={20} color="white">
        {text}
      </Text>
    </Flex>
  );
}

export function SuccessToast({ text }: { text: string }) {
  return (
    <Flex
      p={1}
      bg="green"
      color="black"
      borderRadius={10}
      fontSize={20}
      alignItems="center"
    >
      <Icon glyph="checkmark" style={{ marginRight: 20 }} />
      <Text p={0} m={0} fontSize={20} color="black">
        {text}
      </Text>
    </Flex>
  );
}
