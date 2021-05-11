import { useContext, useEffect } from "react";
import ThemeContext from "../lib/disable_theme";

export default function DisableTheme() {
  const { setThemeEnabled } = useContext(ThemeContext);

  useEffect(() => {
    setThemeEnabled(false);
  });

  return <></>;
}
