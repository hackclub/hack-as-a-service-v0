import { IconButton } from "theme-ui";
import { Moon, Sun } from "react-feather";
import { useColorMode } from "theme-ui";

function ColorButton({ ...props }) {
  const [colorMode, setColorMode] = useColorMode();

  const icon = colorMode === "dark" ? <Sun /> : <Moon />;

  return (
    <IconButton
      as="button"
      {...props}
      sx={{
        backgroundColor: "transparent",
        border: 0,
        color: "inherit",
        cursor: "pointer",
        margin: "0 10px",
      }}
      onClick={(e) => {
        setColorMode(colorMode === "default" ? "dark" : "default");
      }}
    >
      {icon}
    </IconButton>
  );
}

export default ColorButton;
