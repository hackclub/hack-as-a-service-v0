import {
  Button,
  Input,
  InputGroup,
  InputRightElement,
  InputLeftAddon,
  Text,
} from "@chakra-ui/react";

import { useState, useEffect } from "react";
import { KVConfig } from "../types/haas";

export function KVEntry(props: {
  entry: KVConfig;
  id: string;
  onDataChange: Function;
}) {
  const { entry, id, onDataChange } = props;
  const { obscureValue, keyEditable, valueEditable, key, value } = entry[id];
  const [hide, toggleVal] = useState(obscureValue);
  const [val, updateVal] = useState(value);
  const [newKey, updateKey] = useState(key);
  useEffect(() => {
    runCallback();
  }, [val]);

  useEffect(() => {
    runCallback();
  }, [newKey]);

  const handleChange = (evt) => {
    updateVal(evt.target.value);
  };
  const handleKeyChange = (event) => {
    updateKey(event.target.value);
  };
  function runCallback() {
    let obj = {};
    obj[id] = {
      obscureValue,
      keyEditable,
      valueEditable,
      key: newKey,
      value: val,
    };
    onDataChange(obj);
  }
  return (
    <>
      <InputGroup size="md" mb="8px">
        <InputLeftAddon margin="initial" padding="initial">
          {keyEditable ? (
            <Input
              isReadOnly={!keyEditable}
              isDisabled={!keyEditable}
              onChange={handleKeyChange}
              value={newKey}
            />
          ) : (
            <Text px="0.5em">{newKey}</Text>
          )}
        </InputLeftAddon>
        {/* this is here in an attempt to stop browsers from offering to save passwords */}
        {hide && <Input autoComplete="off" display="none" type="password" />}
        <Input
          pr="4.5rem"
          isReadOnly={!valueEditable}
          isDisabled={!valueEditable}
          type={!hide ? "text" : "password"}
          placeholder="Enter your value here..."
          value={val}
          autoComplete={"off"}
          onChange={handleChange}
        />
        <InputRightElement width="4.5rem">
          {obscureValue && (
            <Button
              h="1.75rem"
              size="sm"
              margin="unset"
              px="1.5em"
              mr="0.25em"
              onClick={() => toggleVal(!hide)}
            >
              {!hide ? "Hide" : "Show"}
            </Button>
          )}
        </InputRightElement>
      </InputGroup>
    </>
  );
}
