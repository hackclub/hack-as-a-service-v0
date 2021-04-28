import { ThemeProvider, Theme, merge } from "theme-ui";
import theme from "@hackclub/theme";

const haasTheme: Theme = merge(
  {
    images: {
      avatar: {
        boxShadow: "0 4px 12px 0 rgba(0,0,0,.1)",
      },
    },
  },
  theme as Theme
);

function MyApp({ Component, pageProps }) {
  return (
    <ThemeProvider theme={haasTheme}>
      <Component {...pageProps} />
    </ThemeProvider>
  );
}

export default MyApp;
