import { ThemeProvider, Theme, merge } from "theme-ui";
import theme from "@hackclub/theme";

const haasTheme: Theme = merge(theme as Theme, {
  forms: {
    input: {
      border: "2px solid grey",
    },
  },
  images: {
    avatar: {
      boxShadow: "0 4px 12px 0 rgba(0,0,0,.1)",
    },
  },
});

function MyApp({ Component, pageProps }) {
  return (
    <ThemeProvider theme={haasTheme}>
      <Component {...pageProps} />
    </ThemeProvider>
  );
}

export default MyApp;
