import {
  Button,
  Text,
  Heading,
  Input,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
} from "@chakra-ui/react";

import { useState } from "react";
import _ from "lodash";

export function ConfirmDelete(props: {
  name: string;
  onConfirmation: Function;
  onCancellation: Function;
  buttonText?: string;
  verb?: string;
  isOpen: boolean;
  onOpen: Function;
  onClose: () => void;
}) {
  let ranCancel = false;
  const {
    name,
    onConfirmation,
    onCancellation,
    buttonText,
    isOpen,
    verb,
    onClose,
  } = props;

  const [value, setValue] = useState("");
  const [markedForDeletion, setDelete] = useState(false);

  function handleChange(evt) {
    const v: string = evt.target.value;
    setValue(v);
    v.trim().includes(name.trim()) ? setDelete(true) : setDelete(false);
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={() => {
        onClose();
        onCancellation();
      }}
    >
      <ModalOverlay />
      <ModalContent padding="2em">
        <ModalHeader
          padding="unset"
          display="flex"
          justifyContent="space-between"
          alignItems="center"
        >
          <Heading py="0.5em">
            Are you sure you want to {verb ?? "delete"} {name}?
          </Heading>
          <ModalCloseButton
            margin="initial"
            padding="initial"
            position="unset"
            onClick={() => onCancellation()}
          />
        </ModalHeader>

        <ModalBody margin="initial" padding="initial">
          <Text margin="initial" padding="initial">
            Type {name} into the box below to confirm this action.
          </Text>
          <Input
            margin="initial"
            my="0.5em"
            onChange={handleChange}
            value={value}
            placeholder={name}
          />
          <Button
            margin="initial"
            px="1.5em"
            isDisabled={!markedForDeletion}
            onClick={() => {
              setValue("");
              onClose();
              onConfirmation();
            }}
          >
            {buttonText ??
              verb.toUpperCase().substr(0, 1) + verb.substr(1) ??
              "Delete"}
          </Button>
        </ModalBody>
        <ModalFooter margin="initial" padding="initial">
          <Button
            margin="initial"
            px="1.5em"
            variant="ghost"
            colorScheme="blue"
            onClick={() => {
              ranCancel = false;
              onClose();
              onCancellation();
            }}
          >
            Cancel
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
