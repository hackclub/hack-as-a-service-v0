import { Moon, Sun } from "react-feather";
import { IconButton, useColorMode } from "@chakra-ui/react";

function ColorButton({ ...props }) {
  const { colorMode, toggleColorMode } = useColorMode();

  const icon = colorMode === "dark" ? <Sun /> : <Moon />;

  return (
    <IconButton
      as="button"
      aria-label="Toggle color mode"
      mx="5px"
      background="inherit"
      onClick={toggleColorMode}
      {...props}
    >
      {icon}
    </IconButton>
  );
}

export default ColorButton;
