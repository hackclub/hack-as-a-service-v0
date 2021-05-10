import { createContext } from "react";

export const ThemeContext = createContext({
  themeEnabled: true,
  setThemeEnabled: (_: boolean) => {},
});
export default ThemeContext;
