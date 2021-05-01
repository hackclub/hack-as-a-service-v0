import { IconButton } from "theme-ui";
import { Moon, Sun } from "react-feather";
import { useColorMode } from "theme-ui";

function changeToggleButton({ ...props }) {
  const [colorMode] = useColorMode()

 let icon;
 
  if(colorMode === 'dark') {
    icon = <Sun />
  } else {
    icon = <Moon />
  }
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
    >
      {icon}
    </IconButton>
    )
  ;
}


export default changeToggleButton;
