import { KVEntry } from "./KVEntry";
import {
  Flex,
  Button,
  Text,
  Heading,
  Img,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useDisclosure,
} from "@chakra-ui/react";

import { useState } from "react";
import { IAddon, KVConfig } from "../types/haas";
import { devAddons } from "../lib/dummyData";
import _ from "lodash";

// TODO: Make component look good lol
export function Addon({
  name,
  activated: a,
  description,
  id,
  img,
  config: c,
}: IAddon) {
  const [activated, updateActive] = useState(a);
  const [config, updateConfig] = useState(c);
  const {
    isOpen: manageIsOpen,
    onOpen: manageOnOpen,
    onClose: manageOnClose,
  } = useDisclosure();
  const {
    isOpen: enableIsOpen,
    onOpen: enableOnOpen,
    onClose: enableOnClose,
  } = useDisclosure();
  function saveAddonData() {
    const ogKey = Object.keys(c)[0];
    if (ogKey != Object.keys(config)[0]) {
    }
  }
  return (
    <>
      <Flex
        margin="2em"
        borderRadius="xl"
        borderColor="#00000033"
        borderWidth="thin"
        padding="2"
        alignItems="center"
      >
        <Flex flexDirection="column">
          <Flex alignItems="center">
            <Img objectFit="contain" boxSize="3em" src={img} />
            <Heading paddingLeft="0.25em">{name}</Heading>
          </Flex>
          <Text margin="unset">{description}</Text>
        </Flex>
        <Button px="1.75em" onClick={activated ? manageOnOpen : enableOnOpen}>
          {activated ? "Manage" : "Enable"}
        </Button>
      </Flex>

      <Modal isOpen={manageIsOpen} onClose={manageOnClose}>
        <ModalOverlay />
        <ModalContent padding="2em">
          <ModalHeader
            padding="unset"
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Heading py="0.5em">Addons → {name}</Heading>
            <ModalCloseButton
              margin="initial"
              padding="initial"
              position="unset"
            />
          </ModalHeader>
          <ModalBody margin="initial" padding="initial">
            {Object.entries(config).map((entry) => {
              const kv_id = entry[0];
              const v = entry[1];
              let obj = {};
              obj[kv_id] = v;
              return (
                <KVEntry
                  key={kv_id}
                  id={kv_id}
                  onDataChange={(entry: KVConfig) => {
                    updateConfig({ ...config, ...entry });
                  }}
                  entry={obj}
                />
              );
            })}
          </ModalBody>

          <ModalFooter margin="initial" padding="initial">
            <Button
              margin="initial"
              variant="ghost"
              colorScheme="blue"
              px="1.5em"
              onClick={() => {
                console.log("do the saving thing and close");
                manageOnClose();
              }}
            >
              Save
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Modal isOpen={enableIsOpen} onClose={enableOnClose}>
        <ModalOverlay />
        <ModalContent padding="2em">
          <ModalHeader
            padding="unset"
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Heading py="0.5em">Enable {name}?</Heading>
            <ModalCloseButton
              margin="initial"
              padding="initial"
              position="unset"
            />
          </ModalHeader>
          <ModalBody margin="initial" padding="initial">
            You will be billed for all resources used by this addon.
          </ModalBody>

          <ModalFooter margin="initial" padding="initial">
            <Button
              margin="initial"
              variant="ghost"
              colorScheme="blue"
              px="1.5em"
              onClick={() => {
                const obj = {
                  name,
                  activated,
                  id,
                  config: c,
                  img,
                  description,
                };
                const idx = devAddons.findIndex((o) => _.isEqual(o, obj));
                console.log(idx, obj, devAddons);
                devAddons[idx].activated = true;
                updateActive(true);
                enableOnClose();
              }}
            >
              Confirm
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}
