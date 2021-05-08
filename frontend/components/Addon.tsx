import { KVEntry } from "./KVEntry";
import { ConfirmDelete } from "./ConfirmDelete";
import { Stat } from "./Stat";
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
import { devAddons, devAddonsOriginal } from "../lib/dummyData";
import _ from "lodash";

export function Addon({
  name,
  activated: a,
  description,
  id,
  img,
  storage,
  config: c,
}: IAddon) {
  const [activated, updateActive] = useState(a);
  const [newConfig, updateConfig] = useState(c);
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
  const {
    isOpen: confirmIsOpen,
    onOpen: confirmOnOpen,
    onClose: confirmOnClose,
  } = useDisclosure();

  const [verb, setVerb] = useState("disable");

  function closeAndDiscard() {
    manageOnClose();
    devAddons[id] = devAddonsOriginal[id];
    updateConfig(c);
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

      <Modal isOpen={manageIsOpen} onClose={closeAndDiscard}>
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
            <Stat label="Storage" description={storage} />
            <Text margin="initial" padding="initial" my="0.5em">
              Any unsaved changes will be discarded.
            </Text>
            {Object.entries(newConfig).map((entry) => {
              const kv_id = entry[0];
              const v = entry[1];
              let obj = {};
              obj[kv_id] = v;
              return (
                <KVEntry
                  key={kv_id}
                  id={kv_id}
                  onDataChange={(entry: KVConfig) => {
                    updateConfig({ ...newConfig, ...entry });
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
              mx="1em"
              onClick={() => {
                setVerb("disable");
                manageOnClose();
                confirmOnOpen();
              }}
            >
              Disable Addon
            </Button>
            <Button
              margin="initial"
              mx="1em"
              variant="ghost"
              colorScheme="blue"
              px="1.5em"
              onClick={() => {
                setVerb("wipe");
                manageOnClose();
                confirmOnOpen();
              }}
            >
              Wipe Data Only
            </Button>
            <Button
              margin="initial"
              ml="1em"
              variant="ghost"
              colorScheme="blue"
              px="1.5em"
              onClick={() => {
                devAddons[id]["config"] = newConfig;
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
                  storage,
                  description,
                };
                const idx = devAddons.findIndex((o) => _.isEqual(o, obj));
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

      <ConfirmDelete
        isOpen={confirmIsOpen}
        onClose={confirmOnClose}
        onOpen={confirmOnOpen}
        verb={verb}
        name={name}
        onCancellation={manageOnOpen}
        onConfirmation={() => {
          const obj = {
            name,
            activated,
            id,
            config: c,
            img,
            storage,
            description,
          };
          const idx = devAddons.findIndex((o) => _.isEqual(o, obj));
          if (verb != "wipe") {
            devAddons[idx].activated = false;
            updateActive(false);
          }
        }}
      />
    </>
  );
}
