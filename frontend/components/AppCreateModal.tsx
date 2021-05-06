import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Heading,
  IconButton,
  useDisclosure,
  Button,
  FormControl,
  FormLabel,
  FormHelperText,
  Input,
  Text,
  FormErrorMessage,
} from "@chakra-ui/react";

import { Formik } from "formik";

export default function AppCreateModal({ isOpen, onClose, onSubmit }) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />
      <Formik
        initialValues={{ id: "", name: "" }}
        onSubmit={onSubmit}
        validate={(values) => {
          const errors: { id?: string } = {};
          if (!values.id) {
            errors.id = "This field is required.";
          } else if (!/^[a-z0-9][^/:_A-Z\s]*$/.test(values.id)) {
            errors.id = "Your app ID can't contain spaces or most puncuation.";
          }

          return errors;
        }}
      >
        {({ handleChange, handleBlur, values, handleSubmit, errors }) => (
          <ModalContent>
            <ModalHeader py={1}>
              <Heading as="h1" my={1} fontWeight="normal">
                Create an app
              </Heading>
            </ModalHeader>
            <ModalCloseButton float="right" top={5} right={2} />
            <ModalBody py={1}>
              <form>
                <FormControl id="id" isRequired my={1}>
                  <FormLabel mb={1}>App ID</FormLabel>
                  <Input
                    type="text"
                    onChange={handleChange}
                    onBlur={handleBlur}
                    value={values.id}
                  />
                  {errors.id && (
                    <Text color="red" my={0}>
                      {errors.id}
                    </Text>
                  )}
                  {values.id != "" && !errors.id && (
                    <FormHelperText mt={1}>
                      Your app will be accessible at{" "}
                      <Text display="inline" fontWeight="bold" color="inherit">
                        {values.id}
                      </Text>
                      .haas.hackclub.com
                    </FormHelperText>
                  )}
                </FormControl>

                <FormControl id="name">
                  <FormLabel mb={1}>App Name</FormLabel>
                  <Input
                    type="text"
                    onChange={handleChange}
                    onBlur={handleBlur}
                    value={values.name}
                  />
                  {errors.name && (
                    <Text color="red" my={0}>
                      {errors.name}
                    </Text>
                  )}
                  <FormHelperText mt={1}>
                    An optional name for your app
                  </FormHelperText>
                </FormControl>
              </form>
            </ModalBody>

            <ModalFooter py={1.5}>
              <Button px={2} onClick={onClose}>
                Cancel
              </Button>
              <Button
                variant="cta"
                px={2}
                ml={1}
                onClick={() => handleSubmit()}
                isLoading
              >
                Create
              </Button>
            </ModalFooter>
          </ModalContent>
        )}
      </Formik>
    </Modal>
  );
}
