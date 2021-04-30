import { IconButton } from "theme-ui";
import { Moon } from "react-feather";

const ColorButton = ({ ...props }) => (
  <IconButton
    as="button"
    {...props}
    sx={{
      backgroundColor: "transparent",
      border: 0,
      color: "inherit",
      cursor: "pointer",
      position: "absolute",
      left: 200, // this is not very responsive css, we need to fix it later lmao
      top: 26,
    }}
  >
    <Moon size={26} />
  </IconButton>
);

export default ColorButton;
